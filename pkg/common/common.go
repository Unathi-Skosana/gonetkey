package common

import (
	"encoding/base32"
	"fmt"
	"log"
	"strings"

	"github.com/go-ini/ini"
	"github.com/kolo/xmlrpc"
)

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

const defaultFirewallURL = "https://maties2.sun.ac.za:443/RTAD4-RPC3"

func NewInetkey(username string, password string) InetKey {
	client, err := xmlrpc.NewClient(defaultFirewallURL, nil)

	if err != nil {
		log.Fatal("Failed to initialize client: ", err)
	}

	return InetKey{UserName: username, Password: password, FirewallStatus: false, Client: client}
}

func (s InetKey) OpenConnection() error {
	var err error

	log.Println("Opening connection...")
	s.Response, err = s.Invoke("rtad4inetkey_api_open2", "any", 0)

	if err != nil {
		return err
	}

	message, monthlyUsage, monthlyBytes := s.AccountInfo()

	log.Println(message)
	log.Println(monthlyUsage)
	log.Println(monthlyBytes)

	setConnectionStatus(true)
	return nil

}

func (s InetKey) CloseConnection() error {
	var err error

	log.Println("Closing connection...")
	s.Response, err = s.Invoke("rtad4inetkey_api_close2", "any", 1)

	if err != nil {
		return err
	}

	setConnectionStatus(false)
	return nil
}

func (s InetKey) Invoke(funcName string, platform string, keepAlive int) (*reply, error) {
	req := &request{UserName: s.UserName, UserPwd: s.Password, Platform: platform, KeepAlive: keepAlive}
	res := new(reply)
	err := s.Client.Call(funcName, req, res)

	return res, err
}

func LoadUserCredentials(config string) (string, string) {
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
		encodedPassword := base32.StdEncoding.EncodeToString([]byte(strings.TrimSpace(password)))
		cfg.Section("config").Key("encoded_password_b32").SetValue(encodedPassword)
		cfg.Section("config").Key("password").SetValue("")
		cfg.SaveTo(config)
	} else {
		encodedPassword := cfg.Section("config").Key("encoded_password_b32").String()
		passwordBytes, err = base32.StdEncoding.DecodeString(strings.TrimSpace(encodedPassword))
		if err != nil {
			log.Fatal("decode error:", err)
		}
		password = string(passwordBytes)
	}

	return username, password
}

func (s InetKey) AccountInfo() (string, string, string) {
	code := s.Response.Code
	message := s.Response.Message
	usage := s.Response.Usage
	bytes := s.Response.Bytes

	if code != 0 {
		if strings.Contains(message, "rejected") || strings.Contains(message, "password") {
			log.Fatal(message)
		}
	}

	monthlyUsage := fmt.Sprintf("Monthly usage: R%0.2f", usage)
	monthlyBytes := fmt.Sprintf("Monthly bytes: %d MB", bytes/1024.0/1024.0)

	return message, monthlyUsage, monthlyBytes
}

func setConnectionStatus(status bool) {
	if status {
		log.Println("Connection open. Press <Ctrl> C to close connection")
	} else {
		log.Println("Connection closed.")
	}
}
