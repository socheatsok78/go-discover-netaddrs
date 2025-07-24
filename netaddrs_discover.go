package netaddrs

import (
	"context"
	"fmt"
	"log"

	netaddrs "github.com/hashicorp/go-netaddrs"
)

type Provider struct{}

func (p *Provider) Help() string {
	return `netaddrs:
    Inspired by go-discover, go-netaddrs is a Go (golang) library and command line tool to discover ip addresses of
    nodes in a customizable fashion suitable for any environment. It returns IP addresses (IPv4 or IPv6) given a
      - DNS name, OR
      - custom executable with optional args which:
        - on success - exits with 0 and prints whitespace delimited IP (v4/v6) addresses to stdout.
        - on failure - exits with a non-zero code and optionally prints an error message of up to 1024 bytes to stderr.

    provider:         "netaddrs"
    exec:             "DNS name" or "executable with optional args".
`
}

func (p *Provider) Addrs(args map[string]string, l *log.Logger) ([]string, error) {
	if args["provider"] != "netaddrs" {
		return nil, fmt.Errorf("discover-netaddrs: invalid provider " + args["provider"])
	}

	if args["exec"] == "" {
		return nil, fmt.Errorf("discover-netaddrs: no exec provided")
	}

	addresses, err := netaddrs.IPAddrs(context.Background(), args["exec"], logger{l: l})
	if err != nil {
		return nil, fmt.Errorf("discover-netaddrs: %s", err)
	}

	var addrs []string

	for _, address := range addresses {
		outputAddress := address.IP.String()
		if address.Zone != "" {
			outputAddress += "%" + address.Zone
		}
		addrs = append(addrs, outputAddress)
	}

	return addrs, nil
}

type logger struct {
	l *log.Logger
}

func (l logger) Debug(msg string, args ...interface{}) {
	l.l.Print(msg)
	l.l.Println(args...)
}
