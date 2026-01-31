package gonix

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

var sqlite *SQLiteDB

func Run() error {
	tmpdir, err := os.MkdirTemp("", "gonix_dev_*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpdir)

	sqlitefile := filepath.Join(tmpdir, "gonix.db")

	if sqlite, err = NewSQLiteDB(sqlitefile, LogOptions, LogVersion); err != nil {
		return fmt.Errorf("NewSQLiteDB(%s): %w", sqlitefile, err)
	}
	log.Printf("temp db: %s. ctrl-d to exit", sqlitefile)

	// Blocks and waits for Ctrl+D (EOF)
	_, _ = io.ReadAll(os.Stdin)
	return nil
}
