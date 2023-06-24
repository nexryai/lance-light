package memory

import (
	"database/sql"
	"lance-light/core"
)

const (
	dbPath = "./memory.db"
)

var isDatabaseInitialized bool

// InitDatabase データベースが初期化されていない場合初期化する。一度実行したら二回目以降はスキップする。
func InitDatabase() {
	if !isDatabaseInitialized {
		// 各種定義
		const (
			// ToDo なんかもっといい感じに書き直す（そもそもconstする必要性ない？）
			nftablesLogsTable = `
				CREATE TABLE IF NOT EXISTS nftablesLogs (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				uuid TEXT,
				src TEXT,
				nic TEXT,
				dst TEXT,
				dpt TEXT,
				mac TEXT,
				proto TEXT,
				timestamp INTEGER
			);`

			registryTable = `
				CREATE TABLE IF NOT EXISTS registry (
				key TEXT PRIMARY KEY,
				value TEXT
			);`
		)

		// データベースのオープン
		db, err := sql.Open("sqlite3", dbPath)
		core.ExitOnError(err, "Failed to open database!")
		defer db.Close()

		// テーブルの作成
		_, err = db.Exec(nftablesLogsTable)
		core.ExitOnError(err, "Failed to create table to database!")

		_, err = db.Exec(registryTable)
		core.ExitOnError(err, "Failed to create table to database!")

		core.MsgDebug("Table OK.")
		isDatabaseInitialized = true
	} else {
		core.MsgDebug("Skip InitDatabase!")
	}
}
