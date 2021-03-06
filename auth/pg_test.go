package auth

import (
	"testing"

	"gopkg.in/pg.v4"
)

var (
	db   *pg.DB
	auth Storer

	dbAddr     = "localhost:5432"
	dbUser     = "postgres"
	dbPassword = "postgres"
	dbDatabase = "test"
)

func init() {
	db = pg.Connect(&pg.Options{
		Addr:     dbAddr,
		User:     dbUser,
		Password: dbPassword,
		Database: dbDatabase,
	})
	// verify connection
	_, err := db.Exec(`SELECT 1`)
	if err != nil {
		panic("Error connecting to the database.")
	}

	auth = NewPgStore(db)
}

func TestAuth(t *testing.T) {
	s, err := auth.Create(1337, "127.0.0.1")
	if err != nil {
		t.Error(err)
	}
	s2, err := auth.Get(s.Key)
	if err != nil {
		t.Error(err)
	}
	if s2.UserID != 1337 || s2.IP != "127.0.0.1" {
		t.Error("Retrieved auth does not match input.")
	}
	err = auth.Del(s2.Key)
	if err != nil {
		t.Error(err)
	}
}
