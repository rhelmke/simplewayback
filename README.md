# simplewayback

[![GoDoc](https://godoc.org/github.com/rhelmke/simplewayback?status.svg)](https://godoc.org/github.com/rhelmke/simplewayback)
[![Build Status](https://travis-ci.org/rhelmke/simplewayback.svg?branch=master)](https://travis-ci.org/rhelmke/simplewayback)
[![Go Report Card](https://goreportcard.com/badge/github.com/rhelmke/simplewayback)](https://goreportcard.com/report/github.com/rhelmke/simplewayback)

**Please note that this package was build to suit my needs in another project I am working on. Therefore, further development of this package is highly dependent on the requirements needed by said project. But feel free to create feature requests, bug reports or pull requests! We'll see what we can do.** 

## Description

simplewayback is a simple go package for querying the Wayback Machine CDX API and fetching snapshots.

## Get

```bash
go get -u github.com/rhelmke/simplewayback
```

## Searching the CDX API

The simplest way to work with simplewayback is fetching all results at once:

```go
package main

import (
    "fmt"
    wayback "github.com/rhelmke/simplewayback"
)

func main() {
    // Create new API wrapper for URL "archive.org"
    cdx, err := wayback.NewCDXAPI("archive.org")
    if err != nil {
        fmt.Println(err)
        return
    }
    // this will return []CDXResult
    results, err := cdx.Perform()
    if err != nil {
        fmt.Println(err)
        return
    }

    for _, result := range results {
        fmt.Printf("-----------\nURLKey: %s\nTimestamp: %s\nOriginal: %s\nMimetype: %s\nStatuscode: %d\nDigest: %s\nLength: %d\n",
            result.URLKey,
            result.Timestamp.String(),
            result.Original,
            result.MimeType,
            result.StatusCode,
            result.Digest,
            result.Length,
        )
    }
}
```

Each `CDXResult` represents a snapshot taken by the Wayback Machine. Taking it a step further, we want to fetch the actual snapshot data from a specific CDX Result. We can do so by accessing the `Data`-Attribute of `CDXResult`. `Data` implements the [io.Reader](https://golang.org/pkg/io/#Reader)-Interface and will perform a query to the Wayback API fetching the snapshot of that specific result:

```go
package main

import (
    "bytes"
    "fmt"
    "io/ioutil"
    wayback "github.com/rhelmke/simplewayback"
)

func main() {
    // Create new API wrapper for URL "archive.org"
    cdx, err := wayback.NewCDXAPI("archive.org/robots.txt")
    // Fetch only 3 items
    cdx.SetLimit(3)
    // Use the collapsing filter to only get unique results
    cdx.AddCollapsing(wayback.FieldDigest, 0)
    if err != nil {
        fmt.Println(err)
        return
    }
    // this will return []CDXResult
    results, err := cdx.Perform()
    if err != nil {
        fmt.Println(err)
        return
    }

    // Buffer for printing the results
    var buf bytes.Buffer

    // Iterate all CDXResults, download snapshots and print them
    for _, result := range results {
        buf.Reset()
        snapshot, err := ioutil.ReadAll(result.Data)
        if err != nil {
            fmt.Println(err)
            continue
        }
        buf.WriteString("-----------\n")
        buf.WriteString("Digest: ")
        buf.WriteString(result.Digest)
        buf.WriteString("\nTimestamp: ")
        buf.WriteString(result.Timestamp.String())
        buf.WriteString("\nSnapshot:\n\n\n")
        buf.Write(snapshot)
        fmt.Println(buf.String())
    }
}
```

You might have noticed that you can instruct `simplewayback` to use some of the advanced filters like `collapsing`. For a full set of supported features conduct [documentation](https://godoc.org/github.com/rhelmke/simplewayback) and [CDX API](https://github.com/internetarchive/wayback/tree/master/wayback-cdx-server).

## Raw CDX Results
In case you want to build your own Parser for CDX Search Results, you can do so by invoking `cdx.RawPerform()` instead of `cdx.Perform()`. `RawPerform()` will return an [io.Reader](https://golang.org/pkg/io/#Reader) to query the Wayback Machine:

```go
package main

import (
    "fmt"
    "io/ioutil"
    wayback "github.com/rhelmke/simplewayback"
)

func main() {
    // Create new API wrapper for URL "archive.org"
    cdx, err := wayback.NewCDXAPI("archive.org/robots.txt")
    // Fetch only 3 items
    cdx.SetLimit(3)
    // Use the collapsing filter to only get unique results
    cdx.AddCollapsing(wayback.FieldDigest, 0)
    if err != nil {
        fmt.Println(err)
        return
    }
    // Create a raw CDX reader
    rawReader, err := cdx.RawPerform()
    if err != nil {
        fmt.Println(err)
        return
    }

    result, err := ioutil.ReadAll(rawReader)
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Println(string(result))
}

```
