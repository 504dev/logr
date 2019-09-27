package log

import (
	"encoding/json"
	"github.com/504dev/kidlog/cipher"
	"github.com/504dev/kidlog/models/dashboard"
	"time"
)

type Log struct {
	Timestamp time.Time `json:"timestamp"`
	Dashboard int
	Host      string `json:"host"`
	Project   string `json:"project"`
	Logname   string `json:"logname"`
	Level     int    `json:"level"`
	Message   string `json:"message"`
	Order     int    `json:"order"`
}

type Logs []Log

type LogPackage struct {
	PublicKey  string `json:"publickey"`
	CipherText string `json:"ciphertext"`
	Signature  string `json:"sig"`
	Log        `json:"log"`
}

func (lp *LogPackage) Decrypt() error {
	dash := dashboard.GetByPk(lp.PublicKey)
	privateKey, err := cipher.BytesToPrivateKey([]byte(dash.PrivateKey))
	text, err := cipher.DecryptWithPrivateKey([]byte(lp.CipherText), privateKey)
	if err != nil {
		return err
	}
	var log Log
	err = json.Unmarshal(text, &log)
	if err != nil {
		return err
	}
	return nil
}
