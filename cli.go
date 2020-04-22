package gonetkey

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/go-ini/ini"
	"github.com/kolo/xmlrpc"
	"github.com/mkideal/cli"
)

type fn func(int)

type argT struct {
	cli.Helper
	User     string `cli:"user"   usage:"Student number / Username"`
	Password string `pw:"p,password" usage:"Password" prompt:"Password"`
	Config   string `cli:"config" usage:"Loads username and/or password from file"`
	Retries  string `cli:"retries" usage:"Number of connection retries (default=1)" dft:"1"`
}

// InetKey - Comment
type InetKey struct {
	UserName       string
	Password       string
	FirewallStatus bool
	Client         *xmlrpc.Client
	Response       *reply
}

type request struct {
	UserName  string `xmlrpc:"requser"`
	UserPwd   string `xmlrpc:"reqpwd"`
	Platform  string `xmlrpc:"platform"`
	KeepAlive int    `xmlrpc:"keepalive"`
}

type reply struct {
	Message string  `xmlrpc:"resultmsg"`
	Code    int     `xmlrpc:"resultcode"`
	Bytes   int     `xmlrpc:"monthbytes"`
	Usage   float64 `xmlrpc:"monthusage"`
}

const (
	defaultFirewallURL = "https://maties2.sun.ac.za:443/RTAD4-RPC3"
	reconnectionDelay  = 10
)

func accountInfo(r *reply) {
	code := r.Code
	message := r.Message
	usage := r.Usage
	bytes := r.Bytes

	if code != 0 {
		if strings.Contains(message, "rejected") || strings.Contains(message, "password") {
			log.Fatal(message)
		}
	}
	log.Println(message)
	log.Printf("Monthly usage: R%0.2f\n", usage)
	log.Printf("Monthly bytes: %d MB", bytes/1024.0/1024.0)
}

// Initialise - comment
func Initialise(username string, password string) InetKey {
	client, err := xmlrpc.NewClient(defaultFirewallURL, nil)

	if err != nil {
		log.Fatal("Failed to initialize client: ", err)
	}

	return InetKey{UserName: username, Password: password, FirewallStatus: false, Client: client}
}

// OpenConnection - Comment
func (s InetKey) OpenConnection() error {
	var err error

	log.Println("Opening connection...")
	s.Response, err = s.invoke("rtad4inetkey_api_open2", "any", 0)

	if err != nil {
		return err
	}

	accountInfo(s.Response)
	setConnectionStatus(true)
	return nil

}

// CloseConnection - Comment
func (s InetKey) CloseConnection() error {
	var err error

	log.Println("Closing connection...")
	s.Response, err = s.invoke("rtad4inetkey_api_close2", "any", 1)

	if err != nil {
		return err
	}

	accountInfo(s.Response)
	setConnectionStatus(false)
	return nil
}

func loadUserCredentials(config string) (string, string) {
	var err error
	var username string
	var password string
	var passwordBytes []byte
	var cfg *ini.File

	cfg, err = ini.Load(config)
	if err != nil {
		log.Fatal("Fail to read file: ", err)
	}

	username = cfg.Section("config").Key("username").String()
	password = cfg.Section("config").Key("password").String()

	if password != "" {
		encodedPassword := base64.StdEncoding.EncodeToString([]byte(password))
		cfg.Section("config").Key("encoded_password_b32").SetValue(encodedPassword)
		cfg.Section("config").Key("password").SetValue("")
		cfg.SaveTo("config.ini")
	} else {
		encodedPassword := cfg.Section("config").Key("encoded_password_b32").String()
		passwordBytes, err = base64.StdEncoding.DecodeString(encodedPassword)

		if err != nil {
			log.Fatal("decode error:", err)
		}
		password = string(passwordBytes)
	}

	return username, password
}

func setConnectionStatus(status bool) {
	if status {
		log.Println("Connection open. Press <Ctrl> C to close connection")
	} else {
		log.Println("Connection closed.")
	}
}

func (s InetKey) invoke(funcName string, platform string, keepAlive int) (*reply, error) {
	req := &request{UserName: s.UserName, UserPwd: s.Password, Platform: platform, KeepAlive: keepAlive}
	res := new(reply)
	err := s.Client.Call(funcName, req, res)

	return res, err
}

// Validate implements cli.Validator interface
func (argv *argT) Validate(ctx *cli.Context) error {
	userName := argv.User
	password := argv.Password
	config := argv.Config

	if strings.Compare(config, "") != 0 {
		userName, password = loadUserCredentials(config)
	}

	if strings.Compare(userName, "") == 0 || strings.Compare(password, "") == 0 {
		return fmt.Errorf("Username or Password was not provided")
	}

	return nil
}

// Run - Document this
func Run() {
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
			userName, password = loadUserCredentials(config)
		}

		inetkey := Initialise(userName, password)
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
