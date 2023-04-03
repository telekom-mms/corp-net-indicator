package cmp

import (
	"time"

	"github.com/diamondburned/gotk4/pkg/core/glib"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

type Notification struct {
	text     *gtk.Label
	Revealer *gtk.Revealer
	timer    *time.Timer
}

// creates notification to show error messages
func NewNotification() *Notification {
	n := &Notification{}

	// revealer shows and hides box with notification
	n.Revealer = gtk.NewRevealer()
	n.Revealer.SetHAlign(gtk.AlignCenter)
	n.Revealer.SetVAlign(gtk.AlignStart)

	// box to hold notification
	box := gtk.NewBox(gtk.OrientationHorizontal, 20)
	box.SetCanFocus(false)
	box.SetVAlign(gtk.AlignStart)
	box.AddCSSClass("app-notification")
	box.SetMarginStart(30)
	box.SetMarginEnd(30)
	box.SetOpacity(0.8)

	// notification itself with wrapping
	n.text = gtk.NewLabel("")
	n.text.SetWrap(true)
	n.text.SetHExpand(true)

	// button to close notification
	btn := gtk.NewButtonFromIconName("window-close-symbolic")
	btn.SetHAlign(gtk.AlignEnd)
	btn.SetVAlign(gtk.AlignCenter)
	btn.SetHExpand(false)
	btn.SetVExpand(false)
	btn.ConnectClicked(func() {
		n.Revealer.SetRevealChild(false)
	})

	// build structure
	box.Append(n.text)
	box.Append(btn)
	n.Revealer.SetChild(box)

	return n
}

// shows notification and hides them after 10s
func (n *Notification) Show(text string) {
	if n.timer != nil {
		n.timer.Stop()
	}
	n.timer = time.NewTimer(time.Second * 10)
	n.text.SetLabel(text)
	n.Revealer.SetRevealChild(true)
	go func() {
		<-n.timer.C
		glib.IdleAdd(func() {
			n.Revealer.SetRevealChild(false)
		})
	}()
}
