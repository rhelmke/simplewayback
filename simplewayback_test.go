package simplewayback

import (
	"bytes"
	neturl "net/url"
	"reflect"
	"testing"
	"time"
)

func TestNewCDXAPI(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name    string
		args    args
		want    *CDXAPI
		wantErr bool
	}{
		{"Trigger ErrorInvalidScheme", args{url: "ftp://archive.org"}, nil, true},
		{"Create CDXAPI", args{url: "archive.org"}, &CDXAPI{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewCDXAPI(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewCDXAPI() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if reflect.TypeOf(got).Kind() != reflect.TypeOf(tt.want).Kind() {
				t.Errorf("NewCDXAPI() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCDXAPI_SetAPIKey(t *testing.T) {
	type args struct {
		apiKey string
	}
	tests := []struct {
		name    string
		cdx     *CDXAPI
		args    args
		wantErr bool
	}{
		{"SetAPIKey", &CDXAPI{params: &neturl.Values{}, urlBuf: &bytes.Buffer{}}, args{apiKey: "testkey"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.cdx.SetAPIKey(tt.args.apiKey); (err != nil) != tt.wantErr {
				t.Errorf("CDXAPI.SetAPIKey() error = %v, wantErr %v", err, tt.wantErr)
			} else if tt.cdx.apiKey != tt.args.apiKey {
				t.Errorf("CDXAPI.SetAPIKey() did not set the correct apiKey = %v, want %v", tt.cdx.apiKey, tt.args.apiKey)
			}
		})
	}
}

func TestCDXAPI_APIKey(t *testing.T) {
	cdx := &CDXAPI{params: &neturl.Values{}, urlBuf: &bytes.Buffer{}, apiKey: "testkey"}
	tests := []struct {
		name string
		cdx  *CDXAPI
		want string
	}{
		{"Getter", cdx, "testkey"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cdx.APIKey(); got != tt.want {
				t.Errorf("CDXAPI.APIKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCDXAPI_SetMatchType(t *testing.T) {
	type args struct {
		mType matchType
	}
	cdx := &CDXAPI{params: &neturl.Values{}, urlBuf: &bytes.Buffer{}}
	tests := []struct {
		name    string
		cdx     *CDXAPI
		args    args
		wantErr bool
	}{
		{"InvalidMatchType", cdx, args{mType: -1}, true},
		{"Valid Match Type", cdx, args{mType: MatchTypeDomain}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.cdx.SetMatchType(tt.args.mType); (err != nil) != tt.wantErr {
				t.Errorf("CDXAPI.SetMatchType() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCDXAPI_MatchType(t *testing.T) {
	cdx1 := &CDXAPI{params: &neturl.Values{}, urlBuf: &bytes.Buffer{}}
	cdx1.SetMatchType(MatchTypeDomain)
	cdx2 := &CDXAPI{params: &neturl.Values{}, urlBuf: &bytes.Buffer{}}
	cdx2.SetMatchType(MatchTypeExact)
	cdx3 := &CDXAPI{params: &neturl.Values{}, urlBuf: &bytes.Buffer{}}
	cdx3.SetMatchType(MatchTypePrefix)
	cdx4 := &CDXAPI{params: &neturl.Values{}, urlBuf: &bytes.Buffer{}}
	cdx4.SetMatchType(MatchTypeHost)
	tests := []struct {
		name string
		cdx  *CDXAPI
		want int
	}{
		{"Domain", cdx1, int(MatchTypeDomain)},
		{"Exact", cdx2, int(MatchTypeExact)},
		{"Prefix", cdx3, int(MatchTypePrefix)},
		{"Host", cdx4, int(MatchTypeHost)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cdx.MatchType(); got != tt.want {
				t.Errorf("CDXAPI.MatchType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCDXAPI_ResetMatchType(t *testing.T) {
	cdx := &CDXAPI{params: &neturl.Values{}, urlBuf: &bytes.Buffer{}}
	cdx.SetMatchType(MatchTypeHost)
	tests := []struct {
		name string
		cdx  *CDXAPI
	}{
		{"Default reset", cdx},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.cdx.ResetMatchType()
			if tt.cdx.params.Get("matchType") != "exact" {
				t.Errorf("CDXAPI.ResetMatchType() = %v, want exact", tt.cdx.params.Get("matchType"))
			}
		})
	}
}

func TestCDXAPI_SetOutputFormat(t *testing.T) {
	cdx := &CDXAPI{params: &neturl.Values{}, urlBuf: &bytes.Buffer{}}
	type args struct {
		format outputFormat
	}
	tests := []struct {
		name    string
		cdx     *CDXAPI
		args    args
		wantErr bool
	}{
		{"ErrorInvalidOutputFormat", cdx, args{format: -1}, true},
		{"CDX", cdx, args{format: OutputFormatCDX}, false},
		{"JSON", cdx, args{format: OutputFormatJSON}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.cdx.SetOutputFormat(tt.args.format); (err != nil) != tt.wantErr {
				t.Errorf("CDXAPI.SetOutputFormat() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCDXAPI_OutputFormat(t *testing.T) {
	cdx1 := &CDXAPI{params: &neturl.Values{}, urlBuf: &bytes.Buffer{}}
	cdx1.SetOutputFormat(OutputFormatJSON)
	cdx2 := &CDXAPI{params: &neturl.Values{}, urlBuf: &bytes.Buffer{}}
	tests := []struct {
		name string
		cdx  *CDXAPI
		want int
	}{
		{"OutputFormatJSON", cdx1, int(OutputFormatJSON)},
		{"OutputFormatCDX", cdx2, int(OutputFormatCDX)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cdx.OutputFormat(); got != tt.want {
				t.Errorf("CDXAPI.OutputFormat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCDXAPI_ResetOutputFormat(t *testing.T) {
	tests := []struct {
		name string
		cdx  *CDXAPI
	}{
		{"Default Reset", &CDXAPI{params: &neturl.Values{}, urlBuf: &bytes.Buffer{}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.cdx.SetOutputFormat(OutputFormatJSON)
			tt.cdx.ResetOutputFormat()
			if tt.cdx.params.Get("output") != "" {
				t.Errorf("CDXAPI.ResetOutputFormat() resets param to %v, want ''", tt.cdx.params.Get("output"))
			}
		})
	}
}

func TestCDXAPI_SetURL(t *testing.T) {
	cdx := &CDXAPI{params: &neturl.Values{}, urlBuf: &bytes.Buffer{}}
	type args struct {
		url string
	}
	tests := []struct {
		name    string
		cdx     *CDXAPI
		args    args
		wantErr bool
	}{
		{"ErrorInvalidScheme", cdx, args{"ftp://archive.org"}, true},
		{"HTTP", cdx, args{"hTtP://archive.org"}, false},
		{"No Scheme", cdx, args{"archive.org"}, false},
		{"HTTPS", cdx, args{"hTtPs://archive.org"}, false},
		{"url.Parse", cdx, args{"ü>äasdläüö:\\\\archive;org"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.cdx.SetURL(tt.args.url); (err != nil) != tt.wantErr {
				t.Errorf("CDXAPI.SetURL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCDXAPI_URL(t *testing.T) {
	cdx := &CDXAPI{params: &neturl.Values{}, urlBuf: &bytes.Buffer{}}
	cdx.SetURL("https://archive.org")
	tests := []struct {
		name string
		cdx  *CDXAPI
		want string
	}{
		{"Getter", cdx, "https://archive.org"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cdx.URL(); got != tt.want {
				t.Errorf("CDXAPI.URL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCDXAPI_SetLimit(t *testing.T) {
	cdx := &CDXAPI{params: &neturl.Values{}, urlBuf: &bytes.Buffer{}}
	type args struct {
		limit int
	}
	tests := []struct {
		name    string
		cdx     *CDXAPI
		args    args
		wantErr bool
	}{
		{"Invalid SetLimit()", cdx, args{-1}, true},
		{"Valid SetLimit()", cdx, args{10}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.cdx.SetLimit(tt.args.limit); (err != nil) != tt.wantErr {
				t.Errorf("CDXAPI.SetLimit() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCDXAPI_Limit(t *testing.T) {
	cdx1 := &CDXAPI{params: &neturl.Values{}, urlBuf: &bytes.Buffer{}}
	cdx2 := &CDXAPI{params: &neturl.Values{}, urlBuf: &bytes.Buffer{}}
	cdx2.SetLimit(10)
	cdx3 := &CDXAPI{params: &neturl.Values{}, urlBuf: &bytes.Buffer{}}
	cdx3.params.Add("limit", "ahvsidasd")
	tests := []struct {
		name string
		cdx  *CDXAPI
		want int
	}{
		{"No Limit set", cdx1, -1},
		{"Correct Limit", cdx2, 10},
		{"Limit error", cdx3, -1}, // this should never happen
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cdx.Limit(); got != tt.want {
				t.Errorf("CDXAPI.Limit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCDXAPI_ResetLimit(t *testing.T) {
	cdx := &CDXAPI{params: &neturl.Values{}, urlBuf: &bytes.Buffer{}}
	cdx.SetLimit(10)
	tests := []struct {
		name string
		cdx  *CDXAPI
	}{
		{"Default Reset", cdx},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.cdx.ResetLimit()
			if tt.cdx.params.Get("") != "" {
				t.Errorf("CDXAPI.ResetLimit() resets param to %v, want ''", tt.cdx.params.Get(""))
			}
		})
	}
}

func TestCDXAPI_AddRegexFilter(t *testing.T) {
	cdx := &CDXAPI{params: &neturl.Values{}, urlBuf: &bytes.Buffer{}}
	type args struct {
		fld    field
		regex  string
		negate bool
	}
	tests := []struct {
		name    string
		cdx     *CDXAPI
		args    args
		wantErr bool
	}{
		{"ErrorInvalidField", cdx, args{fld: -1, regex: ".*", negate: false}, true},
		{"CompileError", cdx, args{fld: FieldDigest, regex: "(?.*", negate: false}, true},
		{"NoErrorNonNegate", cdx, args{fld: FieldDigest, regex: "XYA", negate: false}, false},
		{"NoErrorNegate", cdx, args{fld: FieldDigest, regex: "XYA", negate: true}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.cdx.AddRegexFilter(tt.args.fld, tt.args.regex, tt.args.negate); (err != nil) != tt.wantErr {
				t.Errorf("CDXAPI.AddRegexFilter() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCDXAPI_ResetRegexFilters(t *testing.T) {
	cdx := &CDXAPI{params: &neturl.Values{}, urlBuf: &bytes.Buffer{}}
	cdx.AddRegexFilter(FieldDigest, "XYA", false)
	cdx.AddRegexFilter(FieldStatuscode, "200", true)
	tests := []struct {
		name string
		cdx  *CDXAPI
	}{
		{"Default Reset", cdx},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.cdx.ResetRegexFilters()
			if len(tt.cdx.regFilterKeys) != 0 {
				t.Errorf("CDXAPI.ResetRegexFilters() didn't reset the filter keys")
			}
		})
	}
}

func TestCDXAPI_SetTimeFilter(t *testing.T) {
	cdx := &CDXAPI{params: &neturl.Values{}, urlBuf: &bytes.Buffer{}}
	type args struct {
		from time.Time
		to   time.Time
	}
	tests := []struct {
		name    string
		cdx     *CDXAPI
		args    args
		wantErr bool
	}{
		{"ErrorInvalidFromTo", cdx, args{time.Now().Add(3 * time.Hour), time.Now()}, true},
		{"ValidFromTo", cdx, args{time.Now(), time.Now().Add(3 * time.Hour)}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.cdx.SetTimeFilter(tt.args.from, tt.args.to); (err != nil) != tt.wantErr {
				t.Errorf("CDXAPI.SetTimeFilter() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCDXAPI_TimeFilter(t *testing.T) {
	cdx1 := &CDXAPI{params: &neturl.Values{}, urlBuf: &bytes.Buffer{}}
	cdx1.params.Add("from", "asuzdtuaivsd")
	cdx1.params.Add("to", "20060102150405")
	cdx2 := &CDXAPI{params: &neturl.Values{}, urlBuf: &bytes.Buffer{}}
	cdx2.params.Add("from", "20060102150405")
	cdx2.params.Add("to", "asuzdtuaivsd")
	cdx3 := &CDXAPI{params: &neturl.Values{}, urlBuf: &bytes.Buffer{}}
	tm, _ := time.Parse("20060102150405", "20060102150405")
	tm2, _ := time.Parse("20060102150405", "20070102150405")
	cdx3.SetTimeFilter(tm, tm2)
	tests := []struct {
		name  string
		cdx   *CDXAPI
		want  time.Time
		want1 time.Time
	}{
		{"Invalid Time1", cdx1, time.Time{}, time.Time{}},
		{"Invalid Time2", cdx2, time.Time{}, time.Time{}},
		{"Valid Time", cdx3, tm, tm2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.cdx.TimeFilter()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CDXAPI.TimeFilter() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("CDXAPI.TimeFilter() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestCDXAPI_ResetTimeFilter(t *testing.T) {
	cdx := &CDXAPI{params: &neturl.Values{}, urlBuf: &bytes.Buffer{}}
	cdx.SetTimeFilter(time.Now(), time.Now().Add(3*time.Hour))
	tests := []struct {
		name string
		cdx  *CDXAPI
	}{
		{"Default Reset", cdx},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.cdx.ResetTimeFilter()
			if tt.cdx.params.Get("from") != "" || tt.cdx.params.Get("to") != "" {
				t.Errorf("CDXAPI.ResetTimeFilter() didn't reset the underlying values")
			}
		})
	}
}

func TestCDXAPI_AddCollapsing(t *testing.T) {
	cdx := &CDXAPI{params: &neturl.Values{}, urlBuf: &bytes.Buffer{}}
	type args struct {
		fld field
		n   int
	}
	tests := []struct {
		name    string
		cdx     *CDXAPI
		args    args
		wantErr bool
	}{
		{"ErrorInvalidField", cdx, args{-1, 1}, true},
		{"ErrorInvalidNumber", cdx, args{FieldDigest, -1}, true},
		{"No N", cdx, args{FieldDigest, 0}, false},
		{"N", cdx, args{FieldDigest, 10}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.cdx.AddCollapsing(tt.args.fld, tt.args.n); (err != nil) != tt.wantErr {
				t.Errorf("CDXAPI.AddCollapsing() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCDXAPI_ResetCollapsing(t *testing.T) {
	cdx := &CDXAPI{params: &neturl.Values{}, urlBuf: &bytes.Buffer{}}
	cdx.AddCollapsing(FieldDigest, 10)
	cdx.AddCollapsing(FieldDigest, 1)
	tests := []struct {
		name string
		cdx  *CDXAPI
	}{
		{"Default Reset", cdx},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.cdx.ResetCollapsing()
			if len(tt.cdx.collapsingKeys) != 0 {
				t.Errorf("CDXAPI.ResetCollapsing() didn't reset the collapsing keys")
			}
		})
	}
}

func TestCDXAPI_SetGzip(t *testing.T) {
	cdx := &CDXAPI{params: &neturl.Values{}, urlBuf: &bytes.Buffer{}}
	type args struct {
		enabled bool
	}
	tests := []struct {
		name    string
		cdx     *CDXAPI
		args    args
		wantErr bool
	}{
		{"Disable", cdx, args{false}, false},
		{"Enable", cdx, args{true}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.cdx.SetGzip(tt.args.enabled); (err != nil) != tt.wantErr {
				t.Errorf("CDXAPI.SetGzip() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCDXAPI_Gzip(t *testing.T) {
	cdx1 := &CDXAPI{params: &neturl.Values{}, urlBuf: &bytes.Buffer{}}
	cdx1.SetGzip(false)
	cdx2 := &CDXAPI{params: &neturl.Values{}, urlBuf: &bytes.Buffer{}}
	cdx2.SetGzip(true)
	tests := []struct {
		name string
		cdx  *CDXAPI
		want bool
	}{
		{"Disabled", cdx1, false},
		{"Enabled", cdx2, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cdx.Gzip(); got != tt.want {
				t.Errorf("CDXAPI.Gzip() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCDXAPI_ResetGzip(t *testing.T) {
	cdx := &CDXAPI{params: &neturl.Values{}, urlBuf: &bytes.Buffer{}}
	cdx.SetGzip(false)
	tests := []struct {
		name string
		cdx  *CDXAPI
	}{
		{"Default Reset", cdx},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.cdx.ResetGzip()
			if !cdx.Gzip() {
				t.Errorf("CDXAPI.ResetGzip() didn't reset the gzip flag")
			}
		})
	}
}

func TestCDXAPI_SetOffset(t *testing.T) {
	cdx := &CDXAPI{params: &neturl.Values{}, urlBuf: &bytes.Buffer{}}
	type args struct {
		offset int
	}
	tests := []struct {
		name    string
		cdx     *CDXAPI
		args    args
		wantErr bool
	}{
		{"ErrorInvalidNumber", cdx, args{-1}, true},
		{"ValidNumber", cdx, args{1}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.cdx.SetOffset(tt.args.offset); (err != nil) != tt.wantErr {
				t.Errorf("CDXAPI.SetOffset() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCDXAPI_Offset(t *testing.T) {
	cdx1 := &CDXAPI{params: &neturl.Values{}, urlBuf: &bytes.Buffer{}}
	cdx2 := &CDXAPI{params: &neturl.Values{}, urlBuf: &bytes.Buffer{}}
	cdx2.params.Set("offset", "asd")
	cdx3 := &CDXAPI{params: &neturl.Values{}, urlBuf: &bytes.Buffer{}}
	cdx3.SetOffset(10)
	tests := []struct {
		name string
		cdx  *CDXAPI
		want int
	}{
		{"No offset", cdx1, -1},
		{"Conversion error", cdx2, -1}, // should never happen
		{"Valid offset", cdx3, 10},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cdx.Offset(); got != tt.want {
				t.Errorf("CDXAPI.Offset() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCDXAPI_ResetOffset(t *testing.T) {
	cdx := &CDXAPI{params: &neturl.Values{}, urlBuf: &bytes.Buffer{}}
	cdx.SetOffset(10)
	tests := []struct {
		name string
		cdx  *CDXAPI
	}{
		{"Default Reset", cdx},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.cdx.ResetOffset()
			if cdx.Offset() != -1 {
				t.Errorf("CDXAPI.ResetGzip() didn't reset the offset")
			}
		})
	}
}

func TestCDXAPI_SetResumptionKey(t *testing.T) {
	cdx1 := &CDXAPI{params: &neturl.Values{}, urlBuf: &bytes.Buffer{}}
	cdx1.usePagination = true
	cdx2 := &CDXAPI{params: &neturl.Values{}, urlBuf: &bytes.Buffer{}}
	type args struct {
		enabled bool
		key     string
	}
	tests := []struct {
		name    string
		cdx     *CDXAPI
		args    args
		wantErr bool
	}{
		{"ErrorPaginationResumption", cdx1, args{true, "key"}, true},
		{"Enabled", cdx2, args{true, "key"}, false},
		{"Disabled", cdx2, args{false, ""}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.cdx.SetResumptionKey(tt.args.enabled, tt.args.key); (err != nil) != tt.wantErr {
				t.Errorf("CDXAPI.SetResumptionKey() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCDXAPI_ResumptionKeyEnabled(t *testing.T) {
	cdx := &CDXAPI{params: &neturl.Values{}, urlBuf: &bytes.Buffer{}}
	tests := []struct {
		name string
		cdx  *CDXAPI
		want bool
	}{
		{"Enabled", cdx, true},
		{"Disabled", cdx, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cdx.SetResumptionKey(tt.want, "key")
			if got := tt.cdx.ResumptionKeyEnabled(); got != tt.want {
				t.Errorf("CDXAPI.ResumptionKeyEnabled() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCDXAPI_ResumptionKey(t *testing.T) {
	cdx := &CDXAPI{params: &neturl.Values{}, urlBuf: &bytes.Buffer{}}
	tests := []struct {
		name string
		cdx  *CDXAPI
		want string
	}{
		{"Getter", cdx, "key"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cdx.SetResumptionKey(true, tt.want)
			if got := tt.cdx.ResumptionKey(); got != tt.want {
				t.Errorf("CDXAPI.ResumptionKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCDXAPI_ResetResumptionKey(t *testing.T) {
	cdx := &CDXAPI{params: &neturl.Values{}, urlBuf: &bytes.Buffer{}}
	cdx.SetResumptionKey(true, "key")
	tests := []struct {
		name string
		cdx  *CDXAPI
	}{
		{"Default Reset", cdx},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.cdx.ResetResumptionKey()
			if tt.cdx.useResumptionKey != false || tt.cdx.params.Get("resumeKey") != "" {
				t.Errorf("CDXAPI.ResetResumptionKey() didn't reset the resumption key and usage")
			}
		})
	}
}

func TestCDXAPI_SetPagination(t *testing.T) {
	cdx1 := &CDXAPI{params: &neturl.Values{}, urlBuf: &bytes.Buffer{}}
	cdx1.useResumptionKey = true
	cdx2 := &CDXAPI{params: &neturl.Values{}, urlBuf: &bytes.Buffer{}}
	type args struct {
		enabled bool
		page    int
	}
	tests := []struct {
		name    string
		cdx     *CDXAPI
		args    args
		wantErr bool
	}{
		{"ErrorPaginationResumption", cdx1, args{true, 1}, true},
		{"ErrorInvalidNumber", cdx2, args{true, -1}, true},
		{"Enabled", cdx2, args{true, 1}, false},
		{"Disabled", cdx2, args{false, 1}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.cdx.SetPagination(tt.args.enabled, tt.args.page); (err != nil) != tt.wantErr {
				t.Errorf("CDXAPI.SetPagination() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCDXAPI_PaginationEnabled(t *testing.T) {
	cdx1 := &CDXAPI{params: &neturl.Values{}, urlBuf: &bytes.Buffer{}}
	cdx1.SetPagination(true, 3)
	cdx2 := &CDXAPI{params: &neturl.Values{}, urlBuf: &bytes.Buffer{}}
	tests := []struct {
		name string
		cdx  *CDXAPI
		want bool
	}{
		{"Enabled", cdx1, true},
		{"Disabled", cdx2, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cdx.PaginationEnabled(); got != tt.want {
				t.Errorf("CDXAPI.PaginationEnabled() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCDXAPI_PaginationPage(t *testing.T) {
	cdx := &CDXAPI{params: &neturl.Values{}, urlBuf: &bytes.Buffer{}}
	cdx.SetPagination(true, 3)
	tests := []struct {
		name string
		cdx  *CDXAPI
		want int
	}{
		{"Getter", cdx, 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cdx.PaginationPage(); got != tt.want {
				t.Errorf("CDXAPI.PaginationPage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCDXAPI_ResetPagination(t *testing.T) {
	cdx := &CDXAPI{params: &neturl.Values{}, urlBuf: &bytes.Buffer{}}
	cdx.SetPagination(true, 3)
	tests := []struct {
		name string
		cdx  *CDXAPI
	}{
		{"Default Reset", cdx},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.cdx.ResetPagination()
			if tt.cdx.params.Get("page") != "" || tt.cdx.page != -1 || tt.cdx.usePagination {
				t.Errorf("CDXAPI.ResetPagination() didn't reset the pagination info")
			}
		})
	}
}

func TestCDXAPI_buildURL(t *testing.T) {
	cdx1 := &CDXAPI{params: &neturl.Values{}, urlBuf: &bytes.Buffer{}}
	cdx1.SetURL("archive.org")
	cdx1.AddCollapsing(FieldLength, 10)
	cdx1.AddCollapsing(FieldStatuscode, 0)
	cdx1.AddRegexFilter(FieldMimetype, "text/html", true)
	cdx1.SetOffset(1)
	cdx1.SetLimit(1)
	cdx1.SetMatchType(MatchTypeHost)
	tm, _ := time.Parse("20060102150405", "20060102150405")
	tm2, _ := time.Parse("20060102150405", "20070102150405")
	cdx1.SetTimeFilter(tm, tm2)
	cdx1.SetOutputFormat(OutputFormatJSON)
	cdx2 := &CDXAPI{params: &neturl.Values{}, urlBuf: &bytes.Buffer{}}
	type args struct {
		urlDst *bytes.Buffer
	}
	tests := []struct {
		name    string
		cdx     *CDXAPI
		args    args
		wantErr bool
	}{
		{"buildURL", cdx1, args{&bytes.Buffer{}}, false},
		{"ErrorInvalidURL", cdx2, args{&bytes.Buffer{}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			want := "https://web.archive.org/cdx/search/cdx?collapse=length%3A10&collapse=statuscode&filter=%21mimetype%3Atext%2Fhtml&from=20060102150405&limit=1&matchType=host&offset=1&output=json&to=20070102150405&url=archive.org"
			if err := tt.cdx.buildURL(tt.args.urlDst); (err != nil) != tt.wantErr {
				t.Errorf("CDXAPI.buildURL() error = %v, wantErr %v", err, tt.wantErr)
			} else if tt.args.urlDst.String() != "" && tt.args.urlDst.String() != want {
				t.Errorf("CDXAPI.buildURL() built URL = %v, want = %v", tt.args.urlDst.String(), want)
			}
		})
	}
}

func TestCDXRawQuery_Read(t *testing.T) {
	cdx1 := &CDXAPI{params: &neturl.Values{}, urlBuf: &bytes.Buffer{}}
	cdx1.SetURL("archive.org")
	cdx1.SetLimit(1)
	qry, _ := cdx1.RawPerform()
	type args struct {
		p []byte
	}
	tests := []struct {
		name    string
		qry     *CDXRawQuery
		args    args
		want    int
		wantErr bool
	}{
		{"read2", qry, args{make([]byte, 2)}, 2, false},
		{"readToEOF", qry, args{make([]byte, 10000)}, 106, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.qry.Read(tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("CDXRawQuery.Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CDXRawQuery.Read() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCDXAPI_RawPerform(t *testing.T) {
	cdx1 := &CDXAPI{params: &neturl.Values{}, urlBuf: &bytes.Buffer{}}
	cdx2 := &CDXAPI{params: &neturl.Values{}, urlBuf: &bytes.Buffer{}}
	cdx2.SetURL("archive.org")
	cdx2.SetLimit(1)
	tests := []struct {
		name    string
		cdx     *CDXAPI
		want    *CDXRawQuery
		wantErr bool
	}{
		{"ErrorInvalidURL", cdx1, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.cdx.RawPerform()
			if (err != nil) != tt.wantErr {
				t.Errorf("CDXAPI.RawPerform() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CDXAPI.RawPerform() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cdxResultReader_Read(t *testing.T) {
	type args struct {
		p []byte
	}
	tests := []struct {
		name    string
		dr      *cdxResultReader
		args    args
		want    int
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.dr.Read(tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("cdxResultReader.Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("cdxResultReader.Read() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCDXAPI_Perform(t *testing.T) {
	tests := []struct {
		name    string
		cdx     *CDXAPI
		want    []CDXResult
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.cdx.Perform()
			if (err != nil) != tt.wantErr {
				t.Errorf("CDXAPI.Perform() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CDXAPI.Perform() = %v, want %v", got, tt.want)
			}
		})
	}
}
