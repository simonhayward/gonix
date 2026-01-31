package gonix

import (
	"database/sql"
	_ "embed"
	"encoding/json"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

type SQLiteDB struct {
	db *sql.DB
}

//go:embed schema.sql
var sqlSchema string

func NewSQLiteDB(f string, options ...func(*sql.DB)) (*SQLiteDB, error) {
	db, err := openDB(f)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}

	initDb(db)
	for _, option := range options {
		option(db)
	}

	if _, err = db.Exec(sqlSchema); err != nil {
		return nil, err
	}

	return &SQLiteDB{db: db}, nil
}

func openDB(f string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", f)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func initDb(db *sql.DB) {
	for k, v := range map[string]string{
		"foreign_keys": "ON",     // foreign key constraint enforcement
		"busy_timeout": "5000",   // wait for up to 5,000 milliseconds (5 seconds) before failing with a "database is locked"
		"journal_mode": "WAL",    // Write-Ahead Logging mode - allows concurrent reads and writes; readers no longer block writers, and a writer no longer blocks readers
		"synchronous":  "NORMAL", // sync data to the disk at critical moments but not as aggressively as the default
	} {
		_, err := db.Exec(fmt.Sprintf("PRAGMA %s=%s", k, v))

		if err != nil {
			log.Fatalf("%s", err.Error())
		}
	}
}

func LogVersion(db *sql.DB) {
	var version string
	err := db.QueryRow("select sqlite_version()").Scan(&version)

	if err != nil {
		log.Fatalf("%s", err.Error())
	}

	log.Printf("sqlite version: %s", version)
}

func LogOptions(db *sql.DB) {
	for _, opt := range []string{
		"compile_options",
		"synchronous",
		"foreign_keys",
		"journal_mode",
		"incremental_vacuum",
		"busy_timeout",
		"encoding",
		"case_sensitive_like",
	} {
		rows, err := db.Query(fmt.Sprintf("pragma %s", opt))
		if err != nil {
			log.Fatalf("%q: %s", err, opt)
		}
		defer rows.Close()
		for rows.Next() {
			var s string
			err = rows.Scan(&s)
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("%s: %s", opt, s)
		}
	}
}

func InsertRecord[T any](db *sql.DB, table string, data T) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	query := fmt.Sprintf("INSERT INTO %s (data) VALUES (jsonb(?))", table)
	_, err = db.Exec(query, bytes)

	return err
}
