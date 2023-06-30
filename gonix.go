package gonix

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

var db *SQLiteDB

func Run() error {
	tmpdir, err := os.MkdirTemp("", "gonix_dev_*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpdir)

	sqlitefile := filepath.Join(tmpdir, "gonix.db")
	log.Printf("temp db: %s", sqlitefile)

	if db, err = NewSQLiteDB(sqlitefile, LogOptions, LogVersion); err != nil {
		return fmt.Errorf("NewSQLiteDB(%s): %w", sqlitefile, err)
	}
	log.Printf("db: %v", db)
	return nil
}
