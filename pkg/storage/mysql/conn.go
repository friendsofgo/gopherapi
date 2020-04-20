package mysql

import (
	"database/sql"
	"fmt"
)

func NewConn(addr, db string) (*sql.DB, error) {
	conn := fmt.Sprintf("mysql://%s/%s?sslmode=disable", addr, db)
	return sql.Open("mysql", conn)
}
