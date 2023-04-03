package main

import (
	"github.com/godbus/dbus/v5/prop"
)

type agent struct {
	props    *prop.Properties
	simulate bool
}
