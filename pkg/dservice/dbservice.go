package dbservice

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/godbus/dbus"
	"github.com/unathi-skosana/gonetkey/pkg/common"
)

type open string
type close string
type pid string
type stop string
type status string
type json_status string
type usage string
type user string

const (
	BusName    string = "za.ac.sun.gonetkey"
	ObjectPath string = "/za/ac/sun/gonetkey/system"
)

const Pkill = "pkill -9 -f \"go.+gonetkey.go\""

func (f open) Open(inetkey common.InetKey) *dbus.Error {
	inetkey.OpenConnection()
	return nil
}

func (f close) Close(inetkey common.InetKey) *dbus.Error {
	inetkey.CloseConnection()
	return nil
}

func (f pid) Getpid(inetkey common.InetKey) (int, *dbus.Error) {
	return os.Getpid(), nil
}

func (f stop) Shutdown(inetkey common.InetKey) *dbus.Error {
	inetkey.CloseConnection()
	return nil
}

func (f status) Status(inetkey common.InetKey) (string, *dbus.Error) {
	var s string = "close"
	if inetkey.FirewallStatus {
		s = "open"
	}
	return s, nil
}

func (f json_status) JSON_status(inetkey common.InetKey) ([]byte, *dbus.Error) {
	json, err := json.Marshal(inetkey.Response)
	if err != nil {
		panic(err)
	}
	return json, nil
}

func (f usage) Usage(inetkey common.InetKey) (string, *dbus.Error) {
	_, monthlyUsage, _ := inetkey.AccountInfo()
	return monthlyUsage, nil
}

func (f user) User(inetkey common.InetKey) (string, *dbus.Error) {
	return inetkey.UserName, nil
}

func NewDbusConn() *dbus.Conn {
	conn, err := dbus.SessionBus()
	if err != nil {
		panic(err)
	}

	reply, err := conn.RequestName(BusName,
		dbus.NameFlagDoNotQueue)
	if err != nil {
		panic(err)
	}

	if reply != dbus.RequestNameReplyPrimaryOwner {
		fmt.Fprintln(os.Stderr, "name already taken")
		os.Exit(1)
	}

	return conn
}

/*
func main() {
		propsSpec := map[string]map[string]*prop.Prop{
		bus_name: {
			"SomeInt": {
				int32(0),
				true,
				prop.EmitTrue,
				func(c *prop.Change) *dbus.Error {
					fmt.Println(c.Name, "changed to", c.Value)
					return nil
				},
			},
		},
	}

	f := open("Bar")
	conn.Export(f, dbus.ObjectPath(object_path), bus_name)
	props := prop.New(conn, dbus.ObjectPath(object_path), propsSpec)
	n := &introspect.Node{
		Name: object_path,
		Interfaces: []introspect.Interface{
			introspect.IntrospectData,
			prop.IntrospectData,
			{
				Name:       bus_name,
				Methods:    introspect.Methods(f),
				Properties: props.Introspection(bus_name),
			},
		},
	}
	conn.Export(introspect.NewIntrospectable(n), dbus.ObjectPath(object_path),
		"org.freedesktop.DBus.Introspectable")
	fmt.Printf("Listening on %s / %s \n", bus_name, object_path)

	c := make(chan *dbus.Signal)
	conn.Signal(c)
	for _ = range c {
	}
}*/
