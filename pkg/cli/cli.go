package cli

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/mkideal/cli"
	"github.com/unathi-skosana/gonetkey/pkg/common"
	"golang.org/x/crypto/ssh/terminal"
)

const reconnectionDelay = 10

var password string = ""

type argT struct {
	cli.Helper
	Config  string `cli:"config" usage:"Loads username and/or password from file"`
	User    string `cli:"user"   usage:"Student number / Username"`
	Retries string `cli:"retries" usage:"Number of connection retries (default=1)" dft:"1"`
}

func RunCmdClient() {
	os.Exit(cli.Run(new(argT), func(ctx *cli.Context) error {
		var conErr error

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

		inetkey := common.NewInetkey(userName, password)
		retries_left := retries
		c := make(chan os.Signal)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)

		go func() {
			for {
				if retries_left <= 0 {
					log.Fatal("Exceeded the retries limit.")
				}

				retries_left -= 1
				conErr = inetkey.OpenConnection()
				if conErr != nil {
					log.Println("Connection Failed. Retrying connection.")
				} else {
					retries_left = retries
				}

				time.Sleep(reconnectionDelay * time.Second)
			}
		}()

		<-c
		inetkey.CloseConnection()
		os.Exit(1)
		return nil
	}))
}
