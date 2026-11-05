package db

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB
var createSQLTable string

func InitDB(dbPath string) {
	var err error

	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		fmt.Println(err)
	}

	db.SetMaxOpenConns(1)
	createSQLTable = `CREATE TABLE IF NOT EXISTS api_keys (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    key TEXT,
    provider TEXT
	);`

	_, err = db.Exec(createSQLTable)
	if err != nil {
		fmt.Println(err)
	}
}

func SaveDataAK(key string, provider string) {
	_, err := db.Exec("INSERT INTO api_keys (key, provider) VALUES (?, ?)", key, provider)
	if err != nil {
		fmt.Println(err)
	}
}
