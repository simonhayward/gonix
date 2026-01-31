package gonix

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

var testMemoryDb = "file::memory:?cache=shared"

func TestDb(t *testing.T) {
	var err error
	sqlite, err = NewSQLiteDB(testMemoryDb)
	if err != nil {
		t.Fatal(err)
	}
	defer sqlite.db.Close()

	user := User{
		UserID:  1,
		Name:    "test",
		Created: time.Now().UTC(),
	}
	err = InsertRecord(sqlite.db, "user", &user)
	if err != nil {
		t.Fatal(err)
	}

	token := Token{
		Token:    "027864ad-ab88-4fbc-bedf-998d72bf33ac",
		UserID:   user.UserID,
		LastUsed: time.Now().UTC(),
		Created:  time.Now().UTC(),
	}
	err = InsertRecord(sqlite.db, "token", &token)
	if err != nil {
		t.Fatal(err)
	}

	var rawJSON []byte
	err = sqlite.db.QueryRow("SELECT json(data) FROM user WHERE UserID = ?", user.UserID).Scan(&rawJSON)

	if err != nil {
		t.Fatal(err)
	}

	var u User
	err = json.Unmarshal(rawJSON, &u)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(user, u) {
		t.Errorf("Structs do not match!\nExpected: %+v\nGot: %+v", user, u)
	}

}
