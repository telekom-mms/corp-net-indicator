package service

import (
	oc "github.com/telekom-mms/oc-daemon/pkg/client"
	"github.com/telekom-mms/oc-daemon/pkg/vpnstatus"
)

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

	err := v.client.Authenticate()
	if err != nil {
		return err
	}

	return v.client.Connect()
}

func (v *VPNService) Disconnect() error {
	return v.client.Disconnect()
}

// Returns servers to connect to
func (v *VPNService) GetServers() ([]string, error) {
	result, err := v.client.Query()
	return result.Servers, err
}

func (v *VPNService) GetStatus() (*vpnstatus.Status, error) {
	return v.client.Query()
}

func (v *VPNService) Close() {
	close(v.done)
	v.client.Close()
	close(v.statusChan)
}
