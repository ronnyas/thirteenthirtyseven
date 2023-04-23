package game

import (
	"database/sql"
	"fmt"
	"log"
)

func SetDatabase(db *sql.DB) {
	Config.db = db
}

func GetServers() ([]string, error) {
	query := `SELECT DISTINCT serverid FROM config`
	rows, err := Config.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying database: %s", err)
	}
	defer rows.Close()

	var serverIDs []string
	for rows.Next() {
		var serverID string
		if err := rows.Scan(&serverID); err != nil {
			log.Println("Error scanning row:", err)
			continue
		}
		serverIDs = append(serverIDs, serverID)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %s", err)
	}
	return serverIDs, nil
}
