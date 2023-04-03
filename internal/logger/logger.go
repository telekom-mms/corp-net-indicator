package logger

import (
	"fmt"
	"log"
	"os"
)

// holds logger and verbose setting
var (
	IsVerbose = false
	std       = log.New(os.Stderr, "", log.LstdFlags)
)

// setup have to be called on start to configure logger instance
func Setup(prefix string, verbose bool) {
	IsVerbose = verbose
	std.SetPrefix(prefix + " ")
	if verbose {
		std.SetFlags(log.LstdFlags | log.Lshortfile)
	} else {
		std.SetFlags(log.LstdFlags)
	}
}

// default non verbose logging
func Log(v ...any) {
	std.Output(2, fmt.Sprintln(v...))
}

// default non verbose formatted logging
func Logf(format string, v ...any) {
	std.Output(2, fmt.Sprintf(format, v...))
}

// verbose logging
func Verbose(v ...any) {
	if IsVerbose {
		std.Output(2, fmt.Sprintln(v...))
	}
}

// verbose formatted logging
func Verbosef(format string, v ...any) {
	if IsVerbose {
		std.Output(2, fmt.Sprintf(format, v...))
	}
}
