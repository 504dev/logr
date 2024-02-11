package config

import (
	"flag"
)

var args Args
var config Config

func ParseArgs() {
	flag.StringVar(&args.Configpath, "config", "./config.yml", "set the logr configuration file")
	flag.IntVar(&args.Retries, "retries", 0, "set the number of attempts to reconnect to databases")
	flag.Parse()
}

func Init() Args {
	ParseArgs()
	err := config.FromFile(args.Configpath)
	if err != nil {
		panic(err)
	}
	return args
}

func Get() *ConfigData {
	return config.Get()
}

func Set(set func(c *ConfigData)) {
	config.Set(set)
}

func Save() error {
	return config.ToFile(args.Configpath)
}
