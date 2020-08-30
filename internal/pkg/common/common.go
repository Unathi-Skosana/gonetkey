package common

import (
	"encoding/base32"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/go-ini/ini"
	"github.com/kolo/xmlrpc"
)

const (
	BusName    = "za.ac.sun.gonetkey"
	ObjectPath = "/za/ac/sun/gonetkey/system"
)

const ReconnectionDelay = 10

type InetKey struct {
	UserName string
	Password string
	Status   string
	Client   *xmlrpc.Client
	Response *Reply
}

type Request struct {
	UserName  string `xmlrpc:"requser"`
	UserPwd   string `xmlrpc:"reqpwd"`
	Platform  string `xmlrpc:"platform"`
	KeepAlive int    `xmlrpc:"keepalive"`
}

type Reply struct {
	Message string  `xmlrpc:"resultmsg"`
	Code    int     `xmlrpc:"resultcode"`
	Bytes   int     `xmlrpc:"monthbytes"`
	Usage   float64 `xmlrpc:"monthusage"`
}

const defaultFirewallURL = "https://maties2.sun.ac.za:443/RTAD4-RPC3"

func NewInetkey(username string, password string) *InetKey {
	client, err := xmlrpc.NewClient(defaultFirewallURL, nil)

	if err != nil {
		log.Fatal("Failed to initialize client: ", err)
	}

	return &InetKey{UserName: username, Password: password, Status: "uninitialized", Client: client}
}

// Open firewall connection
func (s *InetKey) OpenConnection() error {
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

	s.setConnectionStatus("up")
	return nil

}

// Close firewall connection
func (s *InetKey) CloseConnection() error {
	var err error

	log.Println("Closing connection...")
	s.Response, err = s.Invoke("rtad4inetkey_api_close2", "any", 1)

	if err != nil {
		return err
	}

	s.setConnectionStatus("down")
	return nil
}

// Make a method call over the xmlrpc client
func (s *InetKey) Invoke(funcName string, platform string, keepAlive int) (*Reply, error) {
	req := &Request{UserName: s.UserName, UserPwd: s.Password, Platform: platform, KeepAlive: keepAlive}
	res := new(Reply)
	err := s.Client.Call(funcName, req, res)

	return res, err
}

// Load user credential from a config i.e
// [config]
// username: <username>
// password: <password>
// encoded_password_b32: <encoded_password_b32>
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

// Log account information
// i.e Response message, Monthly usage, bytes, etc
func (s *InetKey) AccountInfo() (string, string, string) {
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

// A goroutine to open a continuous connection
func (s *InetKey) Run(retries int) {
	retries_left := retries

	go func() {
		for {
			if retries_left <= 0 {
				log.Fatal("Error : Exceeded the retries limit.")
			}

			if s.Status == "down" {
				log.Println("Connected closed. Exiting...")
				break
			}

			retries_left -= 1
			conErr := s.OpenConnection()
			if conErr != nil {
				log.Println("Connection failed. Retrying to connect...")
			} else {
				retries_left = retries
			}

			time.Sleep(ReconnectionDelay * time.Second)
		}
	}()
}

// Set connection status connection.
func (s *InetKey) setConnectionStatus(status string) {
	s.Status = status

	if status == "up" {
		log.Println("Connection open. Press <Ctrl> C to close connection")
	} else if status == "down" {
		log.Println("Connection closed.")
	}
}
