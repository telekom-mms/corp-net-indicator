package service

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testSubscribeClient struct {
	m         sync.Mutex
	pingCount int
	status    chan string
	err       bool
	subCalled bool
}

func (c *testSubscribeClient) Ping() error {
	c.m.Lock()
	defer c.m.Unlock()
	if c.pingCount == 0 {
		c.pingCount++
		return fmt.Errorf("Not available")
	}
	return nil
}

func (c *testSubscribeClient) Subscribe() (chan string, error) {
	c.m.Lock()
	defer c.m.Unlock()
	c.subCalled = true
	if c.err {
		return nil, fmt.Errorf("Failed")
	}
	return c.status, nil
}

func TestSubscribe(t *testing.T) {
	pollInterval = 1
	client := &testSubscribeClient{status: make(chan string)}
	status := make(chan string, 10)
	done := make(chan struct{})

	go waitAndSubscribe[string](client, status, done)

	client.status <- "msg"
	assert.Equal(t, 1, client.pingCount)
	assert.Equal(t, "msg", <-status)
	close(client.status)
}

func TestSubscribeErr(t *testing.T) {
	pollInterval = 1
	client := &testSubscribeClient{status: make(chan string), err: true}
	defer func() {
		assert.Equal(t, fmt.Errorf("Failed"), recover())
		assert.Equal(t, 1, client.pingCount)
	}()
	status := make(chan string, 10)
	done := make(chan struct{})

	waitAndSubscribe[string](client, status, done)

	t.Errorf("should have panicked")
}

func TestSubscribeCancel(t *testing.T) {
	pollInterval = 1
	client := &testSubscribeClient{status: make(chan string)}
	status := make(chan string, 10)
	done := make(chan struct{})
	c := make(chan struct{})

	go func() {
		waitAndSubscribe[string](client, status, done)
		close(c)
	}()

	close(done)
	<-c
	assert.Equal(t, 0, client.pingCount)
}

func TestSubscribeCancel2(t *testing.T) {
	pollInterval = 1
	client := &testSubscribeClient{status: make(chan string)}
	status := make(chan string)
	done := make(chan struct{})

	go waitAndSubscribe[string](client, status, done)

	client.status <- "msg"
	close(done)
	assert.Equal(t, 1, client.pingCount)
}
