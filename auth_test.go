package main

import (
	"testing"

	"gopkg.in/pg.v4"
)

type MockAuthStore struct {
	Mem map[string]Auth
}

func init() {
	if db == nil {
		db = pg.Connect(&pg.Options{
			User:     dbUser,
			Password: dbPassword,
			Database: dbTestDatabase,
		})
		// verify connection
		_, err := db.Exec(`SELECT 1`)
		if err != nil {
			panic("Error connecting to the database.")
		}
	}

	auth = NewAuthStore(db)
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
