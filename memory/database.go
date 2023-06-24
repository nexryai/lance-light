package memory

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"lance-light/core"
	"log"
)

const (
	dbPath = "./memory.db"

	nftablesLogsTable = `CREATE TABLE IF NOT EXISTS nftablesLogs (
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

	insertNftablesLog = `
		INSERT INTO nftablesLogs (uuid, src, nic, dst, dpt, mac, proto, timestamp)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
)

func InitDatabase() {
	// データベースのオープン
	db, err := sql.Open("sqlite3", dbPath)
	core.ExitOnError(err, "Failed to open database!")
	defer db.Close()

	// テーブルの作成
	_, err = db.Exec(nftablesLogsTable)
	core.ExitOnError(err, "Failed to create table to database!")

	core.MsgDebug("Table OK.")
}

func CreateNftablesLogRecord(src string, nic string, dst string, dpt string, mac string, proto string, timestamp int) {
	db, err := sql.Open("sqlite3", dbPath)
	core.ExitOnError(err, "Failed to open database!")
	defer db.Close()

	eventUUID := uuid.New().String()
	_, err = db.Exec(insertNftablesLog, eventUUID, src, nic, dst, dpt, mac, proto, timestamp)
	core.ExitOnError(err, "Failed to insert a record to database!")
}

func GetNftablesLogRecord(column string, value string) []NftablesRecord {
	db, err := sql.Open("sqlite3", dbPath)
	core.ExitOnError(err, "Failed to open database!")
	defer db.Close()

	query := fmt.Sprintf(`
		SELECT id, uuid, src, nic, dst, dpt, mac, proto, timestamp
		FROM nftablesLogs
		WHERE %s = "%s"
	`, column, value)

	rows, err := db.Query(query)
	core.ExitOnError(err, "Failed to load database!")
	defer rows.Close()

	var records []NftablesRecord
	for rows.Next() {
		var record NftablesRecord
		err := rows.Scan(
			&record.ID,
			&record.EventUUID,
			&record.SrcIP,
			&record.Nic,
			&record.DstIP,
			&record.DstPort,
			&record.Mac,
			&record.Protocol,
			&record.Timestamp)

		core.ExitOnError(err, "Failed to load records!")
		records = append(records, record)
	}

	log.Println("取得したレコード数:", len(records))
	return records
}
