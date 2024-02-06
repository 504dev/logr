package config

import (
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

type Config struct {
	sync.RWMutex `yaml:"-"`
	Bind         struct {
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
	NoDemo        bool   `yaml:"no_demo"`
}

var args Args
var config Config

func ParseArgs() {
	flag.StringVar(&args.Configpath, "config", "./config.yml", "set service config file")
	flag.Parse()
}

func Get() *Config {
	config.RLock()
	defer config.RUnlock()
	return &config
}

func Set(set func(c *Config)) {
	config.Lock()
	clone := config
	set(&clone)
	config = clone
	config.Unlock()
}

func Save() error {
	d, err := yaml.Marshal(&config)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(args.Configpath, d, 0644)
	if err != nil {
		return err
	}
	return nil
}

func Init() {
	ParseArgs()
	yamlFile, err := ioutil.ReadFile(args.Configpath)
	if err != nil {
		log.Fatalf("yamlFile.Get err   #%v ", err)
	}
	yamlFile = []byte(os.ExpandEnv(string(yamlFile)))
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
}
