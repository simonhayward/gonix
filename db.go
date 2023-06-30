package gonix

import (
	"database/sql"
	_ "embed"
	"fmt"
	"log"
	"strings"
	"time"

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

	requiredCompiledOptions(db)
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
		"foreign_keys": "ON",
		"busy_timeout": "5000",
		"journal_mode": "WAL",
	} {
		_, err := db.Exec(fmt.Sprintf("PRAGMA %s=%s", k, v))

		if err != nil {
			log.Fatalf(err.Error())
		}
	}
}

func requiredCompiledOptions(db *sql.DB) {
	rows, err := db.Query("pragma compile_options")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var compileOptions string
	for rows.Next() {
		var s string
		err = rows.Scan(&s)
		if err != nil {
			log.Fatal(err.Error())
		}
		compileOptions += fmt.Sprintf(" %s ", s)
	}
	for _, check := range []string{
		"ENABLE_MATH_FUNCTIONS",
		"ENABLE_FTS5",
	} {
		if !strings.Contains(compileOptions, fmt.Sprintf(" %s ", check)) {
			log.Fatalf("compile option: %s required", check)
		}
	}
}

func LogVersion(db *sql.DB) {
	var version string
	err := db.QueryRow("select sqlite_version()").Scan(&version)

	if err != nil {
		log.Fatalf(err.Error())
	}

	log.Printf("version: %s", version)
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

type User struct {
	ID      int
	Created time.Time
}

func (s *SQLiteDB) SaveUser(u *User) error {
	result, err := s.db.Exec("INSERT OR REPLACE INTO Users (ID, Created) VALUES (?, ?)", u.ID, u.Created.Unix())
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows != 1 {
		return fmt.Errorf("expected to affect 1 row, affected %d", rows)
	}
	return nil
}
