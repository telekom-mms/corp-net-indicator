package ui

import (
	"os"
	"os/signal"

	"github.com/telekom-mms/corp-net-indicator/internal/logger"
	"github.com/telekom-mms/corp-net-indicator/internal/model"
	"github.com/telekom-mms/corp-net-indicator/internal/service"
	"github.com/telekom-mms/corp-net-indicator/internal/ui/gtkui"
	"github.com/telekom-mms/fw-id-agent/pkg/status"
	"github.com/telekom-mms/oc-daemon/pkg/vpnstatus"
)

// minimal interface to interact with an ui implementation
type StatusWindow interface {
	Open(quickConnect bool, getServers func() ([]string, error), onReady func())
	Close()
	ApplyIdentityStatus(status *status.Status)
	ApplyVPNStatus(status *vpnstatus.Status)
	NotifyError(err error)
}

// holds data channels for updates and a window handle
// is used to free memory after closing window
type Status struct {
	ctx *model.Context

	connectDisconnectClicked chan *model.Credentials
	reLoginClicked           chan bool
	done                     chan struct{}

	window StatusWindow
}

func NewStatus() *Status {
	s := &Status{
		ctx:                      model.NewContext(),
		connectDisconnectClicked: make(chan *model.Credentials),
		reLoginClicked:           make(chan bool),
		done:                     make(chan struct{}),
	}
	s.window = gtkui.NewStatusWindow(s.ctx, s.connectDisconnectClicked, s.reLoginClicked)
	return s
}

func (s *Status) Run(quickConnect bool) {
	// create services
	vSer := service.NewVPNService()
	iSer := service.NewIdentityService()

	// listen to status changes
	vChan := vSer.Subscribe()
	iChan := iSer.Subscribe()

	// catch interrupt and clean up
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	s.window.Open(quickConnect, vSer.GetServers, func() {
		for {
			select {
			// handle window clicks
			case connect := <-s.connectDisconnectClicked:
				s.ctx.Write(func(ctx *model.ContextValues) {
					ctx.VPNInProgress = true
				})
				if connect != nil {
					logger.Verbose("Open dialog to connect to VPN")

					if err := vSer.ConnectWithPasswordAndServer(connect.Password, connect.Server); err != nil {
						s.handleError(err)
					}
				} else {
					logger.Verbose("Tray to disconnect")

					if err := vSer.Disconnect(); err != nil {
						s.handleError(err)
					}
				}
			case <-s.reLoginClicked:
				logger.Verbose("Try to login to identity service")

				s.ctx.Write(func(ctx *model.ContextValues) {
					ctx.IdentityInProgress = true
				})
				if err := iSer.ReLogin(); err != nil {
					s.handleError(err)
				}
			case status := <-iChan:
				logger.Verbosef("Apply identity status: %+v\n", status)

				s.ctx.Write(func(ctx *model.ContextValues) {
					ctx.IdentityInProgress = service.IdentityInProgress(status.LoginState)
					ctx.LoggedIn = status.LoginState.LoggedIn()
				})
				go s.window.ApplyIdentityStatus(status)
			case status := <-vChan:
				logger.Verbosef("Apply vpn status: %+v\n", status)

				s.ctx.Write(func(ctx *model.ContextValues) {
					ctx.VPNInProgress = service.VPNInProgress(status.ConnectionState)
					ctx.Connected = status.ConnectionState.Connected()
					ctx.TrustedNetwork = status.TrustedNetwork.Trusted()
				})
				go s.window.ApplyVPNStatus(status)
			case <-c:
				logger.Verbose("Received SIGINT -> closing")
				defer s.window.Close()
				return
			case <-s.done:
				return
			}
		}
	})

	logger.Verbose("Window closed")
	close(s.done)

	vSer.Close()
	iSer.Close()
}

func (s *Status) handleError(err error) {
	logger.Logf("Error: %v\n", err)

	s.ctx.Write(func(ctx *model.ContextValues) {
		ctx.VPNInProgress = false
		ctx.IdentityInProgress = false
	})
	go s.window.NotifyError(err)
}
