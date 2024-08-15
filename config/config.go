package config

import (
	"flag"
)

var args CommandLineArgs
var config Config

func parseCommandLineArgs(args *CommandLineArgs) {
	flag.StringVar(&args.Configpath, "config", "./config.yml", "set the logr configuration file")
	flag.IntVar(&args.Retries, "retries", 0, "set the number of attempts to reconnect to databases")
	flag.Parse()
}

func MustLoad() {
	parseCommandLineArgs(&args)
	err := config.ReadFromFile(args.Configpath)
	if err != nil {
		panic(err)
	}
}

func GetCommandLineArgs() CommandLineArgs {
	return args
}

func Get() *ConfigData {
	return config.Get()
}

func Set(set func(c *ConfigData)) {
	config.Set(set)
}

func Save() error {
	return config.Save()
}
