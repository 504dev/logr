package types

import (
	"github.com/504dev/kidlog/cipher"
	"time"
)

type Count struct {
	DashId    int    `json:"dash_id"`
	Timestamp int64  `json:"timestamp"`
	Hostname  string `json:"hostname"`
	Logname   string `json:"logname"`
	Keyname   string `json:"keyname"`
	Metrics   struct {
		*Inc
		*Avg
		*Max
		*Min
		*Per
	} `json:"metrics"`
}

type Counts []*Count

func (c *Count) Decrypt(cipherText string, priv string) error {
	return cipher.DecodeAesJson(cipherText, priv, c)
}

func (c *Count) Encrypt(priv string) (string, error) {
	return cipher.EncryptAesJson(c, priv)
}

func (c *Count) AsVector() []interface{} {
	dt := time.Unix(0, c.Timestamp)
	day := dt.Format("2006-01-02")
	values := []interface{}{day, dt, c.DashId, c.Hostname, c.Logname, c.Keyname}
	if c.Metrics.Inc == nil {
		values = append(values, nil)
	} else {
		values = append(values, c.Metrics.Inc.Val)
	}
	if c.Metrics.Max == nil {
		values = append(values, nil)
	} else {
		values = append(values, c.Metrics.Max.Val)
	}
	if c.Metrics.Min == nil {
		values = append(values, nil)
	} else {
		values = append(values, c.Metrics.Min.Val)
	}
	if c.Metrics.Avg == nil {
		values = append(values, 0.0, 0)
	} else {
		values = append(values, c.Metrics.Avg.Sum, c.Metrics.Avg.Num)
	}
	if c.Metrics.Per == nil {
		values = append(values, 0.0, 0.0)
	} else {
		values = append(values, c.Metrics.Per.Taken, c.Metrics.Per.Total)
	}
	return values
}

func (c *Count) now() {
	c.Timestamp = time.Now().UnixNano()
}

func (c *Count) Inc(num float64) *Count {
	if c.Metrics.Inc == nil {
		c.Metrics.Inc = &Inc{}
	}
	c.Metrics.Inc.Val += num
	c.now()
	return c
}

func (c *Count) Max(num float64) *Count {
	if c.Metrics.Max == nil {
		c.Metrics.Max = &Max{Val: num}
	} else if num > c.Metrics.Max.Val {
		c.Metrics.Max.Val = num
	}
	c.now()
	return c
}

func (c *Count) Min(num float64) *Count {
	if c.Metrics.Min == nil {
		c.Metrics.Min = &Min{Val: num}
	} else if num < c.Metrics.Min.Val {
		c.Metrics.Min.Val = num
	}
	c.now()
	return c
}

func (c *Count) Avg(num float64) *Count {
	if c.Metrics.Avg == nil {
		c.Metrics.Avg = &Avg{}
	}
	c.Metrics.Avg.Sum += num
	c.Metrics.Avg.Num += 1
	c.now()
	return c
}

func (c *Count) Per(taken float64, total float64) *Count {
	if c.Metrics.Per == nil {
		c.Metrics.Per = &Per{}
	}
	c.Metrics.Per.Taken += taken
	c.Metrics.Per.Total += total
	c.now()
	return c
}

type Inc struct {
	Val float64 `json:"inc"`
}

type Max struct {
	Val float64 `json:"max"`
}

type Min struct {
	Val float64 `json:"min"`
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
	return m.Val
}

func (m *Min) Value() float64 {
	return m.Val
}

func (a *Avg) Value() float64 {
	return a.Sum / float64(a.Num)
}

func (p *Per) Value() float64 {
	return p.Taken / p.Total * 100
}
