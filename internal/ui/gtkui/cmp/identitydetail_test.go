package cmp

import (
	"testing"

	gtk "github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/telekom-mms/corp-net-indicator/internal/model"
)

func TestIdentityDetailInit(t *testing.T) {
	gtk.Init()

	NewIdentityDetails(model.NewContext(), make(chan bool))
}

// TODO add assertions and tests
