package types

import (
	"github.com/504dev/kidlog/models/log"
	"time"
)

type Where struct {
	Timestamp [2]time.Time `json:"timestamp"`
	Project   string       `json:"project"`
	Host      string       `json:"host"`
	Logname   string       `json:"logname"`
	Level     int          `json:"level"`
	Message   string       `json:"message"`
	Offset    time.Time    `json:"offset"`
	Limit     int          `json:"limit"`
}

func (w *Where) Match(log log.Log) bool {
	if w.Project != "" && w.Project != log.Project {
		return false
	}
	if w.Host != "" && w.Host != log.Host {
		return false
	}
	if w.Logname != "" && w.Logname != log.Logname {
		return false
	}
	if w.Level != 0 && w.Level != log.Level {
		return false
	}
	return true
}

func (w *Where) ToSql() string {
	return ""
}
