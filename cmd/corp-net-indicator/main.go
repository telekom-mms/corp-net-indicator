package main

import (
	"flag"

	"github.com/telekom-mms/corp-net-indicator/internal/config"
	"github.com/telekom-mms/corp-net-indicator/internal/logger"
	"github.com/telekom-mms/corp-net-indicator/internal/tray"
)

var verbose bool

func init() {
	flag.BoolVar(&verbose, "v", false, "verbose logging")
}

// starts indicator as tray
func main() {
	flag.Parse()

	// start as tray
	logger.Setup("TRAY", verbose)
	logger.Logf("Start corp-net-indicator tray [%s-%s]", config.Version, config.Commit)
	tray.New().Run()
}
