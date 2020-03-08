package main

import (
	"github.com/usvc/go-config"
)

var (
	conf = config.Map{
		"server_addr": &config.String{
			Default: "0.0.0.0",
			Usage:   "ip address/hostname to bind to",
		},
		"server_port": &config.Uint{
			Default: 8888,
			Usage:   "port to listen on",
		},
	}
)
