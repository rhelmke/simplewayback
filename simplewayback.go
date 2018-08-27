// Package simplewayback is a simple library for querying the Wayback Machine CDX API and fetching snapshots
package simplewayback

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	neturl "net/url"
	"regexp"
	"strconv"
	"time"
)

const (
	cdxURL  = "https://web.archive.org/cdx/search/cdx?"
	dataURL = "http://web.archive.org/web"
)

type matchType int

type outputFormat int

type field int

// Matchtypes
const (
	// MatchTypeExact instructs the simplewayback package to return results matching exactly example.org/example.html
	MatchTypeExact matchType = iota
	// MatchTypePrefix instructs the simplewayback package to return results for all results under the path example.org/subdir/
	MatchTypePrefix
	// MatchTypeHost instructs the simplewayback package to return results from host example.org
	MatchTypeHost
	// MatchTypeDomain instructs the simplewayback package to return results from host example.org and all subhosts *.example.org
	MatchTypeDomain
)

// Output Formats
const (
	// OutputFormatJSON sets JSON as response format for the archive.org api ([["urlkey","timestamp","original","mimetype","statuscode","digest","length"],...])
	OutputFormatJSON outputFormat = iota
	// OutputFormatCDX sets CDX as response format for the archive.org api (urlkey timestamp original mimetype statuscode digest length
	OutputFormatCDX
)

// Errors
var (
	// ErrorInvalidMatchType...
	ErrorInvalidMatchType     = errors.New("simplewayback: Invalid matchType")
	ErrorInvalidURL           = errors.New("simplewayback: Invalid URL (this field is mandatory)")
	ErrorInvalidOutputFormat  = errors.New("simplewayback: Invalid OutputFormat")
	ErrorInvalidFromTo        = errors.New("simplewayback: Parameter 'to' must be larger than 'from'")
	ErrorInvalidField         = errors.New("simplewayback: Invalid field")
	ErrorInvalidNumber        = errors.New("simplewayback: Integer must be >= 0")
	ErrorPaginationResumption = errors.New("simplewayback: Pagination and Resumption Keys can not be enabled at the same time")
	ErrorInvalidScheme        = errors.New("simplewayback: The provided URL must use 'http', 'https' or '' as scheme")
	ErrorBadResponse          = errors.New("simplewayback: Bad Response from Wayback Machine API (!200)")
)

// RegexFields
const (
	FieldURLKey field = iota
	FieldTimestamp
	FieldOriginal
	FieldMimetype
	FieldStatuscode
	FieldDigest
	FieldLength
)

var fields = map[field]string{
	FieldURLKey:     "urlkey",
	FieldTimestamp:  "timestamp",
	FieldOriginal:   "original",
	FieldMimetype:   "mimetype",
	FieldStatuscode: "statuscode",
	FieldDigest:     "digest",
	FieldLength:     "length",
}

var matchTypes = map[matchType]string{
	MatchTypeExact:  "exact",
	MatchTypePrefix: "prefix",
	MatchTypeHost:   "host",
	MatchTypeDomain: "domain",
}

var outputFormats = map[outputFormat]string{
	OutputFormatJSON: "json",
	OutputFormatCDX:  "cdx",
}

// CDXAPI is a wrapper for polling the CDX-API of Wayback Machine
type CDXAPI struct {
	params           *neturl.Values
	regFilterKeys    []string
	collapsingKeys   []string
	useResumptionKey bool
	resumptionKey    string
	usePagination    bool
	page             int
	apiKey           string
	urlBuf           *bytes.Buffer
}

// NewCDXAPI creates and initializes a new CDX API wrapper
func NewCDXAPI(url string) (*CDXAPI, error) {
	cdx := &CDXAPI{params: &neturl.Values{}, urlBuf: &bytes.Buffer{}}
	if err := cdx.SetURL(url); err != nil {
		return nil, err
	}
	return cdx, nil
}

