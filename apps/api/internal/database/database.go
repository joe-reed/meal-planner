package database

import (
	"database/sql"
)

func CreateDatabase(dbFile string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbFile)

	if err != nil {
		return nil, err
	}

	return db, nil
}
