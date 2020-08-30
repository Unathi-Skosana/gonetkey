package dservice

import (
	"encoding/json"
	"log"
	"os"

	"github.com/godbus/dbus"
	"github.com/godbus/dbus/introspect"
	"github.com/unathi-skosana/gonetkey/internal/pkg/common"
)

type dbusT struct {
	inetkey *common.InetKey
}

type dservice struct {
	serv *dbusT
}

// Open firewall connection
func (s dbusT) Open() *dbus.Error {
	s.inetkey.OpenConnection()
	return nil
}

// Close the firewall connection temporarily - will
// try to reconnect again
func (s dbusT) Close() *dbus.Error {
	s.inetkey.CloseConnection()
	return nil
}

// Get process id
func (s dbusT) GetPID() (int, *dbus.Error) {
	pid := os.Getpid()
	return pid, nil
}

// Get firewall status
func (s dbusT) Status() (string, *dbus.Error) {
	status := s.inetkey.Status
	return status, nil
}

// Get the response as a JSON object
func (s dbusT) JSON_status() ([]byte, *dbus.Error) {
	json, err := json.Marshal(s.inetkey.Response)
	if err != nil {
		panic(err)
	}
	return json, nil
}

// Get response message - rejected or accept
func (s dbusT) Message() (string, *dbus.Error) {
	message, _, _ := s.inetkey.AccountInfo()
	return message, nil
}

// Get monthly usage
func (s dbusT) Usage() (string, *dbus.Error) {
	_, monthlyUsage, _ := s.inetkey.AccountInfo()
	return monthlyUsage, nil
}

// Get monthly bytes
func (s dbusT) Bytes() (string, *dbus.Error) {
	_, _, monthBytes := s.inetkey.AccountInfo()
	return monthBytes, nil
}

// Get username
func (s dbusT) User() (string, *dbus.Error) {
	return s.inetkey.UserName, nil
}

// Run Dbus service
func (d *dservice) Run() {
	go func() {
		conn, err := dbus.SessionBus()
		if err != nil {
			panic(err)
		}

		reply, err := conn.RequestName(common.BusName,
			dbus.NameFlagDoNotQueue)
		if err != nil {
			panic(err)
		}

		if reply != dbus.RequestNameReplyPrimaryOwner {
			log.Fatal("Error : Name already taken")
			os.Exit(1)
		}

		conn.Export(d.serv, common.ObjectPath, common.BusName)
		n := &introspect.Node{
			Name: common.ObjectPath,
			Interfaces: []introspect.Interface{
				introspect.IntrospectData,
				{
					Name:    common.BusName,
					Methods: introspect.Methods(d.serv),
				},
			},
		}
		conn.Export(introspect.NewIntrospectable(n),
			common.ObjectPath,
			"org.freedesktop.DBus.Introspectable")

		log.Printf("Listening on %s / %s\n", common.BusName, common.ObjectPath)

		c := make(chan *dbus.Signal)
		conn.Signal(c)
		for _ = range c {
		}
	}()
}

// Initialize new DBus service object
func NewDbusService(inetkey *common.InetKey) *dservice {
	s := &dbusT{}
	s.inetkey = inetkey
	d := &dservice{}
	d.serv = s
	return d
}
