package types

import "github.com/504dev/kidlog/cipher"

type Count struct {
	DashId    int    `db:"dash_id"   json:"dash_id"`
	Timestamp int64  `db:"timestamp" json:"timestamp"`
	Hostname  string `db:"hostname"  json:"hostname"`
	Logname   string `db:"logname"   json:"logname"`
	Keyname   string `db:"keyname"   json:"keyname"`
	inc       *Inc
	avg       *Avg
	max       *Max
	min       *Min
	per       *Per
}

type Counts []*Count

func (c *Count) Decrypt(cipherText string, priv string) error {
	return cipher.DecodeAesJson(cipherText, priv, c)
}

func (c *Count) Encrypt(priv string) (string, error) {
	return cipher.EncryptAesJson(c, priv)
}

func (c *Count) Inc(num float64) {
	if c.inc == nil {
		c.inc = &Inc{}
	}
	c.inc.Val += num
}

func (c *Count) Max(num float64) {
	if c.max == nil {
		c.max = &Max{}
	}
	if c.max.Val == nil || num > *c.max.Val {
		c.max.Val = &num
	}
}

func (c *Count) Min(num float64) {
	if c.min == nil {
		c.min = &Min{}
	}
	if c.min.Val == nil || num > *c.min.Val {
		c.min.Val = &num
	}
}

func (c *Count) Avg(num float64) {
	if c.avg == nil {
		c.avg = &Avg{}
	}
	c.avg.Sum += num
	c.avg.Num += 1
}

func (c *Count) Per(taken float64, total float64) {
	if c.per == nil {
		c.per = &Per{}
	}
	c.per.Taken += taken
	c.per.Total += total
}

type Inc struct {
	Val float64 `json:"inc"`
}

type Max struct {
	Val *float64 `json:"max"`
}

type Min struct {
	Val *float64 `json:"min"`
}

type Avg struct {
	Sum float64 `json:"avg_sum"`
	Num int     `json:"avg_num"`
}

type Per struct {
	Taken float64 `json:"per_tkn"`
	Total float64 `json:"per_ttl"`
}

type Valuer interface {
	Value() float64
}

func (i *Inc) Value() float64 {
	return i.Val
}

func (m *Max) Value() float64 {
	return *m.Val
}

func (m *Min) Value() float64 {
	return *m.Val
}

func (a *Avg) Value() float64 {
	return a.Sum / float64(a.Num)
}

func (p *Per) Value() float64 {
	return p.Taken / p.Total * 100
}
