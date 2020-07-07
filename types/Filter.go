package types

import (
	"regexp"
	"strings"
	"time"
)

type Filter struct {
	DashId    int      `json:"dash_id"`
	Hostname  string   `json:"hostname"`
	Logname   string   `json:"logname"`
	Level     string   `json:"level"`
	Pid       int      `json:"pid"`
	Version   string   `json:"version"`
	Message   string   `json:"message"`
	Timestamp [2]int64 `json:"timestamp"`
	Pattern   string   `json:"pattern"`
	Offset    int64    `json:"offset"`
	Limit     int      `json:"limit"`
	Keyname   string   `json:"keyname"`
}

func (f *Filter) Match(log *Log) bool {
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
	if f.Timestamp[0] != 0 && log.Timestamp < f.Timestamp[0] {
		return false
	}
	if f.Timestamp[1] != 0 && log.Timestamp > f.Timestamp[1] {
		return false
	}
	if f.Message != "" && !regexp.MustCompile(f.Message).MatchString(log.Message) {
		return false
	}
	if f.Pattern != "" {
		s := strings.Split(f.Pattern, "T")
		dt := time.Unix(0, log.Timestamp).UTC()
		day := dt.Format("2006-01-02")
		tm := dt.Format("15:04:05")
		if !regexp.MustCompile("^" + s[0]).MatchString(day) {
			return false
		}
		if len(s) > 1 && !regexp.MustCompile("^"+s[1]).MatchString(tm) {
			return false
		}
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
		sql += " AND match(formatDateTime(day, '%F', 'UTC'), ?)"
		values = append(values, "^"+s[0])
		if len(s) > 1 {
			sql += " AND match(formatDateTime(toDateTime(timestamp/1e9), '%T', 'UTC'), ?)"
			values = append(values, "^"+s[1])
		}
	}
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
	if f.Message != "" {
		sql += " AND match(message, ?)"
		values = append(values, f.Message)
	}
	return sql, values
}
