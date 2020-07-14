package config

import (
	"flag"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

type argsT struct {
	Configpath string
}

type confT struct {
	Bind struct {
		Http string `yaml:"http"`
		Udp  string `yaml:"udp"`
	} `yaml:"bind"`
	OAuth struct {
		Github struct {
			ClientId     string `yaml:"client_id"`
			ClientSecret string `yaml:"client_secret"`
			Org          string `yaml:"org"`
		} `yaml:"github"`
		JwtSecret string `yaml:"jwt_secret"`
	} `yaml:"oauth"`
	Clickhouse string `yaml:"clickhouse"`
	Mysql      string `yaml:"mysql"`
}

var args argsT
var config confT

func ParseArgs() {
	flag.StringVar(&args.Configpath, "config", "./config.yml", "set service config file")
	flag.Parse()
}

func Args() *argsT {
	return &args
}
func Get() *confT {
	return &config
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
