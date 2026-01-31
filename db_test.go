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

	expectedUser, expectedToken := User{
		UserID:  1,
		Name:    "test",
		Created: time.Now().UTC(),
	}, Token{
		Token:    "027864ad-ab88-4fbc-bedf-998d72bf33ac",
		UserID:   1,
		LastUsed: time.Now().UTC(),
		Created:  time.Now().UTC(),
	}

	err = InsertRecord(sqlite.db, "user", &expectedUser)
	if err != nil {
		t.Fatal(err)
	}
	err = InsertRecord(sqlite.db, "token", &expectedToken)
	if err != nil {
		t.Fatal(err)
	}

	var dataUser, dataToken []byte
	err = sqlite.db.QueryRow(`
		SELECT 
		json(user.data), 
		json(token.data) 
		FROM user 
		LEFT JOIN token ON user.userid = token.userid WHERE user.userid = ?`,
		expectedUser.UserID).Scan(&dataUser, &dataToken)

	if err != nil {
		t.Fatal(err)
	}

	var (
		gotToken Token
		gotUser  User
	)
	json.Unmarshal(dataUser, &gotUser)
	json.Unmarshal(dataToken, &gotToken)

	if !reflect.DeepEqual(expectedToken, gotToken) {
		t.Errorf("Structs do not match!\nExpected: %+v\nGot: %+v", expectedToken, gotToken)
	}
	if !reflect.DeepEqual(expectedUser, gotUser) {
		t.Errorf("Structs do not match!\nExpected: %+v\nGot: %+v", expectedUser, gotUser)
	}
}