// SetAPIKey sets an optional API key
func (cdx *CDXAPI) SetAPIKey(apiKey string) error {
	cdx.apiKey = apiKey
	return nil
}

// APIKey getter
func (cdx *CDXAPI) APIKey() string {
	return cdx.apiKey
}

// SetMatchType where mType = MatchTypeExact | MatchTypePrefix | MatchTypeHost | MatchTypeDomain
func (cdx *CDXAPI) SetMatchType(mType matchType) error {
	if _, ok := matchTypes[mType]; !ok {
		return ErrorInvalidMatchType
	}
	cdx.params.Set("matchType", matchTypes[mType])
	return nil
}

// MatchType getter
func (cdx *CDXAPI) MatchType() int {
	retval := MatchTypeExact
	switch cdx.params.Get("matchType") {
	case matchTypes[MatchTypeExact]:
		retval = MatchTypeExact
	case matchTypes[MatchTypePrefix]:
		retval = MatchTypePrefix
	case matchTypes[MatchTypeHost]:
		retval = MatchTypeHost
	case matchTypes[MatchTypeDomain]:
		retval = MatchTypeDomain
	}
	return int(retval)
}

// ResetMatchType resets the MatchType (default: MatchTypeExact)
func (cdx *CDXAPI) ResetMatchType() {
	cdx.params.Set("matchType", matchTypes[MatchTypeExact])
}

// SetOutputFormat sets the output format
func (cdx *CDXAPI) SetOutputFormat(format outputFormat) error {
	if _, ok := outputFormats[format]; !ok {
		return ErrorInvalidOutputFormat
	}
	if format == OutputFormatCDX {
		cdx.params.Del("output")
	} else {
		cdx.params.Set("output", outputFormats[format])
	}
	return nil
}

// OutputFormat getter
func (cdx *CDXAPI) OutputFormat() int {
	if cdx.params.Get("output") == outputFormats[OutputFormatJSON] {
		return int(OutputFormatJSON)
	}
	return int(OutputFormatCDX)
}

// ResetOutputFormat resets the output format (default: OutputFormatCDX)
func (cdx *CDXAPI) ResetOutputFormat() {
	cdx.params.Del("output")
}

// SetURL to search for
func (cdx *CDXAPI) SetURL(url string) error {
	parsed, err := neturl.Parse(url)
	if err != nil {
		return err
	}
	if !(parsed.Scheme == "http" || parsed.Scheme == "https" || parsed.Scheme == "") {
		return ErrorInvalidScheme
	}
	cdx.params.Set("url", url)
	return nil
}

// URL getter
func (cdx *CDXAPI) URL() string {
	return cdx.params.Get("url")
}

// SetLimit sets a limit
func (cdx *CDXAPI) SetLimit(limit int) error {
	if limit <= 0 {
		return ErrorInvalidNumber
	}
	cdx.params.Set("limit", strconv.Itoa(limit))
	return nil
}

// Limit getter
func (cdx *CDXAPI) Limit() int {
	lims := cdx.params.Get("limit")
	if lims == "" {
		return -1
	}
	lim, err := strconv.Atoi(lims)
	if err != nil {
		return -1
	}
	return lim
}

// ResetLimit resets the limit (default: no limit)
func (cdx *CDXAPI) ResetLimit() {
	cdx.params.Del("limit")
}

// AddRegexFilter to the wayback machine query.
// Regex filtering: It is possible to filter on a specific field or the entire CDX line (which is space delimited).
// Filtering by specific field is often simpler
func (cdx *CDXAPI) AddRegexFilter(fld field, regex string, negate bool) error {
	if _, ok := fields[fld]; !ok {
		return ErrorInvalidField
	}
	if _, err := regexp.Compile(regex); err != nil {
		return err
	}
	var buf bytes.Buffer
	if negate {
		buf.WriteString("!")
	}
	buf.WriteString(fields[fld])
	buf.WriteString(":")
	buf.WriteString(regex)
	key := fmt.Sprintf("filter%d", len(cdx.regFilterKeys))
	cdx.regFilterKeys = append(cdx.regFilterKeys, key)
	cdx.params.Set(key, buf.String())
	return nil
}

