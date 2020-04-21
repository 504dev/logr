package mysql

import (
	"fmt"
)

type Model struct {
	sql string
}

type Sampler interface {
	First() interface{}
	Len() int
}

func (m *Model) FindAllByField(result interface{}, fieldname string, val interface{}, limit int) error {
	conn := Conn()
	sqltext := fmt.Sprintf("%v WHERE %v = ?", m.sql, fieldname)
	if limit > 0 {
		sqltext = fmt.Sprintf("%v LIMIT %v", sqltext, limit)
	}
	fmt.Println(sqltext, val)
	return conn.Select(result, sqltext, val)
}

func CreateModel(sql string) *Model {
	return &Model{sql}
}
