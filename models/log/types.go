package log

import (
	"github.com/504dev/kidlog/cipher"
	"github.com/504dev/kidlog/models/dashboard"
)

type Log struct {
	DashId    int    `db:"dash_id"   json:"dash_id"`
	Timestamp int64  `db:"timestamp" json:"timestamp"`
	Hostname  string `db:"hostname"  json:"hostname"`
	Logname   string `db:"logname"   json:"logname"`
	Level     int    `db:"level"     json:"level"`
	Message   string `db:"message"   json:"message"`
}

type Logs []*Log

type LogPackage struct {
	PublicKey  string `json:"publickey"`
	CipherText string `json:"ciphertext"`
	*Log       `json:"log"`
}

func (lp *LogPackage) DecryptLog() error {
	dash, _ := dashboard.GetByPub(lp.PublicKey)
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
	dash, _ := dashboard.GetByPub(lp.PublicKey)
	cipherText, err := cipher.EncryptAesJson(lp.Log, dash.PrivateKey)
	if err != nil {
		return err
	}
	lp.CipherText = cipherText
	return nil
}

type Filter struct {
	Timestamp [2]int64 `json:"timestamp"`
	DashId    int      `json:"dash_id"`
	Hostname  string   `json:"hostname"`
	Logname   string   `json:"logname"`
	Level     int      `json:"level"`
	Message   string   `json:"message"`
	Offset    int64    `json:"offset"`
	Limit     int      `json:"limit"`
}

func (f *Filter) Match(log Log) bool {
	if f.Hostname != "" && f.Hostname != log.Hostname {
		return false
	}
	if f.Logname != "" && f.Logname != log.Logname {
		return false
	}
	if f.Level != 0 && f.Level != log.Level {
		return false
	}
	return true
}

func (f *Filter) ToSql() (string, []interface{}) {
	sql := "where dash_id = ?"
	values := []interface{}{f.DashId}
	if f.Hostname != "" {
		sql += " AND hostname = ?"
		values = append(values, f.Hostname)
	}
	if f.Logname != "" {
		sql += " AND logname = ?"
		values = append(values, f.Logname)
	}
	if f.Level != 0 {
		sql += "level = ?"
		values = append(values, f.Level)
	}
	return sql, values
}
