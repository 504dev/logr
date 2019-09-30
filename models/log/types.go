package log

import (
	"github.com/504dev/kidlog/cipher"
	"github.com/504dev/kidlog/models/dashboard"
	"time"
)

type Log struct {
	DashId    int       `json:"dashid"`
	Timestamp time.Time `json:"timestamp"`
	Host      string    `json:"host"`
	Project   string    `json:"project"`
	Logname   string    `json:"logname"`
	Level     int       `json:"level"`
	Message   string    `json:"message"`
}

type Logs []*Log

type LogPackage struct {
	PublicKey  string `json:"publickey"`
	CipherText string `json:"ciphertext"`
	*Log       `json:"log"`
}

func (lp *LogPackage) DecryptLog() error {
	dash := dashboard.GetByPub(lp.PublicKey)
	log := &Log{}
	err := cipher.DecodeAesJson(lp.CipherText, dash.PrivateKey, log)
	if err != nil {
		return err
	}
	lp.Log = log
	lp.DashId = dash.Id
	return nil
}

func (lp *LogPackage) EncryptLog() error {
	dash := dashboard.GetByPub(lp.PublicKey)
	cipherText, err := cipher.EncryptAesJson(lp.Log, dash.PrivateKey)
	if err != nil {
		return err
	}
	lp.CipherText = cipherText
	return nil
}
