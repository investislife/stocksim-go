// Copyright (C) 2014 Adam Woodbeck.
// All rights reserved. Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
//
// Some bits of code (.e.g, embeddedStatic(), mimeTypeForFile(), etc.) are
// copyright (C) 2014 Jakob Borg and Contributors (see the Syncthing
// CONTRIBUTORS file).

package main

import (
	"flag"
	"fmt"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/w0uld/stocksim-go/auto"
)

var (
	Version      = "unknown-dev"
	BuildEnv     = "default"
	BuildStamp   = "0"
	BuildDate    time.Time
	BuildHost    = "unknown"
	BuildUser    = "unknown"
	LongVersion  string
	lastModified = time.Now().UTC().Format(http.TimeFormat)
	showVersion  bool
)

const (
	defaultAddress = "http://127.0.0.1:8888"
	usage          = "stocksim [options]"
	extraUsage     = `The following enviroment variables are interpreted by stocksim:

  SSGUIADDRESS  Override the default listen address set in config. Expects protocol
                type followed by hostname or an IP address, followed by a port, such
                as "` + defaultAddress + `".

  SSGUIASSETS   Directory to load GUI assets from. Overrides compiled in assets.

  GOMAXPROCS    Set the maximum number of CPU cores to use. Defaults to all
                available CPU cores.`
)

func init() {
	if Version != "unknown-dev" {
		// If not a generic dev build, version string should come from git describe
		exp := regexp.MustCompile(`^v\d+\.\d+\.\d+(-[a-z0-9]+)*(\+\d+-g[0-9a-f]+)?(-dirty)?$`)
		if !exp.MatchString(Version) {
			fmt.Sprintf("Invalid version string %q;\n\tdoes not match regexp %v", Version, exp)
		}
	}

	stamp, _ := strconv.Atoi(BuildStamp)
	BuildDate = time.Unix(int64(stamp), 0)

	date := BuildDate.UTC().Format("2006-01-02 15:04:05 MST")
	LongVersion = fmt.Sprintf("stocksim %s (%s %s-%s %s) %s@%s %s", Version, runtime.Version(), runtime.GOOS, runtime.GOARCH, BuildEnv, BuildUser, BuildHost, date)
}

func main() {
	var address string
	fullAddress := os.Getenv("STGUIADDRESS")

	if fullAddress == "" {
		fullAddress = defaultAddress
	}

	addressParts := strings.SplitN(fullAddress, "://", 2)
	switch addressParts[0] {
	case "http":
		address = addressParts[1]
	case "https":
		log.Fatal("HTTPS is not yet supported.")
	default:
		log.Fatal("Unidentified protocol", addressParts[0])
	}

	flag.BoolVar(&showVersion, "version", false, "Show version")
	flag.Usage = usageFor(flag.CommandLine, usage, extraUsage)
	flag.Parse()

	if showVersion {
		fmt.Println(LongVersion)
		return
	}

	if len(os.Getenv("GOMAXPROCS")) == 0 {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}

	// Main handler
	mux := http.NewServeMux()

	// Serve compiled in assets unless an asset directory was set (for development)
	mux.Handle("/", embeddedStatic(os.Getenv("SSGUIASSETS")))

	fmt.Printf("Listening on %s ...\n", address)
	http.ListenAndServe(address, mux)
}

func embeddedStatic(assetDir string) http.Handler {
	assets := auto.Assets()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		file := r.URL.Path

		if file[0] == '/' {
			file = file[1:]
		}

		if len(file) == 0 {
			file = "index.html"
		}

		if assetDir != "" {
			p := filepath.Join(assetDir, filepath.FromSlash(file))
			_, err := os.Stat(p)
			if err == nil {
				http.ServeFile(w, r, p)
				return
			}
		}

		bs, ok := assets[file]
		if !ok {
			http.NotFound(w, r)
			return
		}

		contentType := mimeTypeForFile(file)
		if len(contentType) != 0 {
			w.Header().Set("Content-Type", contentType)
		}
		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(bs)))
		w.Header().Set("Last-Modified", lastModified)

		w.Write(bs)
	})
}

func mimeTypeForFile(file string) string {
	// We use a built in table of the common types since the system
	// TypeByExtension might be unreliable. But if we don't know, we delegate
	// to the system.
	ext := filepath.Ext(file)
	switch ext {
	case ".htm", ".html":
		return "text/html"
	case ".css":
		return "text/css"
	case ".js":
		return "application/javascript"
	case ".json":
		return "application/json"
	case ".png":
		return "image/png"
	case ".ttf":
		return "application/x-font-ttf"
	case ".woff":
		return "application/x-font-woff"
	default:
		return mime.TypeByExtension(ext)
	}
}
