package gonix

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

var sqlite *SQLiteDB

func Run() error {
	sqlitefile, exists := os.LookupEnv("SQLITE_PATH")
	if exists == false || sqlitefile == "" {
		tmpdir, err := os.MkdirTemp("", "gonix_dev_*")
		if err != nil {
			return err
		}
		defer os.RemoveAll(tmpdir)
		sqlitefile = filepath.Join(tmpdir, "gonix.db")
	}

	if _, err := NewSQLiteDB(sqlitefile, LogOptions, LogVersion); err != nil {
		return fmt.Errorf("NewSQLiteDB(%s): %w", sqlitefile, err)
	}
	log.Printf("db: %s", sqlitefile)

	statusHandler := func(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, "ok\n")
	}

	log.Println("Listening on :8080")
	http.HandleFunc("/.status", statusHandler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		return err
	}

	return nil
}