// ResetRegexFilters flushes all regex filters (default: no regex filters)
func (cdx *CDXAPI) ResetRegexFilters() {
	for i := range cdx.regFilterKeys {
		cdx.params.Del(cdx.regFilterKeys[i])
	}
	cdx.regFilterKeys = []string{}
}

// SetTimeFilter for the wayback machine query
func (cdx *CDXAPI) SetTimeFilter(from time.Time, to time.Time) error {
	if to.Sub(from) < 0 {
		return ErrorInvalidFromTo
	}
	cdx.params.Set("from", from.Format("20060102150405"))
	cdx.params.Set("to", to.Format("20060102150405"))
	return nil
}

// TimeFilter getter
func (cdx *CDXAPI) TimeFilter() (time.Time, time.Time) {
	from, err := time.Parse("20060102150405", cdx.params.Get("from"))
	if err != nil {
		return time.Time{}, time.Time{}
	}
	to, err := time.Parse("20060102150405", cdx.params.Get("to"))
	if err != nil {
		return time.Time{}, time.Time{}
	}
	return from, to
}

// ResetTimeFilter resets the time filter (default: no time filters)
func (cdx *CDXAPI) ResetTimeFilter() {
	cdx.params.Del("from")
	cdx.params.Del("to")
}

// AddCollapsing adds collapsing options to the Wayback Machine.
// A new form of filtering is the option to 'collapse' results based on a field, or a substring of a field. Collapsing is
// done on adjacent cdx lines where all captures after the first one that are duplicate are filtered out. This is useful
// for filtering out captures that are 'too dense' or when looking for unique captures.
func (cdx *CDXAPI) AddCollapsing(fld field, n int) error {
	if _, ok := fields[fld]; !ok {
		return ErrorInvalidField
	}
	if n < 0 {
		return ErrorInvalidNumber
	}
	var buf bytes.Buffer
	buf.WriteString(fields[fld])
	if n > 0 {
		buf.WriteString(":")
		buf.WriteString(fmt.Sprint(n))
	}
	key := fmt.Sprintf("collapse%d", len(cdx.collapsingKeys))
	cdx.collapsingKeys = append(cdx.collapsingKeys, key)
	cdx.params.Set(key, buf.String())
	return nil
}

// ResetCollapsing flushes all collapsing filters (default: no collapsing)
func (cdx *CDXAPI) ResetCollapsing() {
	for i := range cdx.collapsingKeys {
		cdx.params.Del(cdx.collapsingKeys[i])
	}
	cdx.collapsingKeys = []string{}
}

// SetGzip for gzipped response from archive.org
func (cdx *CDXAPI) SetGzip(enabled bool) error {
	if enabled {
		cdx.params.Del("gzip")
	} else {
		cdx.params.Set("gzip", "false")
	}
	return nil
}

// Gzip getter
func (cdx *CDXAPI) Gzip() bool {
	return !(cdx.params.Get("gzip") == "false")
}

// ResetGzip resets gzip (default: true)
func (cdx *CDXAPI) ResetGzip() {
	cdx.params.Del("gzip")
}

// SetOffset for querying data
func (cdx *CDXAPI) SetOffset(offset int) error {
	if offset <= 0 {
		return ErrorInvalidNumber
	}
	cdx.params.Set("offset", strconv.Itoa(offset))
	return nil
}

// Offset getter
func (cdx *CDXAPI) Offset() int {
	offs := cdx.params.Get("offset")
	if offs == "" {
		return -1
	}
	res, err := strconv.Atoi(offs)
	if err != nil {
		return -1
	}
	return res
}

// ResetOffset resets the offset (default: no offset)
func (cdx *CDXAPI) ResetOffset() {
	cdx.params.Del("offset")
}

