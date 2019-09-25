package types

import "time"

type Log struct {
	Timestamp time.Time `json:"timestamp"`
	Host      string    `json:"host"`
	Project   string    `json:"project"`
	Logname   string    `json:"logname"`
	Level     string    `json:"level"`
	Message   string    `json:"message"`
	Order     int       `json:"order"`
}

type Logs []Log

type LogPackage struct {
	Public   string `json:"public"`
	LogCrypt string `json:"logcrypt"`
	*Log     `json:"log"`
}

func (pl *LogPackage) Decrypt() error {
	// get private key by public
	// decrypt LogCrypt field
	// set Log field
	return nil
}
