package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	_ "embed"
)

var Conn *sql.DB

//go:embed schema.sql
var schema string

func CreateTables() error {
	_, err := Conn.Exec(schema)
	return err
}

func Init(databasePath string) error {
	var err error
	// TODO: implement addition of params in other way
	// e.g cut the path by the ?, and then append params
	Conn, err = sql.Open("sqlite3", databasePath + "?_fk=true")
	return err
}
