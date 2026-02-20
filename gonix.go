package gonix

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"time"

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

	mux := http.NewServeMux()
	mux.HandleFunc("/.status", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "ok")
	})

	var wait time.Duration
	addr := fmt.Sprintf(":%s", "8080")

	srv := &http.Server{
		Addr:         addr,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      mux,
	}
	go func() {
		log.Printf("listening on %s\n", addr)
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)

	log.Println("shutting down application")
	return nil
}
