package cmp

import (
	"testing"

	gtk "github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/telekom-mms/corp-net-indicator/internal/model"
)

func TestOpenAndClose(t *testing.T) {
	gtk.Init()

	d := newLoginDialog(&gtk.Window{}, func() ([]string, error) { return []string{}, nil })
	d.open(func(c *model.Credentials) {})
	d.close()
}

// TODO add assertions and tests
