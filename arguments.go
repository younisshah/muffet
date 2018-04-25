package main

import (
	"os"
	"strconv"

	"github.com/docopt/docopt-go"
)

const usage = `Muffet, the web repairgirl

Usage:
	muffet [-c <concurrency>] [-p <address>] [-v] <url>

Options:
	-c, --concurrency <concurrency>  Roughly maximum number of concurrent HTTP connections. [default: 512]
	-h, --help  Show this help.
	-p, --proxy-address <address>  Set an HTTP proxy address. (e.g. foo.com:4242)
	-v, --verbose  Show successful results too.`

type arguments struct {
	concurrency  int
	url          string
	proxyAddress string
	verbose      bool
}

func getArguments(ss []string) (arguments, error) {
	if ss == nil {
		ss = os.Args[1:]
	}

	args, err := docopt.ParseArgs(usage, ss, "0.1.0")

	if err != nil {
		return arguments{}, err
	}

	c, err := strconv.ParseInt(args["--concurrency"].(string), 10, 32)

	if err != nil {
		return arguments{}, err
	}

	p, _ := args["--proxy-address"].(string)

	return arguments{
		int(c),
		args["<url>"].(string),
		p,
		args["--verbose"].(bool),
	}, nil
}
