package auth

import (
	"errors"
	"time"

	"gopkg.in/pg.v4"
)

// PgStore stores authentication info.
type PgStore struct {
	mem map[string]Auth
	db  *pg.DB
	// auth length Expiry time.Time
}

// NewPgStore ensures the `auths` table exists in the database and returns an initialized PgStore.
func NewPgStore(database *pg.DB) *PgStore {
	// init tables
	_, err := database.Exec(`CREATE TABLE IF NOT EXISTS auths (
		key text PRIMARY KEY,
		user_id INT NOT NULL,
		ip TEXT,
		created TIMESTAMP,
		expires TIMESTAMP)`)
	if err != nil {
		panic("Error initializing auth table.")
	}
	// _, err := database.Exec(`CREATE TABLE IF NOT EXISTS oauths (...)`)

	return &PgStore{make(map[string]Auth), database}
}

// Create creates an Auth provided with a user ID and IP address, or returns an error.
func (s PgStore) Create(userid int, ip string) (auth Auth, err error) {
	key, err := GenerateAuthKey()
	if err != nil {
		return auth, err
	}

	auth.Key = key
	auth.UserID = userid
	auth.IP = ip
	auth.Created = time.Now()
	auth.Expires = time.Now().AddDate(1, 0, 0)

	err = s.db.Create(&auth)

	return auth, err
}

// Get returns an Auth with the provided key, or an error.
func (s PgStore) Get(key string) (auth Auth, err error) {
	err = s.db.Model(&auth).Where("key = ?", key).Select()
	if err != nil {
		return auth, errors.New("Auth not found.")
	}

	return auth, err
}

// Del deletes the Auth with the provided key, or returns an error.
func (s PgStore) Del(key string) error {
	var auth Auth
	_, err := s.db.Model(&auth).Where("key = ?", key).Delete()
	if err != nil {
		return errors.New("Error deleting session.")
	}

	return err
}
