package cmp

import (
	"testing"

	gtk "github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/telekom-mms/corp-net-indicator/internal/model"
)

func TestInitVPNDetail(t *testing.T) {
	gtk.Init()

	NewVPNDetail(
		model.NewContext(),
		make(chan *model.Credentials),
		&gtk.Window{},
		func() ([]string, error) {
			return []string{}, nil
		},
		func() (int64, error) {
			return 0, nil
		},
		nil,
	)
}

// TODO add assertions and tests
