package service

import (
	"github.com/telekom-mms/corp-net-indicator/internal/logger"
	ic "github.com/telekom-mms/fw-id-agent/pkg/client"
	"github.com/telekom-mms/fw-id-agent/pkg/status"
)

type ErrReLogin struct{ BaseError }
type ErrGetIdentityStatus struct{ BaseError }

func IdentityInProgress(state status.LoginState) bool {
	return state == status.LoginStateLoggingIn || state == status.LoginStateLoggingOut
}

type IdentityService struct {
	client     ic.Client
	done       chan struct{}
	statusChan chan *status.Status
}

var newIdentityClient = func() (ic.Client, error) {
	return ic.NewClient()
}

func NewIdentityService() *IdentityService {
	client, err := newIdentityClient()
	if err != nil {
		panic(err)
	}
	return &IdentityService{
		client:     client,
		statusChan: make(chan *status.Status, 10),
		done:       make(chan struct{}),
	}
}

// attaches to DBUS properties changed signal, maps to status and delivers them by returned channel
func (i *IdentityService) Subscribe() <-chan *status.Status {
	logger.Verbose("Start listening to identity status")
	go waitAndSubscribe[*status.Status](i.client, i.statusChan, i.done)
	return i.statusChan
}

// retrieves identity status
func (i *IdentityService) GetStatus() (*status.Status, error) {
	status, err := i.client.Query()
	return status, wrapErr(err, &ErrGetIdentityStatus{})
}

// triggers identity agent login
func (i *IdentityService) ReLogin() error {
	return wrapErr(i.client.ReLogin(), &ErrReLogin{})
}

// closes resources
func (i *IdentityService) Close() {
	close(i.done)
	i.client.Close()
	close(i.statusChan)
}
