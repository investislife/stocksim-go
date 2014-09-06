// Copyright (C) 2014 Adam Woodbeck.
// All rights reserved. Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
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
	listenAddress string
)

func main() {
	flag.StringVar(&listenAddress, "listen-address", "http://127.0.0.1:8888", "Override the listening address.")
	flag.Usage = usageFor(flag.CommandLine, usage, extraUsage)
	flag.Parse()

}
