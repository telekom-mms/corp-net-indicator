package cmp

import (
	gtk "github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/telekom-mms/corp-net-indicator/internal/assets"
)

type statusIcon struct {
	gtk.Image
}

// creates new status icon
func NewStatusIcon(status bool) *statusIcon {
	icon := &statusIcon{*gtk.NewImage()}
	icon.SetStatus(status)
	return icon
}

// changes icon -> true = green check, false = red alert
func (i *statusIcon) SetStatus(status bool) {
	if status {
		i.SetFromPixbuf(assets.GetPixbuf(assets.SVGCheck))
	} else {
		i.SetFromPixbuf(assets.GetPixbuf(assets.SVGAlert))
	}
}

// sets a icon to minus circle as value should be ignored
func (i *statusIcon) SetIgnore() {
	i.SetFromPixbuf(assets.GetPixbuf(assets.SVGMinus))
}
