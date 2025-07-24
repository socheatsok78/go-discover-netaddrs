package netaddrs

import (
	"context"
	"fmt"
	"log"

	netaddrs "github.com/hashicorp/go-netaddrs"
)

type Provider struct{}

func (p *Provider) Help() string {
	return `go-netaddrs:

    provider:         "netaddrs"
    args:             "DNS name" or "executable with optional args".
`
}

func (p *Provider) Addrs(args map[string]string, l *log.Logger) ([]string, error) {
	if args["provider"] != "netaddrs" {
		return nil, fmt.Errorf("discover-netaddrs: invalid provider " + args["provider"])
	}

	if args["args"] == "" {
		return nil, fmt.Errorf("discover-netaddrs: no args provided")
	}

	addresses, err := netaddrs.IPAddrs(context.Background(), args["args"], logger{l: l})
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
