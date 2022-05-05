package db

import (
	"database/sql"
	"encoding/json"
	"fmt"

	_ "github.com/lib/pq"
)

type Conn struct {
	db *sql.DB
}

func OpenConn(host string, port int, dbname string, user string, password string) (*Conn, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=disable", host, port, dbname, user, password))
	if err != nil {
		return nil, err
	}
	return &Conn{
		db: db,
	}, nil
}

func (c *Conn) Ping() error {
	return c.db.Ping()
}

func (c *Conn) Close() error {
	return c.db.Close()
}

func (c *Conn) Write(provider string, json_data interface{}) error {
	sql := "insert into prices (provider, json_data) values ($1, $2)"
	query, err := c.db.Prepare(sql)
	if err != nil {
		return err
	}
	defer query.Close()

	byte_data, err := json.Marshal(json_data)
	if err != nil {
		return err
	}
	_, err = query.Exec(provider, byte_data)

	if err != nil {
		return err
	}
	return nil
}

func (c *Conn) ReadLast(provider string, fsyms string, tsyms string) (*string, error) {
	var json_str string
	query := `SELECT json_data->'RAW'->$1->$2 FROM prices WHERE provider=$3 ORDER BY created DESC LIMIT 1;`

	row := c.db.QueryRow(query, fsyms, tsyms, provider)
	if err := row.Scan(&json_str); err != nil {
		return nil, err
	}

	return &json_str, nil
}
