package log

import (
	"encoding/base64"
	"encoding/json"
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
	Order     int       `json:"order"`
}

type Logs []*Log

type LogPackage struct {
	PublicKey  string `json:"publickey"`
	CipherText string `json:"ciphertext"`
	*Log       `json:"log"`
}

func (lp *LogPackage) Decrypt() error {
	dash := dashboard.GetByPub(lp.PublicKey)
	priv64, _ := base64.StdEncoding.DecodeString(dash.PrivateKey)
	cipher64, _ := base64.StdEncoding.DecodeString(lp.CipherText)
	text, err := cipher.DecryptAes(cipher64, priv64)
	if err != nil {
		return err
	}
	err = json.Unmarshal(text, &lp.Log)
	if err != nil {
		return err
	}
	lp.DashId = dash.Id
	return nil
}

func (lp *LogPackage) Encrypt() error {
	dash := dashboard.GetByPub(lp.PublicKey)
	privateKeyBytes, _ := base64.StdEncoding.DecodeString(dash.PrivateKey)
	jsonMsg, err := json.Marshal(lp.Log)
	if err != nil {
		return err
	}
	cipherBytes, err := cipher.EncryptAes(jsonMsg, privateKeyBytes)
	if err != nil {
		return err
	}
	cipherText := base64.StdEncoding.EncodeToString(cipherBytes)
	lp.CipherText = cipherText
	lp.Log = nil
	return nil
}
