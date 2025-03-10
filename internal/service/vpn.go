package service

import (
	"crypto/x509"
	"encoding/pem"
	"os/exec"
	"time"

	"github.com/telekom-mms/corp-net-indicator/internal/logger"
	oc "github.com/telekom-mms/oc-daemon/pkg/client"
	"github.com/telekom-mms/oc-daemon/pkg/vpnstatus"
)

const THIRTY_DAYS = 30 * 24 * 60 * 60

type ErrConnect struct{ BaseError }
type ErrDisconnect struct{ BaseError }
type ErrGetServers struct{ BaseError }
type ErrGetCertDate struct{ BaseError }
type ErrGetVPNStatus struct{ BaseError }

func VPNInProgress(state vpnstatus.ConnectionState) bool {
	return state == vpnstatus.ConnectionStateConnecting || state == vpnstatus.ConnectionStateDisconnecting
}

type VPNService struct {
	client     oc.Client
	done       chan struct{}
	statusChan chan *vpnstatus.Status
}

var newVPNClient = func() (oc.Client, error) {
	return oc.NewClient(oc.LoadUserSystemConfig())
}

func NewVPNService() *VPNService {
	client, err := newVPNClient()
	if err != nil {
		panic(err)
	}
	return &VPNService{
		client:     client,
		done:       make(chan struct{}),
		statusChan: make(chan *vpnstatus.Status),
	}
}

func (v *VPNService) Subscribe() <-chan *vpnstatus.Status {
	go waitAndSubscribe[*vpnstatus.Status](v.client, v.statusChan, v.done)
	return v.statusChan
}

// triggers VPN connect
func (v *VPNService) ConnectWithPasswordAndServer(password string, server string) error {
	config := v.client.GetConfig()
	config.Password = password
	config.VPNServer = server
	v.client.SetConfig(config)
	// v.client.SetLogin(&logininfo.LoginInfo{})

	err := wrapErr(v.client.Authenticate(), &ErrConnect{})
	if err != nil {
		return err
	}

	return wrapErr(v.client.Connect(), &ErrConnect{})
}

func (v *VPNService) Disconnect() error {
	return wrapErr(v.client.Disconnect(), &ErrDisconnect{})
}

// Returns servers to connect to
func (v *VPNService) GetServers() ([]string, error) {
	result, err := v.client.Query()
	return result.Servers, wrapErr(err, &ErrGetServers{})
}

func (v *VPNService) GetCertExpireDate() (int64, bool, error) {
	// read cert
	out, err := exec.Command("p11tool", "--export", v.client.GetConfig().ClientCertificate).Output()
	if err != nil {
		logger.Logf("Warning: Can't read client certificate: %v", err)
		return -1, false, wrapErr(err, &ErrGetCertDate{})
	}
	// decode
	der, _ := pem.Decode(out)
	// parse
	cert, err := x509.ParseCertificate(der.Bytes)
	if err != nil {
		logger.Logf("Warning: Can't parse client certificate: %v", err)
		return -1, false, wrapErr(err, &ErrGetCertDate{})
	}
	// set value
	notAfter := cert.NotAfter.Unix()
	return notAfter, (notAfter - THIRTY_DAYS) < time.Now().Unix(), nil
}

func (v *VPNService) GetStatus() (*vpnstatus.Status, error) {
	status, err := v.client.Query()
	return status, wrapErr(err, &ErrGetVPNStatus{})
}

func (v *VPNService) Close() {
	close(v.done)
	v.client.Close()
	close(v.statusChan)
}
