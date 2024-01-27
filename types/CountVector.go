package types

import (
	_types "github.com/504dev/logr-go-client/types"
	"time"
)

type CountVector _types.Count

func (c *CountVector) AsVector() []interface{} {
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
