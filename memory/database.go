package memory

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"lance-light/core"
)

func CreateNftablesLogRecord(src string, nic string, dst string, dpt string, mac string, proto string, timestamp int) {
	//定義
	const (
		insertNftablesLog = `
			INSERT INTO nftablesLogs (uuid, src, nic, dst, dpt, mac, proto, timestamp)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		`
	)

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

	// クエリの作成
	query := fmt.Sprintf(`
		SELECT uuid, src, nic, dst, dpt, mac, proto, timestamp
		FROM nftablesLogs
		WHERE %s = "%s"
	`, column, value)

	// クエリ実行
	rows, err := db.Query(query)
	core.ExitOnError(err, "Failed to load database!")
	defer rows.Close()

	// マッピングする
	var records []NftablesRecord
	for rows.Next() {
		var record NftablesRecord
		err := rows.Scan(
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

	return records
}
