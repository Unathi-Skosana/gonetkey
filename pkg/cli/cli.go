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
)

const reconnectionDelay = 10

type argT struct {
	cli.Helper
	User     string `cli:"user"   usage:"Student number / Username"`
	Password string `pw:"p,password" usage:"Password" prompt:"Password"`
	Config   string `cli:"config" usage:"Loads username and/or password from file"`
	Retries  string `cli:"retries" usage:"Number of connection retries (default=1)" dft:"1"`
}

func (argv *argT) Validate(ctx *cli.Context) error {
	userName := argv.User
	password := argv.Password
	config := argv.Config

	if strings.Compare(config, "") != 0 {
		userName, password = common.LoadUserCredentials(config)
	}

	if strings.Compare(userName, "") == 0 || strings.Compare(password, "") == 0 {
		return fmt.Errorf("Username or Password was not provided")
	}

	return nil
}

func RunCmdClient() {
	os.Exit(cli.Run(new(argT), func(ctx *cli.Context) error {
		var conErr error

		argv := ctx.Argv().(*argT)

		// student number
		userName := argv.User

		// password
		password := argv.Password

		// optional config file containing password and username
		config := argv.Config

		// number of retries
		retries, _ := strconv.Atoi(argv.Retries)

		if strings.Compare(config, "") != 0 {
			userName, password = common.LoadUserCredentials(config)
		}

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
