package gorm

import (
	"database/sql"
	"errors"
	"fmt"
)

type Map = map[string] Value
type List = []Map

type config struct {
	Type string
	DBName string
	User string
	Password string
	Host string
	Port int

}
var configs map[string] config
func SetConfig(cfg map[string] config) {
	configs = cfg
}

func Get(name ...string) *DB {
	c := "default"
	if len(name) > 0 {
		c = name[0]
	}
	config, ok := configs[c]
	if !ok {
		panic(errors.New(c + " database config not found"))
	}

	db, err := Open(config.Type, fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", config.User, config.Password, config.Host, config.Port, config.Name))
	if err != nil {
		panic(errors.New(c + " database config not found"))
	}
	return db
}

func (db *DB)LeftJoin(table, on string, args ...interface{}) *DB {
	db.Joins("LEFT JOIN " + table + " ON " + on, args...)
	return db
}

func (db *DB)RightJoin(table, on string, args ...interface{}) *DB {
	db.Joins("RIGHT JOIN " + table + " ON " + on, args...)
	return db
}
func (db *DB)InnerJoin(table, on string, args ...interface{}) *DB {
	db.Joins("INNER JOIN " + table + " ON " + on, args...)
	return db
}
func (db *DB) GetOne(where ...interface{}) (Map, error) {
	rows, err := db.Limit(1).Rows()
	defer rows.Close()

	if err != nil {
		return nil, err
	}
	if !rows.Next() {
		return nil, nil
	}
	colums, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	clen := len(colums)

	scan := make([]interface{}, clen)
	values := make([]sql.RawBytes, clen)
	for k := range values {
		scan[k] = &values[k]
	}
	row := make(map[string]Value)
	err = rows.Scan(scan...)
	if err != nil {
		return nil, err
	}
	for i, s := range values {
		d := make([]byte, len(s))
		copy(d, s)
		row[colums[i]] = &Var{d:d}
	}

	return row, nil
}



func (s *DB) GetAll(where ...interface{}) (List, error) {
	rows, err := s.Limit(1).Rows()
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	colums, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	clen := len(colums)

	scan := make([]interface{}, clen)
	values := make([]sql.RawBytes, clen)
	for k := range values {
		scan[k] = &values[k]
	}
	list := make(List, 0)
	for rows.Next() {
		row := make(map[string]Value)
		err = rows.Scan(scan...)
		if err != nil {
			return nil, err
		}
		for i, s := range values {
			d := make([]byte, len(s))
			copy(d, s)
			row[colums[i]] = &Var{d: d}
		}
		list = append(list, row)
	}
	return list, err
}

func (s *DB) GetCount() uint64 {
	v := uint64(0)
	s.Count(&v)
	return v
}

func (s *DB)Exist() (bool, error) {
	var any = ""
	err := s.Row().Scan(&any)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
