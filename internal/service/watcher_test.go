package service

import (
	"testing"

	"github.com/godbus/dbus/v5"
	"github.com/stretchr/testify/assert"
)

type testConn struct {
	matchOptions []dbus.MatchOption
	signalChan   chan<- *dbus.Signal
	closeCalled  bool
}

func (c *testConn) AddMatchSignal(options ...dbus.MatchOption) error {
	c.matchOptions = options
	return nil
}
func (c *testConn) Signal(ch chan<- *dbus.Signal) {
	c.signalChan = ch
}
func (c *testConn) Close() error {
	c.closeCalled = true
	return nil
}
func setupDbusConn() *testConn {
	conn := &testConn{}
	connectDbus = func() (dbusConn, error) {
		return conn, nil
	}
	return conn
}

func TestWatcher(t *testing.T) {
	testConn := setupDbusConn()
	w := NewWatcher()

	c := w.Listen()
	go func() {
		testConn.signalChan <- &dbus.Signal{Name: SIGNAL}
	}()
	<-c
	w.Close()

	assert.Equal(t, []dbus.MatchOption{
		dbus.WithMatchInterface(IFACE),
	}, testConn.matchOptions)
	assert.Equal(t, true, testConn.closeCalled)
}
