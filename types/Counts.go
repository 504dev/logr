package types

import (
	"fmt"
	_types "github.com/504dev/logr-go-client/types"
)

type Counts []*_types.Count

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
