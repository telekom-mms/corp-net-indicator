package service

import (
	dbus "github.com/godbus/dbus/v5"
	"github.com/telekom-mms/corp-net-indicator/internal/logger"
)

const (
	IFACE  string = "org.freedesktop.login1.Session"
	SIGNAL string = IFACE + ".Unlock"
)

type dbusConn interface {
	AddMatchSignal(options ...dbus.MatchOption) error
	Signal(ch chan<- *dbus.Signal)
	Close() error
}

var connectDbus = func() (dbusConn, error) {
	return dbus.ConnectSystemBus()
}

type Watcher struct {
	conn   dbusConn
	signal chan struct{}
}

// Creates new watcher
func NewWatcher() *Watcher {
	conn, err := connectDbus()
	if err != nil {
		panic(err)
	}
	return &Watcher{conn: conn, signal: make(chan struct{}, 1)}
}

// Listen to login events
func (w *Watcher) Listen() <-chan struct{} {
	logger.Verbose("Listen to user actions")
	// setup signal
	opts := []dbus.MatchOption{
		dbus.WithMatchInterface(IFACE),
	}
	err := w.conn.AddMatchSignal(opts...)
	if err != nil {
		panic(err)
	}

	// create signal channel
	c := make(chan *dbus.Signal, 10)
	w.conn.Signal(c)

	go func() {
		for sig := range c {
			if sig.Name == SIGNAL {
				select {
				case w.signal <- struct{}{}:
				default:
				}
			}
		}
	}()

	return w.signal
}

// Cleanup dbus connection and channel
func (w *Watcher) Close() {
	w.conn.Close()
	close(w.signal)
}
