package types

import "time"

type Log struct {
	Timestamp time.Time `json:timestamp`
	Host      string    `json:host`
	Project   string    `json:project`
	Logname   string    `json:logname`
	Level     string    `json:level`
	Message   string    `json:message`
	Order     int       `json:order`
}

type Logs []Log
