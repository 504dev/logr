package types

import (
	"github.com/504dev/kidlog/models/log"
	"time"
)

type Where struct {
	Timestamp [2]time.Time `json:"timestamp"`
	Host      string       `json:"host"`
	Project   string       `json:"project"`
	Logname   string       `json:"logname"`
	Level     string       `json:"level"`
	Message   string       `json:"message"`
}

func (w *Where) Match(log log.Log) bool {
	return true
}

func (w *Where) ToSql() string {
	return ""
}
