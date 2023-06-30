package gonix

import (
	"testing"
	"time"
)

var testMemoryDb = "file::memory:?cache=shared"

func TestUser(t *testing.T) {
	var err error
	db, err = NewSQLiteDB(testMemoryDb)
	if err != nil {
		t.Fatal(err)
	}
	defer db.db.Close()

	err = db.SaveUser(&User{ID: 1, Created: time.Now().UTC()})
	if err != nil {
		t.Fatal(err)
	}

}
