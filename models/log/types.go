package log

import (
	"encoding/base64"
	"encoding/json"
	"errors"
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
	dash := dashboard.GetByPub(lp.PublicKey)
	priv64, _ := base64.StdEncoding.DecodeString(dash.PrivateKey)
	sig64, _ := base64.StdEncoding.DecodeString(lp.Signature)
	text64, _ := base64.StdEncoding.DecodeString(lp.CipherText)
	privateKey, err := cipher.BytesToPrivateKey(priv64)
	if err != nil {
		return err
	}
	if !cipher.CheckSig(sig64, &privateKey.PublicKey) {
		return errors.New("Bad sig")
	}
	text, err := cipher.DecryptWithPrivateKey(text64, privateKey)
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
