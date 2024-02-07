package config

import (
	"crypto/sha256"
	"encoding/base64"
	"flag"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

type Args struct {
	Configpath string
}
type ConfigData struct {
	Bind struct {
		Http string `yaml:"http"`
		Udp  string `yaml:"udp"`
	} `yaml:"bind"`
	OAuth struct {
		Github struct {
			Org          string `yaml:"org"           json:"org"`
			ClientId     string `yaml:"client_id"     json:"client_id"`
			ClientSecret string `yaml:"client_secret" json:"-"`
		} `yaml:"github"`
		JwtSecret string `yaml:"jwt_secret"`
	} `yaml:"oauth"`
	Clickhouse    string `yaml:"clickhouse"`
	Mysql         string `yaml:"mysql"`
	AllowNoCipher bool   `yaml:"allow_no_cipher"`
	DemoDash      bool   `yaml:"demo_dash"`
}

func (c *ConfigData) GetJwtSecret() string {
	if c.OAuth.JwtSecret != "" {
		return c.OAuth.JwtSecret
	}
	hash := sha256.Sum256([]byte(c.OAuth.Github.ClientSecret))
	return base64.StdEncoding.EncodeToString(hash[:])
}
func (c *ConfigData) NeedSetup() bool {
	return c.OAuth.Github.ClientId == "" || c.OAuth.Github.ClientSecret == ""
}

type Config struct {
	sync.RWMutex
	Data ConfigData
}

var args Args
var config Config

func ParseArgs() {
	flag.StringVar(&args.Configpath, "config", "./config.yml", "set service config file")
	flag.Parse()
}

func Init() {
	ParseArgs()
	yamlFile, err := ioutil.ReadFile(args.Configpath)
	if err != nil {
		log.Fatalf("yamlFile.Get err   #%v ", err)
	}
	yamlFile = []byte(os.ExpandEnv(string(yamlFile)))
	err = yaml.Unmarshal(yamlFile, &config.Data)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
}

func Get() *ConfigData {
	config.RLock()
	defer config.RUnlock()
	return &config.Data
}

func Set(set func(c *ConfigData)) {
	config.Lock()
	clone := config.Data
	set(&clone)
	config.Data = clone
	config.Unlock()
}

func Save() error {
	d, err := yaml.Marshal(&config.Data)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(args.Configpath, d, 0644)
	if err != nil {
		return err
	}
	return nil
}
