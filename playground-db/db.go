package playDB

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

const dbFile = "playground.db"

type TcpLogItem struct {
	Timestamp string
	ClientIp  string
	Msg       string
}

var tcpLogs []TcpLogItem

func OpenTcpLogsDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbFile)
	if err == nil {
		createTableTcpLogsQuery := `CREATE TABLE IF NOT EXISTS TcpLogs (
		id INTEGER PRIMARY KEY, 
		timestamp TEXT DEFAULT CURRENT_TIMESTAMP, 
		clientIp TEXT, 
		msg TEXT);`
		var res int
		db.QueryRow(createTableTcpLogsQuery).Scan(&res)
		fmt.Println("DB-Setup return: ", res)
	}

	return db, err
}

func InsertTcpLog(db *sql.DB, clientIp string, msg string) (sql.Result, error) {
	insertLogItemWithIpAndMsg := `INSERT INTO TcpLogs (clientIp, msg) VALUES ($1, $2);`
	return db.Exec(insertLogItemWithIpAndMsg, clientIp, msg)
}
