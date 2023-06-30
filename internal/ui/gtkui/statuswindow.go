package gtkui

import (
	"github.com/diamondburned/gotk4/pkg/core/glib"
	gio "github.com/diamondburned/gotk4/pkg/gio/v2"
	gtk "github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/telekom-mms/corp-net-indicator/internal/config"
	"github.com/telekom-mms/corp-net-indicator/internal/i18n"
	"github.com/telekom-mms/corp-net-indicator/internal/logger"
	"github.com/telekom-mms/corp-net-indicator/internal/model"
	"github.com/telekom-mms/corp-net-indicator/internal/service"
	"github.com/telekom-mms/corp-net-indicator/internal/ui/gtkui/cmp"
	"github.com/telekom-mms/fw-id-agent/pkg/status"
	"github.com/telekom-mms/oc-daemon/pkg/vpnstatus"
)

// holds all window parts
type statusWindow struct {
	ctx              *model.Context
	quickConnect     bool
	initiallyOpened  bool
	vpnActionClicked chan *model.Credentials
	reLoginClicked   chan bool

	window       *gtk.ApplicationWindow
	notification *cmp.Notification

	identityDetail *cmp.IdentityDetails
	vpnDetail      *cmp.VPNDetail

	service *service.VPNService
}

// creates new status window
func NewStatusWindow(ctx *model.Context, vpnActionClicked chan *model.Credentials, reLoginClicked chan bool) *statusWindow {
	return &statusWindow{vpnActionClicked: vpnActionClicked, reLoginClicked: reLoginClicked, ctx: ctx}
}

// opens a new status window
// initialization is done with given status data
func (sw *statusWindow) Open(quickConnect bool, service *service.VPNService, onReady func()) {
	sw.quickConnect = quickConnect
	sw.service = service
	app := gtk.NewApplication("de.telekom-mms.corp-net-indicator", gio.ApplicationFlagsNone)
	app.ConnectActivate(func() {
		sw.window = gtk.NewApplicationWindow(app)
		sw.window.SetTitle("Corporate Network Status")
		sw.window.SetResizable(false)

		// header menu
		popover := gtk.NewPopover()
		aboutBtn := gtk.NewButtonWithLabel(i18n.L.Sprintf("About"))
		aboutBtn.ConnectClicked(func() {
			// about dialog
			aboutDialog := gtk.NewAboutDialog()
			aboutDialog.SetDestroyWithParent(true)
			aboutDialog.SetModal(true)
			aboutDialog.SetTransientFor(&sw.window.Window)
			aboutDialog.SetProgramName("Corporate Network Status")
			aboutDialog.SetComments(i18n.L.Sprintf("Program to show corporate network status."))
			aboutDialog.SetLogoIconName("applications-internet")
			commit := config.Commit
			if len(commit) > 11 {
				commit = config.Commit[0:11]
			}
			aboutDialog.SetVersion(config.Version + " (" + commit + ")")
			aboutDialog.SetCopyright("© 2023 The MMS Linux Dev Team")
			aboutDialog.SetAuthors([]string{"Hans Wippel", "Martin Schmitt", "Jan Dittberner", "Stefan Schubert"})

			aboutDialog.Show()
			popover.Hide()
		})
		popover.SetChild(aboutBtn)
		menuBtn := gtk.NewMenuButton()
		menuBtn.SetIconName("open-menu-symbolic")
		menuBtn.SetPopover(popover)

		// important to get rounded bottom corners
		headerBar := gtk.NewHeaderBar()
		headerBar.SetShowTitleButtons(true)
		icon := gtk.NewButtonFromIconName("applications-internet")
		icon.SetCanFocus(false)
		icon.SetCanTarget(false)
		headerBar.PackStart(icon)
		headerBar.PackEnd(menuBtn)
		sw.window.SetTitlebar(headerBar)

		// box for holding all detail boxes
		details := gtk.NewBox(gtk.OrientationVertical, 0)
		details.SetMarginTop(30)
		details.SetMarginBottom(30)
		details.SetMarginStart(60)
		details.SetMarginEnd(60)

		// create details
		sw.identityDetail = cmp.NewIdentityDetails(sw.ctx, sw.reLoginClicked)
		sw.vpnDetail = cmp.NewVPNDetail(sw.ctx, sw.vpnActionClicked, &sw.window.Window,
			func() ([]string, error) {
				servers, err := sw.service.GetServers()
				if err != nil {
					go sw.NotifyError(err)
				}
				return servers, err
			},
			func() (int64, error) {
				date, err := sw.service.GetCertExpireDate()
				if err != nil {
					go sw.NotifyError(err)
				}
				return date, err
			},
			sw.identityDetail)

		// append all boxes
		details.Append(sw.identityDetail)
		details.Append(sw.vpnDetail)

		// create notification and overlay for them
		sw.notification = cmp.NewNotification()
		overlay := gtk.NewOverlay()
		// details are added as overlay child
		overlay.SetChild(details)
		overlay.AddOverlay(sw.notification.Revealer)

		// show window
		sw.window.SetChild(overlay)
		sw.window.Show()

		// call on ready
		go onReady()
	})

	// this call blocks until window is closed
	if code := app.Run([]string{}); code > 0 {
		logger.Log("Failed to open window")
	}
}

// applies identity status
func (sw *statusWindow) ApplyIdentityStatus(status *status.Status) {
	if sw.window == nil {
		return
	}
	sw.identityDetail.Apply(status)
}

// applies vpn status
func (sw *statusWindow) ApplyVPNStatus(status *vpnstatus.Status) {
	if sw.window == nil {
		return
	}
	sw.vpnDetail.Apply(status, func(vpnConnected bool) {
		if sw.quickConnect {
			if vpnConnected {
				logger.Verbose("Closing window after quick connect")
				sw.Close()
			}
			if !sw.initiallyOpened && !vpnConnected {
				sw.initiallyOpened = true
				logger.Verbose("Open window on quick connect")
				sw.vpnDetail.OnActionClicked()
			}
		}
	})
}

// closes window
func (sw *statusWindow) Close() {
	if sw.window == nil {
		return
	}
	sw.vpnDetail.Close()
	sw.window.Close()
	sw.window.Destroy()
}

// triggers notification to show for given error
func (sw *statusWindow) NotifyError(err error) {
	if sw.window == nil {
		return
	}
	var msg string
	switch err.(type) {
	case *service.ErrConnect:
		msg = i18n.L.Sprintf("Could not connect. Please Retry.")
	case *service.ErrDisconnect:
		msg = i18n.L.Sprintf("Could not disconnect. Please Retry.")
	case *service.ErrGetVPNStatus:
		msg = i18n.L.Sprintf("Could not query current VPN status.")
	case *service.ErrGetServers:
		msg = i18n.L.Sprintf("Could not query server list.")
	case *service.ErrGetCertDate:
		msg = i18n.L.Sprintf("Could not query certification expire date.")
	case *service.ErrReLogin:
		msg = i18n.L.Sprintf("Could not refresh identity login. Please Retry.")
	case *service.ErrGetIdentityStatus:
		msg = i18n.L.Sprintf("Could not query current Identity status.")
	default:
		msg = i18n.L.Sprintf("Error: [%v]", err)
	}
	glib.IdleAdd(func() {
		sw.vpnDetail.SetButtonsAfterProgress()
		sw.notification.Show(msg)
	})
}
