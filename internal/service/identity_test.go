package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	ic "github.com/telekom-mms/fw-id-agent/pkg/client"
	"github.com/telekom-mms/fw-id-agent/pkg/status"
)

type testIdentityClient struct {
	ic.Client
	queryCalled   bool
	reLoginCalled bool
	closeCalled   bool
	status        chan *status.Status
}

func (c *testIdentityClient) Query() (*status.Status, error) {
	c.queryCalled = true
	return nil, nil
}
func (c *testIdentityClient) ReLogin() error {
	c.reLoginCalled = true
	return nil
}
func (c *testIdentityClient) Close() error {
	c.closeCalled = true
	return nil
}
func (c *testIdentityClient) Ping() error {
	return nil
}
func (c *testIdentityClient) Subscribe() (chan *status.Status, error) {
	return c.status, nil
}
func setupIdentityClient() *testIdentityClient {
	testC := &testIdentityClient{status: make(chan *status.Status)}
	newIdentityClient = func() (ic.Client, error) { return testC, nil }
	return testC
}

func TestGetIdentityStatus(t *testing.T) {
	testC := setupIdentityClient()
	c := NewIdentityService()

	status, err := c.GetStatus()
	c.Close()
	assert.Nil(t, err)
	assert.Nil(t, status)
	assert.Equal(t, true, testC.queryCalled)
	assert.Equal(t, true, testC.closeCalled)
}

func TestReLogin(t *testing.T) {
	testC := setupIdentityClient()
	c := NewIdentityService()

	err := c.ReLogin()
	c.Close()
	assert.Nil(t, err)
	assert.Equal(t, true, testC.reLoginCalled)
	assert.Equal(t, true, testC.closeCalled)
}

func TestIdentitySubscribe(t *testing.T) {
	testC := setupIdentityClient()
	c := NewIdentityService()

	status := &status.Status{}
	statusChan := c.Subscribe()
	testC.status <- status
	assert.Equal(t, status, <-statusChan)
	c.Close()
	assert.Equal(t, true, testC.closeCalled)
}
