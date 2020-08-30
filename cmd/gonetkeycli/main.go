package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/mkideal/cli"
	"github.com/unathi-skosana/gonetkey/internal/pkg/common"
	"golang.org/x/crypto/ssh/terminal"
)

var password string = ""

type argT struct {
	cli.Helper
	Config  string `cli:"config" usage:"Loads username and/or password from file"`
	User    string `cli:"user"   usage:"Student number / Username"`
	Retries string `cli:"retries" usage:"Number of connection retries (default=1)" dft:"1"`
}

func main() {
	os.Exit(cli.Run(new(argT), func(ctx *cli.Context) error {
		argv := ctx.Argv().(*argT)

		// optional config file containing password and username
		config := argv.Config

		userName := ""
		password := ""

		if strings.Compare(config, "") != 0 {
			// loaded from config
			userName, password = common.LoadUserCredentials(config)
		} else {
			// student number
			userName = argv.User

			// get password
			fmt.Print("Password: ")
			bytePassword, _ := terminal.ReadPassword(0)
			password = string(bytePassword)
		}

		// number of retries
		retries, _ := strconv.Atoi(argv.Retries)

		c := make(chan os.Signal)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)

		inetkey := common.NewInetkey(userName, password)
		inetkey.Run(retries)

		<-c
		inetkey.CloseConnection()
		os.Exit(1)
		return nil
	}))
}
