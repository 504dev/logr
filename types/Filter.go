package types

import (
	_types "github.com/504dev/logr-go-client/types"
	"regexp"
	"strings"
	"time"
)

type Filter struct {
	DashId    int       `json:"dash_id,omitempty"`
	Hostname  string    `json:"hostname,omitempty"`
	Logname   string    `json:"logname,omitempty"`
	Level     string    `json:"level,omitempty"`
	Pid       int       `json:"pid,omitempty"`
	Version   string    `json:"version,omitempty"`
	Message   string    `json:"message,omitempty"`
	Timestamp *[2]int64 `json:"timestamp,omitempty"`
	Pattern   string    `json:"pattern,omitempty"`
	Offset    int64     `json:"offset,omitempty"`
	Limit     int       `json:"limit,omitempty"`
	Keyname   string    `json:"keyname,omitempty"`
}

func (f *Filter) Match(log *_types.Log) bool {
	if f.DashId != 0 && f.DashId != log.DashId {
		return false
	}
	if f.Hostname != "" && f.Hostname != log.Hostname {
		return false
	}
	if f.Logname != "" && f.Logname != log.Logname {
		return false
	}
	if f.Level != "" && f.Level != log.Level {
		return false
	}
	if f.Version != "" && f.Version != log.Version {
		return false
	}
	if f.Pid != 0 && f.Pid != log.Pid {
		return false
	}
	if f.Timestamp != nil {
		if f.Timestamp[0] != 0 && log.Timestamp < f.Timestamp[0] {
			return false
		}
		if f.Timestamp[1] != 0 && log.Timestamp > f.Timestamp[1] {
			return false
		}
	}
	if f.Message != "" {
		re, err := regexp.Compile(f.Message)
		if err != nil {
			return false
		}
		if !re.MatchString(log.Message) {
			return false
		}
	}
	if f.Pattern != "" {
		s := strings.Split(f.Pattern, "T")
		dt := time.Unix(0, log.Timestamp).UTC()
		day := dt.Format("2006-01-02")
		tm := dt.Format("15:04:05")
		reDay, err := regexp.Compile("^" + s[0])
		if err != nil {
			return false
		}
		if !reDay.MatchString(day) {
			return false
		}
		if len(s) > 1 {
			reTime, err := regexp.Compile("^" + s[1])
			if err != nil {
				return false
			}
			if !reTime.MatchString(tm) {
				return false
			}
		}
	}
	return true
}

func (f *Filter) ToSql() (string, []interface{}, error) {
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
	if f.Level != "" {
		sql += " AND level = ?"
		values = append(values, f.Level)
	}
	if f.Version != "" {
		sql += " AND version = ?"
		values = append(values, f.Version)
	}
	if f.Pid != 0 {
		sql += " AND pid = ?"
		values = append(values, f.Pid)
	}
	if f.Pattern != "" {
		s := strings.Split(f.Pattern, "T")
		r := "^" + s[0]
		if _, err := regexp.Compile(r); err != nil {
			return "", []interface{}{}, err
		}
		sql += " AND match(formatDateTime(day, '%F', 'UTC'), ?)"
		values = append(values, r)
		if len(s) > 1 {
			r := "^" + s[1]
			if _, err := regexp.Compile(r); err != nil {
				return "", []interface{}{}, err
			}
			sql += " AND match(formatDateTime(toDateTime(timestamp/1e9), '%T', 'UTC'), ?)"
			values = append(values, r)
		}
	}
	if f.Timestamp != nil {
		if f.Timestamp[0] != 0 {
			sql += " AND timestamp > ?"
			values = append(values, f.Timestamp[0])
		}
		to := f.Timestamp[1]
		if f.Offset != 0 {
			to = f.Offset
		}
		if to != 0 {
			sql += " AND timestamp < ?"
			values = append(values, to)
		}
	}
	if f.Message != "" {
		if _, err := regexp.Compile(f.Message); err != nil {
			return "", []interface{}{}, err
		}
		sql += " AND match(message, ?)"
		values = append(values, f.Message)
	}
	return sql, values, nil
}