// SetResumptionKey mode
func (cdx *CDXAPI) SetResumptionKey(enabled bool, key string) error {
	if cdx.usePagination {
		return ErrorPaginationResumption
	}
	cdx.useResumptionKey = enabled
	if enabled {
		cdx.params.Set("resumeKey", key)
		cdx.params.Set("showResumeKey", "true")
	} else {
		cdx.params.Del("resumeKey")
		cdx.params.Del("showResumeKey")
	}
	return nil
}

// ResumptionKeyEnabled checks whether the resumptionKey feature is enabled or not
func (cdx *CDXAPI) ResumptionKeyEnabled() bool {
	return cdx.useResumptionKey
}

// ResumptionKey getter
func (cdx *CDXAPI) ResumptionKey() string {
	return cdx.params.Get("resumeKey")
}

// ResetResumptionKey resets the resumption key settings (default: enabled=false, key="")
func (cdx *CDXAPI) ResetResumptionKey() {
	cdx.SetResumptionKey(false, "")
}

// SetPagination mode
func (cdx *CDXAPI) SetPagination(enabled bool, page int) error {
	if cdx.useResumptionKey {
		return ErrorPaginationResumption
	}
	if page < 0 {
		return ErrorInvalidNumber
	}
	if enabled {
		cdx.page = page
		cdx.params.Set("page", strconv.Itoa(page))
	} else {
		cdx.page = -1
		cdx.params.Del("page")
	}
	cdx.usePagination = enabled
	return nil
}

// PaginationEnabled checks whether the pagination features is enabled
func (cdx *CDXAPI) PaginationEnabled() bool {
	return cdx.usePagination
}

// PaginationPage getter
func (cdx *CDXAPI) PaginationPage() int {
	return cdx.page
}

// ResetPagination resets the pagination (default: enabled=false, page=-1)
func (cdx *CDXAPI) ResetPagination() {
	cdx.SetPagination(false, 0)
}

func (cdx *CDXAPI) buildURL(urlDst *bytes.Buffer) error {
	if cdx.params.Get("url") == "" {
		return ErrorInvalidURL
	}
	urlDst.Reset()
	urlDst.WriteString(cdxURL)

	// this is why setting collapse- and regex-filters is not straightforward.
	// The CDX-API requires that "collapse=" and "filter=" are placed multiple times
	// into the query URL-part. But url.Values is basically a map[string]string.
	// Adding multiple filters of the same type results in overwriting the previous
	// filter. So we need to hold unique Keys in url.Values that will be replaced
	// using following lines of code.
	encoded := cdx.params.Encode()
	if len(cdx.collapsingKeys) > 0 {
		encoded = regexp.MustCompile(`(collapse\d+)=`).ReplaceAllString(encoded, "collapse=")
	}
	if len(cdx.regFilterKeys) > 0 {
		encoded = regexp.MustCompile(`(filter\d+)=`).ReplaceAllString(encoded, "filter=")
	}
	urlDst.WriteString(encoded)
	return nil
}

// CDXRawQuery implements the reader interface to raw read a single search result
type CDXRawQuery struct {
	resp *http.Response
}

// Read interface implementation for CDXAPI
func (qry *CDXRawQuery) Read(p []byte) (int, error) {
	return qry.resp.Body.Read(p)
}

// RawPerform queries the CDX API and returns a CDXRawQuery that can be read using the Reader interface
func (cdx *CDXAPI) RawPerform() (*CDXRawQuery, error) {
	if err := cdx.buildURL(cdx.urlBuf); err != nil {
		return nil, err
	}
	client := http.Client{}
	req, err := http.NewRequest("GET", cdx.urlBuf.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "simplewayback/0.1 (https://github.com/rhelmke/simplewayback)")
	if cdx.apiKey != "" {
		req.AddCookie(&http.Cookie{Name: "cdx-auth-token", Value: cdx.apiKey})
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, ErrorBadResponse
	}
	return &CDXRawQuery{resp: resp}, nil
}

