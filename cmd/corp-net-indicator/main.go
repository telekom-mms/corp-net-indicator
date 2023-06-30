package main

import (
	"flag"

	"github.com/telekom-mms/corp-net-indicator/internal/logger"
	"github.com/telekom-mms/corp-net-indicator/internal/ui"
)

var quickConnect bool
var verbose bool

func init() {
	flag.BoolVar(&quickConnect, "quick", false, "quick connect to vpn")
	flag.BoolVar(&verbose, "v", false, "verbose logging")
}

// starts indicator as window
func main() {
	flag.Parse()

	// start as window
	logger.Setup("WINDOW", verbose)
	ui.NewStatus().Run(quickConnect)
}
