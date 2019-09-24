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
	Bind       string `yaml:"bind"`
	Clickhouse string `yaml:"clickhouse"`
}

var args argsT
var config confT

func ParseArgs() {
	flag.StringVar(&args.Configpath, "config", "./config.yaml", "config path")
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
