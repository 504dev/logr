package types

import (
	"fmt"
	"github.com/504dev/logr/cipher"
	"sync"
	"time"
)

type Count struct {
	sync.RWMutex
	DashId    int    `db:"dash_id"   json:"dash_id,omitempty"`
	Timestamp int64  `db:"timestamp" json:"timestamp"`
	Hostname  string `db:"hostname"  json:"hostname,omitempty"`
	Logname   string `db:"logname"   json:"logname,omitempty"`
	Keyname   string `db:"keyname"   json:"keyname"`
	Version   string `db:"version"   json:"version,omitempty"`
	Metrics   `json:"metrics"`
}

type Metrics struct {
	*Inc
	*Max
	*Min
	*Avg
	*Per
	*Time
}

func (m Metrics) ToMap() map[string]interface{} {
	res := map[string]interface{}{}
	if m.Inc != nil {
		res["inc"] = m.Inc.Value()
	}
	if m.Max != nil {
		res["max"] = m.Max.Value()
	}
	if m.Min != nil {
		res["min"] = m.Min.Value()
	}
	if m.Avg != nil {
		res["avg"] = m.Avg.Value()
	}
	if m.Per != nil {
		res["per"] = m.Per.Value()
	}
	if m.Time != nil {
		res["time"] = m.Time.Value()
	}
	return res
}

func (c *Count) Decrypt(cipherText string, priv string) error {
	c.RLock()
	defer c.RUnlock()
	return cipher.DecodeAesJson(cipherText, priv, c)
}

func (c *Count) Encrypt(priv string) (string, error) {
	c.RLock()
	defer c.RUnlock()
	return cipher.EncryptAesJson(c, priv)
}

func (c *Count) AsVector() []interface{} {
	c.RLock()
	defer c.RUnlock()
	dt := time.Unix(c.Timestamp, 0).UTC()
	day := dt.Format("2006-01-02")
	values := []interface{}{day, dt, c.DashId, c.Hostname, c.Logname, c.Keyname, c.Version}
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
		values = append(values, nil, nil)
	} else {
		values = append(values, c.Metrics.Avg.Sum, c.Metrics.Avg.Num)
	}
	if c.Metrics.Per == nil {
		values = append(values, nil, nil)
	} else {
		values = append(values, c.Metrics.Per.Taken, c.Metrics.Per.Total)
	}
	if c.Metrics.Time == nil {
		values = append(values, nil)
	} else {
		values = append(values, c.Metrics.Time.Duration)
	}
	return values
}

func (c *Count) now() {
	c.Timestamp = time.Now().Unix()
}

func (c *Count) Inc(num float64) *Count {
	c.Lock()
	defer c.Unlock()
	if c.Metrics.Inc == nil {
		c.Metrics.Inc = &Inc{}
	}
	c.Metrics.Inc.Val += num
	c.now()
	return c
}

func (c *Count) Max(num float64) *Count {
	c.Lock()
	defer c.Unlock()
	if c.Metrics.Max == nil {
		c.Metrics.Max = &Max{Val: num}
	} else if num > c.Metrics.Max.Val {
		c.Metrics.Max.Val = num
	}
	c.now()
	return c
}

func (c *Count) Min(num float64) *Count {
	c.Lock()
	defer c.Unlock()
	if c.Metrics.Min == nil {
		c.Metrics.Min = &Min{Val: num}
	} else if num < c.Metrics.Min.Val {
		c.Metrics.Min.Val = num
	}
	c.now()
	return c
}

func (c *Count) Avg(num float64) *Count {
	c.Lock()
	defer c.Unlock()
	if c.Metrics.Avg == nil {
		c.Metrics.Avg = &Avg{}
	}
	c.Metrics.Avg.Sum += num
	c.Metrics.Avg.Num += 1
	c.now()
	return c
}

func (c *Count) Per(taken float64, total float64) *Count {
	c.Lock()
	defer c.Unlock()
	if c.Metrics.Per == nil {
		c.Metrics.Per = &Per{}
	}
	c.Metrics.Per.Taken += taken
	c.Metrics.Per.Total += total
	c.now()
	return c
}

func (c *Count) Time(duration time.Duration) func() time.Duration {
	c.Lock()
	defer c.Unlock()
	if c.Metrics.Time == nil {
		c.Metrics.Time = &Time{}
	}
	c.Metrics.Time.Duration += duration.Nanoseconds()
	c.now()
	ts := time.Now()
	var delta *time.Duration
	return func() time.Duration {
		if delta == nil {
			tmp := time.Since(ts)
			delta = &tmp
			num := float64(tmp.Nanoseconds()) / float64(duration.Nanoseconds())
			c.Avg(num).Min(num).Max(num)
		}
		return *delta
	}
}

type Counts []*Count

type Serie struct {
	Hostname string           `json:"hostname"`
	Keyname  string           `json:"keyname"`
	Kind     string           `json:"kind"`
	Data     [][2]interface{} `json:"data"`
}
type Series []*Serie

func (cs Counts) Format() Series {
	m := map[string]*Serie{}
	for _, c := range cs {
		for k, v := range c.ToMap() {
			key := fmt.Sprintf("%v:%v:%v", k, c.Keyname, c.Hostname)
			if _, ok := m[key]; !ok {
				m[key] = &Serie{Hostname: c.Hostname, Keyname: c.Keyname, Kind: k}
			}
			m[key].Data = append(m[key].Data, [2]interface{}{c.Timestamp, v})
		}
	}
	res := make(Series, 0, len(m))
	for _, s := range m {
		res = append(res, s)
	}
	return res
}

type Inc struct {
	Val float64 `db:"inc,omitempty" json:"inc,omitempty"`
}

type Max struct {
	Val float64 `db:"max,omitempty" json:"max,omitempty"`
}

type Min struct {
	Val float64 `db:"min,omitempty" json:"min,omitempty"`
}

type Avg struct {
	Sum float64 `db:"avg_sum" json:"avg_sum,omitempty"`
	Num int     `db:"avg_num" json:"avg_num,omitempty"`
}

type Per struct {
	Taken float64 `db:"per_tkn" json:"per_tkn,omitempty"`
	Total float64 `db:"per_ttl" json:"per_ttl,omitempty"`
}

type Time struct {
	Duration int64 `db:"time_dur" json:"time_dur,omitempty"`
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

func (t *Time) Value() int64 {
	return t.Duration
}
