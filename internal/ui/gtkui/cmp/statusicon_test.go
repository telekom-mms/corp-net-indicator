package cmp

import (
	"testing"

	gtk "github.com/diamondburned/gotk4/pkg/gtk/v4"
)

func TestIconInit(t *testing.T) {
	gtk.Init()

	NewStatusIcon(true)
}

// TODO add assertions and tests
