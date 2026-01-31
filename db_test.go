package gonix

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

var testMemoryDb = "file::memory:?cache=shared"

func TestDBInsert(t *testing.T) {
	var err error
	sqlite, err = NewSQLiteDB(testMemoryDb)
	if err != nil {
		t.Fatal(err)
	}
	defer sqlite.db.Close()

	expectedUser := User{
		UserID:  1,
		Name:    "test",
		Created: time.Now().UTC(),
	}
	err = InsertRecord(sqlite.db, "user", &expectedUser)
	if err != nil {
		t.Fatal(err)
	}

	expectedToken := Token{
		Token:    "027864ad-ab88-4fbc-bedf-998d72bf33ac",
		UserID:   expectedUser.UserID,
		LastUsed: time.Now().UTC(),
		Created:  time.Now().UTC(),
	}
	err = InsertRecord(sqlite.db, "token", &expectedToken)
	if err != nil {
		t.Fatal(err)
	}

	var rawJSON []byte
	err = sqlite.db.QueryRow("SELECT json(data) FROM user WHERE UserID = ?", expectedUser.UserID).Scan(&rawJSON)

	if err != nil {
		t.Fatal(err)
	}

	var gotUser User
	err = json.Unmarshal(rawJSON, &gotUser)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(expectedUser, gotUser) {
		t.Errorf("Structs do not match!\nExpected: %+v\nGot: %+v", expectedUser, gotUser)
	}

	err = sqlite.db.QueryRow("SELECT json(data) FROM token WHERE token = ?", expectedToken.Token).Scan(&rawJSON)

	if err != nil {
		t.Fatal(err)
	}

	var gotToken Token
	err = json.Unmarshal(rawJSON, &gotToken)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(expectedToken, gotToken) {
		t.Errorf("Structs do not match!\nExpected: %+v\nGot: %+v", expectedToken, gotToken)
	}
}
