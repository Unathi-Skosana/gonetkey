package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/godbus/dbus"
	"github.com/mkideal/cli"
	"github.com/unathi-skosana/gonetkey/pkg/common"
)

type argT struct {
	cli.Helper
	Open    bool `cli:"open" usage:"Open firewall" dft:"false"`
	Close   bool `cli:"close" usage:"Close firewall" dft:"false"`
	Pid     bool `cli:"pid" usage:"Get pid" dft:"false"`
	Status  bool `cli:"status"   usage:"Get boolean status" dft:"false"`
	Json    bool `cli:"json"   usage:"Get JSON status" dft:"false"`
	Message bool `cli:"message"   usage:"Get message from client" dft:"false"`
	Usage   bool `cli:"usage"   usage:"Get monthly usage" dft:"false"`
	Bytes   bool `cli:"bytes"   usage:"Get monthly bytes" dft:"false"`
	User    bool `cli:"user"   usage:"Get current user" dft:"false"`
	Kill    bool `cli:"kill"   usage:"Kill command" dft:"false"`
}

func main() {
	os.Exit(cli.Run(new(argT), func(ctx *cli.Context) error {
		conn, err := dbus.ConnectSessionBus()
		if err != nil {
			panic(err)
		}

		defer conn.Close()

		obj := conn.Object(common.BusName, common.ObjectPath)
		argv := ctx.Argv().(*argT)

		open := argv.Open
		close := argv.Close
		pid := argv.Pid
		status := argv.Status
		jstatus := argv.Json
		message := argv.Message
		usage := argv.Usage
		bytes := argv.Bytes
		user := argv.User
		kill := argv.Kill

		if open {
			obj.Call(fmt.Sprintf("%s.Open", common.BusName), 0)
		} else if close {
			obj.Call(fmt.Sprintf("%s.Close", common.BusName), 0)
		} else if pid {
			var i int
			obj.Call(fmt.Sprintf("%s.GetPID", common.BusName), 0).Store(&i)
			fmt.Printf("PID: %d\n", i)
		} else if status {
			var b bool
			obj.Call(fmt.Sprintf("%s.Status", common.BusName), 0).Store(&b)
			fmt.Printf("Firewall status: %v\n", b)
		} else if jstatus {
			var j []byte
			var res *common.Reply

			obj.Call(fmt.Sprintf("%s.JSON_status", common.BusName), 0).Store(&j)
			json.MarshalIndent(j, "", "    ")

			if err := json.Unmarshal(j, res); err != nil {
				panic(err)
			}
			fmt.Println(res)
		} else if message {
			var s string
			obj.Call(fmt.Sprintf("%s.Message", common.BusName), 0).Store(&s)
		} else if usage {
			var s string
			obj.Call(fmt.Sprintf("%s.Usage", common.BusName), 0).Store(&s)
		} else if bytes {
			var s string
			obj.Call(fmt.Sprintf("%s.Bytes", common.BusName), 0).Store(&s)
		} else if user {
			var s string
			obj.Call(fmt.Sprintf("%s.User", common.BusName), 0).Store(&s)
		} else if kill {
			var i int
			obj.Call(fmt.Sprintf("%s.GetPID", common.BusName), 0).Store(&i)
			exec.Command(fmt.Sprintf("kill -9 %d", i))
		} else {
			fmt.Println("See usage: gonetkeyd -h")
		}

		return nil
	}))
}
