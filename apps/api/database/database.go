package database

import (
	"database/sql"
	sqlStore "github.com/hallgren/eventsourcing/eventstore/sql"
)

func CreateDatabase(dbFile string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return nil, err
	}

	es := sqlStore.Open(db)

	rows, err := db.Query("SELECT name FROM sqlite_master WHERE type='table' AND name='events'")
	if err != nil {
		return nil, err
	}

	if !rows.Next() {
		err = es.Migrate()
	}

	if err != nil {
		return nil, err
	}

	err = rows.Close()

	if err != nil {
		return nil, err
	}

	if _, err := db.Exec("CREATE TABLE IF NOT EXISTS ingredients (id VARCHAR(255) NOT NULL PRIMARY KEY, name VARCHAR(255) NOT NULL)"); err != nil {
		return nil, err
	}

	return db, nil
}
