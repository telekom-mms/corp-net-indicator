package service

import (
	"time"

	"github.com/telekom-mms/corp-net-indicator/internal/logger"
)

var pollInterval time.Duration = 5

type client[T interface{}] interface {
	Ping() error
	Subscribe() (chan T, error)
}

func waitAndSubscribe[T interface{}](client client[T], statusChan chan T, done chan struct{}) {
	for {
		select {
		case <-done:
			return
		default:
		}
		if client.Ping() == nil {
			break
		}
		logger.Verbosef("Wait %d seconds for service to come up...", pollInterval)
		time.Sleep(time.Second * pollInterval)
	}
	c, err := client.Subscribe()
	if err != nil {
		panic(err)
	}
	for status := range c {
		select {
		case statusChan <- status:
		case <-done:
			return
		}
	}
}

type WrappedError interface {
	Error() string
	setError(err error)
}

type BaseError struct {
	wrapped error
}

func (e *BaseError) Error() string {
	return e.wrapped.Error()
}

func (e *BaseError) setError(err error) {
	e.wrapped = err
}

func wrapErr[T WrappedError](toWrap error, wrapper T) error {
	if toWrap != nil {
		logger.Logf("Client error: %v", toWrap)
		wrapper.setError(toWrap)
		return wrapper
	}
	return nil
}
