package tray

import (
	"errors"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"sync/atomic"

	"github.com/slytomcat/systray"
	"github.com/telekom-mms/corp-net-indicator/internal/assets"
	"github.com/telekom-mms/corp-net-indicator/internal/i18n"
	"github.com/telekom-mms/corp-net-indicator/internal/logger"
	"github.com/telekom-mms/corp-net-indicator/internal/model"
	"github.com/telekom-mms/corp-net-indicator/internal/service"
	"github.com/telekom-mms/oc-daemon/pkg/vpnstatus"
)

type tray struct {
	ctx *model.Context

	statusItem   *systray.MenuItem
	actionItem   *systray.MenuItem
	startSystray func()
	quitSystray  func()

	vSer *service.VPNService

	window    atomic.Pointer[os.Process]
	closeChan chan struct{}

	windowInitiallyOpened bool
}

// starts tray
func New() *tray {
	t := &tray{ctx: model.NewContext()}
	// create tray
	t.startSystray, t.quitSystray = systray.RunWithExternalLoop(t.onReady, func() {})
	return t
}

// init tray
func (t *tray) onReady() {
	// set up menu
	t.statusItem = systray.AddMenuItem(i18n.L.Sprintf("Status"), i18n.L.Sprintf("Show Status"))
	t.statusItem.SetIcon(assets.GetIcon(assets.Status))
	t.actionItem = systray.AddMenuItem(i18n.L.Sprintf("Connect VPN"), i18n.L.Sprintf("Connect to VPN"))
	t.actionItem.SetIcon(assets.GetIcon(assets.Connect))
	t.actionItem.Hide()
}

// opens corp-net-indicator window
func (t *tray) OpenWindow(quickConnect bool) {
	t.closeWindow()
	self, err := os.Executable()
	if err != nil {
		logger.Log(err)
		return
	}
	trayBin := strings.TrimSuffix(self, "-tray")
	if _, err := os.Stat(trayBin); err != nil {
		logger.Log(err)
		return
	}
	var cmd *exec.Cmd
	args := []string{}
	if quickConnect {
		args = append(args, "-quick")
	}
	if logger.IsVerbose {
		args = append(args, "-v")
	}
	cmd = exec.Command(trayBin, args...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	closeChan := make(chan struct{})
	err = cmd.Start()
	if err != nil {
		logger.Log(err)
		return
	}
	go func() {
		defer func() {
			t.window.Store(nil)
			close(closeChan)
		}()

		_, err := cmd.Process.Wait()
		if err != nil {
			logger.Verbose(err)
		}

		logger.Verbose("Waited for closing window")
	}()
	t.closeChan = closeChan
	t.window.Store(cmd.Process)
}

func (t *tray) closeWindow() {
	window := t.window.Load()
	if window == nil {
		return
	}
	err := window.Signal(os.Interrupt)
	if err != nil && !errors.Is(err, os.ErrProcessDone) {
		logger.Verbosef("SIGINT not working: %v\n", err)
		err = window.Kill()
		if err != nil && !errors.Is(err, os.ErrProcessDone) {
			logger.Verbosef("SIGKILL not working: %v\n", err)
		}
	}
	<-t.closeChan
}

func (t *tray) Run() {
	// start tray
	t.startSystray()
	// create services
	t.vSer = service.NewVPNService()
	iSer := service.NewIdentityService()
	wSer := service.NewWatcher()

	// listen to status changes
	vChan := t.vSer.Subscribe()
	iChan := iSer.Subscribe()

	// catch user login
	wChan := wSer.Listen()

	// catch interrupt and clean up
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	// error channel
	e := make(chan error, 1)
	// main loop
	for {
		select {
		// handle tray menu clicks
		case <-t.statusItem.ClickedCh:
			t.windowInitiallyOpened = true
			logger.Verbose("Open window to connect")

			t.OpenWindow(false)
		case <-t.actionItem.ClickedCh:
			t.windowInitiallyOpened = true
			if t.ctx.Read().Connected {
				logger.Verbose("Try to disconnect")

				t.actionItem.Disable()
				go func() {
					e <- t.vSer.Disconnect()
				}()
			} else {
				logger.Verbose("Open window to quick connect")

				t.OpenWindow(true)
			}
		// handle disconnect errors
		case err := <-e:
			if err != nil {
				logger.Logf("Error: %v\n", err)
				t.actionItem.Enable()
			}
		// handle status updates
		case status := <-iChan:
			logger.Verbosef("Apply identity status: %+v\n", status)

			ctx := t.ctx.Write(func(ctx *model.ContextValues) {
				ctx.IdentityInProgress = service.IdentityInProgress(status.LoginState)
				ctx.LoggedIn = status.LoginState.LoggedIn()
			})
			t.apply(ctx)
		case status := <-vChan:
			logger.Verbosef("Apply vpn status: %+v\n", status)

			ctx := t.ctx.Write(func(ctx *model.ContextValues) {
				ctx.VPNInProgress = service.VPNInProgress(status.ConnectionState)
				ctx.Connected = status.ConnectionState.Connected()
				ctx.TrustedNetwork = status.TrustedNetwork.Trusted()
			})
			t.apply(ctx)
			// open window, if needed
			if !t.windowInitiallyOpened {
				t.windowInitiallyOpened = t.openWindowIfNeeded(status)
			}
		case <-wChan:
			logger.Verbose("Watcher signal received")
			status, err := t.vSer.GetStatus()
			if err != nil {
				logger.Logf("Error: %v\n", err)
				os.Exit(1)
			}
			t.openWindowIfNeeded(status)
		case <-c:
			logger.Verbose("Received SIGINT -> closing")

			t.closeWindow()
			t.vSer.Close()
			iSer.Close()
			wSer.Close()
			t.quitSystray()
			return
		}
	}
}

func (t *tray) apply(ctx model.ContextValues) {
	// icon
	if ctx.LoggedIn && (ctx.Connected || ctx.TrustedNetwork) {
		systray.SetIcon(assets.GetIcon(assets.Umbrella))
	} else {
		if ctx.Connected || ctx.TrustedNetwork {
			systray.SetIcon(assets.GetIcon(assets.ShieldOn))
		} else {
			systray.SetIcon(assets.GetIcon(assets.ShieldOff))
		}
	}
	// action menu item
	if ctx.VPNInProgress {
		t.actionItem.Disable()
	} else {
		t.actionItem.Enable()
	}
	if ctx.Connected {
		t.actionItem.SetTitle(i18n.L.Sprintf("Disconnect VPN"))
		t.actionItem.SetIcon(assets.GetIcon(assets.Disconnect))
		t.actionItem.Show()
	} else {
		t.actionItem.SetTitle(i18n.L.Sprintf("Connect VPN"))
		t.actionItem.SetIcon(assets.GetIcon(assets.Connect))
	}
	if ctx.TrustedNetwork {
		t.actionItem.Hide()
	} else {
		t.actionItem.Show()
	}
}

// opens window if needed
func (t *tray) openWindowIfNeeded(status *vpnstatus.Status) bool {
	if status.TrustedNetwork == vpnstatus.TrustedNetworkNotTrusted &&
		status.ConnectionState <= vpnstatus.ConnectionStateDisconnected {
		t.OpenWindow(true)
		return true
	}
	if status.TrustedNetwork == vpnstatus.TrustedNetworkTrusted {
		go func() {
			_, warn, _ := t.vSer.GetCertExpireDate()
			if warn {
				t.OpenWindow(false)
			}
		}()
	}
	return false
}
