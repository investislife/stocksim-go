#!/usr/bin/env bash
set -euo pipefail
IFS=$'\n\t'

case "${1:-default}" in
	default)
		# go run build.go "assets"
		# go run build.go "deps"
		go run build.go
		;;

	clean)
		go run build.go "$1"
		;;

	test)
		ulimit -t 60 &>/dev/null || true
		ulimit -d 512000 &>/dev/null || true
		ulimit -m 512000 &>/dev/null || true

		go run build.go "$1"
		;;

	tar)
		go run build.go "$1"
		;;

	deps)
		go run build.go "$1"
		;;

	assets)
		go run build.go "$1"
		;;

	all)
		go run build.go -goos linux -goarch amd64 tar
		go run build.go -goos linux -goarch 386 tar
		go run build.go -goos linux -goarch armv5 tar
		go run build.go -goos linux -goarch armv6 tar
		go run build.go -goos linux -goarch armv7 tar

		go run build.go -goos freebsd -goarch amd64 tar
		go run build.go -goos freebsd -goarch 386 tar

		go run build.go -goos darwin -goarch amd64 tar

		go run build.go -goos windows -goarch amd64 zip
		go run build.go -goos windows -goarch 386 zip
		;;

	setup)
		echo "Don't worry, just build."
		;;

	*)
		echo "Unknown build command $1"
		;;
esac
