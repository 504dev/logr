package types

import (
	"github.com/504dev/logr/cipher"
)

type Log struct {
	DashId    int    `db:"dash_id"   json:"dash_id"`
	Timestamp int64  `db:"timestamp" json:"timestamp"`
	Hostname  string `db:"hostname"  json:"hostname"`
	Logname   string `db:"logname"   json:"logname"`
	Level     string `db:"level"     json:"level"`
	Message   string `db:"message"   json:"message"`
	Pid       int    `db:"pid"       json:"pid"`
	Version   string `db:"version"   json:"version,omitempty"`
	Initiator string `db:"initiator" json:"initiator,omitempty"`
}

type Logs []*Log

func (log *Log) Decrypt(cipherText string, priv string) error {
	return cipher.DecodeAesJson(cipherText, priv, log)
}

func (log *Log) Encrypt(priv string) (string, error) {
	return cipher.EncryptAesJson(log, priv)
}
