package cmp

import (
	"github.com/diamondburned/gotk4/pkg/core/glib"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/telekom-mms/corp-net-indicator/internal/i18n"
	"github.com/telekom-mms/corp-net-indicator/internal/model"
	"github.com/telekom-mms/corp-net-indicator/internal/util"
	"github.com/telekom-mms/fw-id-agent/pkg/status"
)

type IdentityDetails struct {
	detail
	ctx *model.Context

	reLoginClicked chan bool

	loggedInImg      *statusIcon
	keepAliveAtLabel *gtk.Label
	krbEndTimeLabel  *gtk.Label
	reLoginBtn       *gtk.Button
	reLoginSpinner   *gtk.Spinner
}

// creates new identity details box
func NewIdentityDetails(
	// shared context
	ctx *model.Context,
	// channel to notify reLogin clicks
	reLoginClicked chan bool) *IdentityDetails {

	id := &IdentityDetails{detail: newDetail(), ctx: ctx, reLoginClicked: reLoginClicked}

	// keep all elements in struct to instrument their behavior
	id.reLoginBtn = gtk.NewButtonWithLabel(i18n.L.Sprintf("ReLogin"))
	id.reLoginBtn.SetHAlign(gtk.AlignEnd)
	// set click handler
	id.reLoginBtn.ConnectClicked(id.onReLoginClicked)
	id.loggedInImg = NewStatusIcon(false)
	id.reLoginSpinner = gtk.NewSpinner()
	id.reLoginSpinner.SetHAlign(gtk.AlignEnd)
	id.keepAliveAtLabel = gtk.NewLabel(util.DefaultValue)
	id.krbEndTimeLabel = gtk.NewLabel(util.DefaultValue)

	// build details box and attach them to the values and actions
	id.
		buildBase(i18n.L.Sprintf("Identity Details")).
		addRow(i18n.L.Sprintf("Logged in"), id.reLoginSpinner, id.reLoginBtn, id.loggedInImg).
		addRow(i18n.L.Sprintf("Last Refresh"), id.keepAliveAtLabel).
		addRow(i18n.L.Sprintf("Kerberos ticket valid until"), id.krbEndTimeLabel)

	id.reLoginBtn.SetSensitive(false)
	id.reLoginSpinner.Start()

	return id
}

// applies new status to identity details
func (id *IdentityDetails) Apply(status *status.Status) {
	glib.IdleAdd(func() {
		ctx := id.ctx.Read()
		// quick path for in progress updates
		if ctx.IdentityInProgress || ctx.VPNInProgress {
			if ctx.IdentityInProgress {
				id.reLoginSpinner.Start()
			}
			id.reLoginBtn.SetSensitive(false)
			return
		}
		// set new status values
		loggedIn := status.LoginState.LoggedIn()
		id.loggedInImg.SetStatus(loggedIn)
		id.setReLoginBtn(loggedIn)
		id.keepAliveAtLabel.SetText(util.FormatDate(status.LastKeepAlive))
		id.krbEndTimeLabel.SetText(util.FormatDate(status.KerberosTGT.StartTime))
		// set button state
		id.setButtonAndLoginState()
	})
}

// action handler to trigger login to identity service
func (id *IdentityDetails) onReLoginClicked() {
	go func() {
		glib.IdleAdd(func() {
			id.reLoginSpinner.Start()
			id.reLoginBtn.SetSensitive(false)
		})
		id.reLoginClicked <- true
	}()
}

// sets button state -> true = activated, false deactivated
func (id *IdentityDetails) setReLoginBtn(status bool) {
	id.reLoginBtn.SetSensitive(status)
	id.reLoginSpinner.Stop()
}

// set logged in icon state and button state
func (id *IdentityDetails) setButtonAndLoginState() {
	ctx := id.ctx.Read()
	if !ctx.Connected && !ctx.TrustedNetwork {
		id.loggedInImg.SetStatus(false)
		id.setReLoginBtn(false)
	} else if !ctx.IdentityInProgress {
		id.setReLoginBtn(true)
	}
}
