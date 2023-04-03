package main

import (
	"flag"

	"github.com/telekom-mms/corp-net-indicator/internal/logger"
	"github.com/telekom-mms/corp-net-indicator/internal/tray"
	"github.com/telekom-mms/corp-net-indicator/internal/ui"
)

var runAsTray bool
var quickConnect bool
var verbose bool

func init() {
	flag.BoolVar(&runAsTray, "tray", false, "start as tray icon")
	flag.BoolVar(&quickConnect, "quick", false, "quick connect to vpn")
	flag.BoolVar(&verbose, "v", false, "verbose logging")
}

// starts indicator
// indicator can start as tray or window according flag value of -tray
func main() {
	flag.Parse()

	if runAsTray {
		// start as tray
		logger.Setup("TRAY", verbose)
		tray.New().Run()
	} else {
		// start as window
		logger.Setup("WINDOW", verbose)
		ui.NewStatus().Run(quickConnect)
	}
}
