package cockroach

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func NewConn(addr, db string) (*sql.DB, error) {
	conn := fmt.Sprintf("postgresql://%s/%s?sslmode=disable", addr, db)
	return sql.Open("postgres", conn)
}
