package dservice

import (
	"encoding/json"
	"log"
	"os"

	"github.com/godbus/dbus"
	"github.com/godbus/dbus/introspect"
	"github.com/unathi-skosana/gonetkey/pkg/common"
)

type dbusT struct {
	inetkey *common.InetKey
}

type dservice struct {
	serv *dbusT
}

func (s dbusT) Open() *dbus.Error {
	s.inetkey.OpenConnection()
	return nil
}

func (s dbusT) Close() *dbus.Error {
	s.inetkey.CloseConnection()
	return nil
}

func (s dbusT) GetPID() (int, *dbus.Error) {
	pid := os.Getpid()
	return pid, nil
}

func (s dbusT) Status() (bool, *dbus.Error) {
	status := s.inetkey.FirewallStatus
	return status, nil
}

func (s dbusT) JSON_status() ([]byte, *dbus.Error) {
	json, err := json.Marshal(s.inetkey.Response)
	if err != nil {
		panic(err)
	}
	return json, nil
}

func (s dbusT) Message() (string, *dbus.Error) {
	message, _, _ := s.inetkey.AccountInfo()
	return message, nil
}

func (s dbusT) Usage() (string, *dbus.Error) {
	_, monthlyUsage, _ := s.inetkey.AccountInfo()
	return monthlyUsage, nil
}

func (s dbusT) Bytes() (string, *dbus.Error) {
	_, _, monthBytes := s.inetkey.AccountInfo()
	return monthBytes, nil
}

func (s dbusT) User() (string, *dbus.Error) {
	return s.inetkey.UserName, nil
}

func (d *dservice) Run() {

	go func() {
		conn, err := dbus.ConnectSessionBus()
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

func NewDbusService(inetkey *common.InetKey) *dservice {
	s := &dbusT{}
	s.inetkey = inetkey
	d := &dservice{}
	d.serv = s
	return d
}
