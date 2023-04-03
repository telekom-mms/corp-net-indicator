package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	oc "github.com/telekom-mms/oc-daemon/pkg/client"
	"github.com/telekom-mms/oc-daemon/pkg/vpnstatus"
)

type testVPNClient struct {
	oc.Client
	queryCalled      bool
	authCalled       bool
	connectCalled    bool
	disconnectCalled bool
	closeCalled      bool
	status           chan *vpnstatus.Status
	config           *oc.Config
}

func (c *testVPNClient) Query() (*vpnstatus.Status, error) {
	c.queryCalled = true
	return &vpnstatus.Status{Servers: []string{"server"}}, nil
}
func (c *testVPNClient) Authenticate() error {
	c.authCalled = true
	return nil
}
func (c *testVPNClient) Connect() error {
	c.connectCalled = true
	return nil
}
func (c *testVPNClient) Disconnect() error {
	c.disconnectCalled = true
	return nil
}
func (c *testVPNClient) Close() error {
	c.closeCalled = true
	return nil
}
func (c *testVPNClient) Ping() error {
	return nil
}
func (c *testVPNClient) Subscribe() (chan *vpnstatus.Status, error) {
	return c.status, nil
}
func (c *testVPNClient) GetConfig() *oc.Config {
	return c.config
}
func (c *testVPNClient) SetConfig(config *oc.Config) {
	c.config = config
}
func setupVPNClient() *testVPNClient {
	testC := &testVPNClient{status: make(chan *vpnstatus.Status), config: &oc.Config{}}
	newVPNClient = func() (oc.Client, error) { return testC, nil }
	return testC
}

func TestGetVPNStatus(t *testing.T) {
	setupVPNClient()
	c := NewVPNService()
	defer c.Close()

	status, err := c.GetStatus()
	assert.Nil(t, err)
	assert.Equal(t, &vpnstatus.Status{Servers: []string{"server"}}, status)
}

func TestVPNSubscribe(t *testing.T) {
	testC := setupVPNClient()
	c := NewVPNService()

	status := &vpnstatus.Status{}
	statusChan := c.Subscribe()
	testC.status <- status
	assert.Equal(t, status, <-statusChan)
	c.Close()
	assert.Equal(t, true, testC.closeCalled)
}

func TestConnect(t *testing.T) {
	testC := setupVPNClient()
	c := NewVPNService()

	err := c.ConnectWithPasswordAndServer("pass", "server")
	assert.Nil(t, err)
	assert.Equal(t, true, testC.connectCalled)
	assert.Equal(t, true, testC.authCalled)
	assert.Equal(t, &oc.Config{Password: "pass", VPNServer: "server"}, testC.config)

	c.Close()
	assert.Equal(t, true, testC.closeCalled)
}

func TestDisconnect(t *testing.T) {
	testC := setupVPNClient()
	c := NewVPNService()

	err := c.Disconnect()
	assert.Nil(t, err)
	assert.Equal(t, true, testC.disconnectCalled)

	c.Close()
	assert.Equal(t, true, testC.closeCalled)
}

func TestGetServers(t *testing.T) {
	testC := setupVPNClient()
	c := NewVPNService()

	servers, err := c.GetServers()
	assert.Nil(t, err)
	assert.Equal(t, []string{"server"}, servers)
	assert.Equal(t, true, testC.queryCalled)

	c.Close()
	assert.Equal(t, true, testC.closeCalled)
}
