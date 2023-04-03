package main

import (
	"fmt"
	"log"
	"sync/atomic"
	"time"

	"github.com/godbus/dbus/v5"
	"github.com/godbus/dbus/v5/introspect"
	"github.com/godbus/dbus/v5/prop"
	"github.com/telekom-mms/fw-id-agent/pkg/status"
	"github.com/telekom-mms/oc-daemon/pkg/vpnstatus"
)

const V_DBUS_SERVICE_NAME = "com.telekom_mms.oc_daemon.Daemon"
const V_DBUS_OBJECT_PATH = "/com/telekom_mms/oc_daemon/Daemon"

type vpnAgent struct {
	agent
}

var counter atomic.Uint32

func (a vpnAgent) Connect(_, _, _, _, _ string) *dbus.Error {
	log.Println("VPN: Connect called!")
	a.props.SetMust(V_DBUS_SERVICE_NAME, "ConnectionState", vpnstatus.ConnectionStateConnecting)
	go func() {
		var now int64 = 0
		if a.simulate {
			time.Sleep(time.Second * 5)
			now = time.Now().Unix()
		}
		a.props.SetMustMany(V_DBUS_SERVICE_NAME, map[string]interface{}{
			"ConnectionState": vpnstatus.ConnectionStateConnected,
			"ConnectedAt":     now,
			"OCRunning":       vpnstatus.OCRunningRunning,
		})
		if iA != nil && a.simulate {
			time.Sleep(time.Second * 5)
			err := iA.ReLogin()
			if err != nil {
				log.Print(err)
			}
		}
	}()
	if false {
		return &dbus.Error{}
	}
	return nil
}

func (a vpnAgent) Disconnect() *dbus.Error {
	if a.simulate {
		counter.Add(1)
		if counter.Load()%3 == 0 {
			counter.Store(0)
			return dbus.MakeFailedError(fmt.Errorf("Disconnect failed"))
		}
	}
	log.Printf("VPN: Disconnect called!\n")
	a.props.SetMust(V_DBUS_SERVICE_NAME, "ConnectionState", vpnstatus.ConnectionStateDisconnecting)
	go func() {
		if a.simulate {
			time.Sleep(time.Second * 5)
		}
		a.props.SetMustMany(V_DBUS_SERVICE_NAME, map[string]interface{}{
			"ConnectionState": vpnstatus.ConnectionStateDisconnected,
			"ConnectedAt":     0,
			"OCRunning":       vpnstatus.OCRunningNotRunning,
		})
		if iA != nil && a.simulate {
			time.Sleep(time.Second * 5)
			iA.props.SetMustMany(I_DBUS_SERVICE_NAME, map[string]interface{}{
				"LoginState":      status.LoginStateLoggedOut,
				"LastKeepAliveAt": 0,
			})
		}
	}()
	return nil
}

func NewVPNServer(simulate bool) *dbus.Conn {
	conn, err := dbus.ConnectSystemBus()
	if err != nil {
		panic(err)
	}

	a := vpnAgent{agent{simulate: simulate}}

	var now int64 = 0
	if a.simulate {
		now = time.Now().Unix()
	}

	// identity properties
	a.props, err = prop.Export(conn, V_DBUS_OBJECT_PATH, prop.Map{
		V_DBUS_SERVICE_NAME: {
			"TrustedNetwork":  {Value: vpnstatus.TrustedNetworkNotTrusted, Writable: false, Emit: prop.EmitTrue, Callback: nil},
			"ConnectionState": {Value: vpnstatus.ConnectionStateUnknown, Writable: false, Emit: prop.EmitTrue, Callback: nil},
			"IP":              {Value: "127.0.0.1", Writable: false, Emit: prop.EmitTrue, Callback: nil},
			"Device":          {Value: "vpn-tun0", Writable: false, Emit: prop.EmitTrue, Callback: nil},
			"ConnectedAt":     {Value: now, Writable: false, Emit: prop.EmitTrue, Callback: nil},
			// "CertExpiresAt":   {Value: now + 60*60*24*365, Writable: false, Emit: prop.EmitTrue, Callback: nil},
			"Servers": {Value: []string{
				"server1.example.com",
				"server2.example.com",
				"server3.example.com",
			}, Writable: false, Emit: prop.EmitTrue, Callback: nil},
			"OCRunning": {Value: vpnstatus.OCRunningNotRunning, Writable: false, Emit: prop.EmitTrue, Callback: nil},
			"VPNConfig": {Value: "", Writable: false, Emit: prop.EmitTrue, Callback: nil},
		},
	})
	if err != nil {
		panic(err)
	}
	// identity methods
	err = conn.Export(a, V_DBUS_OBJECT_PATH, V_DBUS_SERVICE_NAME)
	if err != nil {
		panic(err)
	}
	// vpn introspection
	n := &introspect.Node{
		Name: V_DBUS_OBJECT_PATH,
		Interfaces: []introspect.Interface{
			introspect.IntrospectData,
			prop.IntrospectData,
			{
				Name:       V_DBUS_SERVICE_NAME,
				Methods:    introspect.Methods(a),
				Properties: a.props.Introspection(V_DBUS_SERVICE_NAME),
			},
		},
	}
	err = conn.Export(introspect.NewIntrospectable(n), V_DBUS_OBJECT_PATH, "org.freedesktop.DBus.Introspectable")
	if err != nil {
		log.Panicf("failed to export node introspection: %s\n", err)
	}

	reply, err := conn.RequestName(V_DBUS_SERVICE_NAME,
		dbus.NameFlagDoNotQueue)
	if err != nil {
		panic(err)
	}
	if reply != dbus.RequestNameReplyPrimaryOwner {
		log.Panic("name already taken")
	}
	log.Printf("Listening on interface - %v and path %v ...\n", V_DBUS_SERVICE_NAME, V_DBUS_OBJECT_PATH)

	return conn
}