// CDXResult represents a single from the CDX API Response
type CDXResult struct {
	URLKey     string    `json:"url_key"`
	Timestamp  time.Time `json:"timestamp"`
	Original   string    `json:"original"`
	MimeType   string    `json:"mime_type"`
	StatusCode int       `json:"status_code"`
	Digest     string    `json:"digest"`
	Length     int       `json:"length"`
	Data       io.Reader `json:"-"`
}

// CDXResultReader can be used to perform a request to the wayback machine and
// fetch the snapshot data of a specific CDXResult.
type cdxResultReader struct {
	resp      *http.Response
	original  string
	timestamp time.Time
	eof       bool
}

// Read implements the Reader interface for CDXResultReader
func (dr *cdxResultReader) Read(p []byte) (int, error) {
	if dr.eof {
		return 0, io.EOF
	}
	if dr.resp == nil {
		client := http.Client{}
		req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s/%s", dataURL, dr.timestamp.Format("20060102150405"), dr.original), nil)
		if err != nil {
			return 0, err
		}
		fmt.Printf("%s/%s/%s\n", dataURL, dr.timestamp.Format("20060102150405"), dr.original)
		req.Header.Set("User-Agent", "simplewayback/0.1 (https://github.com/rhelmke/simplewayback)")
		req.Header.Del("Accept-Encoding")
		req.Header.Set("Accept", "*/*")
		dr.resp, err = client.Do(req)
		if err != nil {
			dr.eof = true
			return 0, err
		}
		if dr.resp.StatusCode != http.StatusOK {
			return 0, ErrorBadResponse
		}
	}
	return dr.resp.Body.Read(p)
}

// Perform queries the CDX API and returns a set of results
func (cdx *CDXAPI) Perform() ([]CDXResult, error) {
	isJSON := cdx.params.Get("output") == "json"
	// it's nice to have cdx and json support. But I don't think it's necessary
	// to implement parsing support for both output formats when this method
	// returns a []CDXResult-Type either ways. So we force json.
	if !isJSON {
		cdx.params.Set("output", "json")
	}
	qry, err := cdx.RawPerform()
	if err != nil {
		return []CDXResult{}, err
	}
	qryRes, err := ioutil.ReadAll(qry)
	if err != nil {
		return []CDXResult{}, err
	}
	splitBuf := [][]string{}
	if err := json.Unmarshal(qryRes, &splitBuf); err != nil {
		return []CDXResult{}, err
	}
	result := []CDXResult{}
	for i := 1; i < len(splitBuf); i++ {
		// resumption key stuff
		if len(splitBuf[i]) == 0 {
			continue
		}
		// resumption key stuff
		if len(splitBuf[i]) == 1 {
			if cdx.ResumptionKeyEnabled() {
				cdx.params.Set("resumeKey", splitBuf[i][0])
			}
			continue
		}
		// convert unknown status code to 0
		if splitBuf[i][4] == "-" {
			splitBuf[i][4] = "0"
		}
		if splitBuf[i][7] == "-" {
			splitBuf[i][7] == "0"
		}
		// parse time
		t, err := time.Parse("20060102150405", splitBuf[i][1])
		if err != nil {
			return []CDXResult{}, err
		}
		code, err := strconv.Atoi(splitBuf[i][4])
		if err != nil {
			return []CDXResult{}, err
		}
		ln, err := strconv.Atoi(splitBuf[i][6])
		if err != nil {
			return []CDXResult{}, err
		}
		result = append(result, CDXResult{URLKey: splitBuf[i][0], Timestamp: t, Original: splitBuf[i][2], MimeType: splitBuf[i][3], StatusCode: code, Digest: splitBuf[i][5], Length: ln, Data: &cdxResultReader{original: splitBuf[i][2], timestamp: t}})
	}
	// act as changing the output never happened :D
	if !isJSON {
		cdx.params.Del("output")
	}
	return result, nil
}

