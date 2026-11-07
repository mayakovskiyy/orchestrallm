package db

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"strconv"
	"strings"
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

func GetDataAK() (string, string) {
	var key string
	var provider string

	err := db.QueryRow("SELECT key, provider FROM api_keys").Scan(&key, &provider)
	if err != sql.ErrNoRows {
		fmt.Println(err)
	} else if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Successfully got the encrypted API key and API provider.")
	}

	return key, provider
}

func UpdateDataAK(key string, provider string, id string) {
	cid, err := strconv.Atoi(strings.TrimSpace(id))
	if err != nil {
		fmt.Println(err)
	}

	db.Exec("UPDATE api_keys SET key = ?, provider = ? WHERE id = ?;", key, provider, cid)
}
