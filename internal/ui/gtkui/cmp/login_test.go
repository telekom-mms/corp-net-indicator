package cmp

import (
	"testing"

	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/stretchr/testify/assert"
	"github.com/telekom-mms/corp-net-indicator/internal/model"
)

func TestOpenAndClose(t *testing.T) {
	gtk.Init()

	d := newLoginDialog(&gtk.Window{}, func() ([]string, error) { return []string{}, nil })
	err := d.open(func(c *model.Credentials) {})
	assert.Nil(t, err)
	d.close()
}

// TODO add assertions and tests
