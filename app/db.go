package app

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func OpenSQLite3(dbFile string) (*sqlx.DB, error) {
	db, err := sqlx.Open("sqlite3", dbFile)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
