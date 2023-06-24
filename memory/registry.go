package memory

import (
	"database/sql"
	"lance-light/core"
)

func SetRegistryValue(key string, value string) {
	InitDatabase()
	//定義
	const (
		insertNftablesLog = `
			INSERT OR REPLACE INTO registry (key, value)
			VALUES (?, ?)
		`
	)

	db, err := sql.Open("sqlite3", dbPath)
	core.ExitOnError(err, "Failed to open database!")
	defer db.Close()

	_, err = db.Exec(insertNftablesLog, key, value)
	core.ExitOnError(err, "Failed to insert a record to database!")
}

func GetRegistryValue(key string) string {
	InitDatabase()

	// クエリ定義
	const (
		selectRegistryValue = `
			SELECT value
			FROM registry
			WHERE key = ?
		`
	)

	// クエリ実行
	db, err := sql.Open("sqlite3", dbPath)
	core.ExitOnError(err, "Failed to open database!")
	defer db.Close()

	row := db.QueryRow(selectRegistryValue, key)
	var value string
	err = row.Scan(&value)
	core.ExitOnError(err, "Failed to load a record to database!")

	return value
}
