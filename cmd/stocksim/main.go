// Copyright (C) 2014 Adam Woodbeck.
// All rights reserved. Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
//
// Some bits of code are copyright (C) 2014 Jakob Borg and Contributors
// (see the Syncthing CONTRIBUTORS file).

package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"time"
)

const (
	usage      = "stocksim [options]"
	extraUsage = `The following enviroment variables are interpreted by stocksim:

  SSGUIADDRESS  Override the default listen address set in config. Expects protocol
                type followed by hostname or an IP address, followed by a port, such
                as "https://127.0.0.1:8888".

  SSGUIASSETS   Directory to load GUI assets from. Overrides compiled in assets.

  GOMAXPROCS    Set the maximum number of CPU cores to use. Defaults to all
                available CPU cores.`
)

var (
	Version       = "unknown-dev"
	BuildEnv      = "default"
	BuildStamp    = "0"
	BuildDate     time.Time
	BuildHost     = "unknown"
	BuildUser     = "unknown"
	LongVersion   string
	GoArchExtra   string // "", "v5", "v6", "v7"
	listenAddress string
	showVersion   bool
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
	flag.StringVar(&listenAddress, "listen-address", "http://127.0.0.1:8888", "Override the listening address.")
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
}
