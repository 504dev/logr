package log

import (
	"encoding/base64"
	"encoding/json"
	"github.com/504dev/kidlog/cipher"
	"github.com/504dev/kidlog/models/dashboard"
	"time"
)

type Log struct {
	Timestamp time.Time `json:"timestamp"`
	DashId    int       `json:"dashid"`
	Host      string    `json:"host"`
	Project   string    `json:"project"`
	Logname   string    `json:"logname"`
	Level     int       `json:"level"`
	Message   string    `json:"message"`
	Order     int       `json:"order"`
}

type Logs []*Log

type LogPackage struct {
	PublicKey  string `json:"publickey"`
	CipherText string `json:"ciphertext"`
	Signature  string `json:"sig"`
	*Log       `json:"log"`
}

func (lp *LogPackage) Decrypt() error {
	dash := dashboard.GetByPub(lp.PublicKey)
	priv64, _ := base64.StdEncoding.DecodeString(dash.PrivateKey)
	sig64, _ := base64.StdEncoding.DecodeString(lp.Signature)
	cipher64, _ := base64.StdEncoding.DecodeString(lp.CipherText)
	privateKey, err := cipher.BytesToPrivateKey(priv64)
	if err != nil {
		return err
	}
	text, err := cipher.DecryptWithPrivateKey(cipher64, privateKey)
	if err != nil {
		return err
	}
	if err = cipher.CheckSig(text, sig64, &privateKey.PublicKey); err != nil {
		return err
	}
	err = json.Unmarshal(text, &lp.Log)
	if err != nil {
		return err
	}
	lp.DashId = dash.Id
	return nil
}
