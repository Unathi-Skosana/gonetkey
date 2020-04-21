package gonetkey

import (
	"encoding/base64"
	"log"
	"os"

	"github.com/go-ini/ini"
	"github.com/kolo/xmlrpc"
	"github.com/mkideal/cli"
	"golang.org/x/crypto/ssh/terminal"
)

type argT struct {
	User    string `cli:"user"   usage:"Student Number / Username"`
	Config  string `cli:"config" usage:"Loads username/password from file"`
	Retries string `cli:"retries" usage:"Number of connection retries (default=1)" dft:"1"`
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
	defaultFirewallURL = "maties2.sun.ac.za:443/RTAD4-RPC3"
	reconnectionDelay  = 60 * 10
)

// Initialise - comment
func Initialise(username string, password string) InetKey {
	client, err := xmlrpc.NewClient(defaultFirewallURL, nil)

	if err != nil {
		log.Fatal("Connection error : ", err)
	}

	defer client.Close()

	return InetKey{UserName: username, Password: password, FirewallStatus: false, Client: client}
}

// OpenConnection - Comment
func (s InetKey) OpenConnection() {
	log.Println("Opening connection...")
	s.Response = s.invoke("rtad4inetkey_api_open2", "any", 0)
	setConnectionStatus(true)

}

// CloseConnection - Comment
func (s InetKey) CloseConnection() {
	log.Println("Closing connection...")
	s.Response = s.invoke("rtad4inetkey_api_close2", "any", 1)
	setConnectionStatus(false)
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

func (s InetKey) invoke(funcName string, platform string, keepAlive int) *reply {
	req := &request{UserName: s.UserName, UserPwd: s.Password, Platform: platform, KeepAlive: keepAlive}
	res := new(reply)

	err := s.Client.Call(funcName, req, res)

	if err != nil {
		log.Println(err)
	}

	return res
}

// Run - Document this
func Run() {
	os.Exit(cli.Run(new(argT), func(ctx *cli.Context) error {
		argv := ctx.Argv().(*argT)
		userName := argv.User
		password := ""
		config := argv.Config
		retries := argv.Retries

		if config == "" {
			userName, password = loadUserCredentials(config)
		}

		if userName != "" && password == "" {
			var err error

			passwordBytes, err := terminal.ReadPassword(0)

			if err != nil {
				log.Fatal("Reading error : ", err)
			}
			password = string(passwordBytes)
		}

		ctx.String("%s\n", userName)
		ctx.String("%s\n", password)
		ctx.String("%s\n", config)
		ctx.String("%s\n", retries)
		return nil
	}))
}
