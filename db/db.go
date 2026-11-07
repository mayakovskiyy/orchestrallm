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

func SaveDataAK(key string, provider string) (string, error) {
	result, err := db.Exec("INSERT INTO api_keys (key, provider) VALUES (?, ?)", key, provider)
	if err != nil {
		return "", err
	}
	id, err := result.LastInsertId()
	return strconv.FormatInt(id, 10), err
}

type ApiKey struct {
	ID   string
	Key  string
	Name string
}

func GetDataAK() ([]ApiKey, error) {
	rows, err := db.Query("SELECT id, key, provider FROM api_keys ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	aks := []ApiKey{}
	for rows.Next() {
		var ak ApiKey
		if err := rows.Scan(&ak.ID, &ak.Key, &ak.Name); err != nil {
			return nil, err
		}
		aks = append(aks, ak)
	}
	return aks, rows.Err()
}

func UpdateDataAK(key string, provider string, id string) {
	cid, err := strconv.Atoi(strings.TrimSpace(id))
	if err != nil {
		fmt.Println(err)
	}

	db.Exec("UPDATE api_keys SET key = ?, provider = ? WHERE id = ?;", key, provider, cid)
}

func DeleteDataAK(id string) error {
	result, err := db.Exec("DELETE FROM api_keys WHERE id = ?;", id)
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return sql.ErrNoRows
	}
	return nil
}
