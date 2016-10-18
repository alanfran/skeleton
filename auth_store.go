package main

import (
  "time"
  "errors"
  "crypto/rand"
  "encoding/base64"
  "gopkg.in/pg.v4"
)

type Auth struct {
  UserID     int
  Key     string
  IP      string
  Created time.Time
  Expires time.Time
}

type Oauth struct {
	UserID      int
	Provider string
	Token    string
	Refresh  string
	Expires   time.Time
}

type AuthStore struct {
  mem       map[string]Auth
  db        *pg.DB
  // auth length Expiry time.Time
}

func NewAuthStore(database *pg.DB) *AuthStore {
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

  return &AuthStore{ make(map[string]Auth), database }
}

// input: (uid, ip)
func (s AuthStore) Create(uid int, ip string) (auth Auth, err error) {
  key, err := GenerateAuthKey()
  if err != nil {
    return auth, err
  }

  auth.Key = key
  auth.UserID = uid
  auth.IP = ip
  auth.Created = time.Now()
  auth.Expires = time.Now().AddDate(1, 0, 0)

  err = s.db.Create(&auth)

  return auth, err
}

func (s AuthStore) Get(key string) (auth Auth, err error) {
  err = s.db.Model(&auth).Where("key = ?", key).Select()
  if err != nil {
    return auth, errors.New("Auth not found.")
  }

  return auth, err
}

func (s AuthStore) Del(key string) error {
  var auth Auth
  _, err := s.db.Model(&auth).Where("key = ?", key).Delete()
  if err != nil {
    return errors.New("Error deleting session.")
  }

  return err
}

func GenerateAuthKey() (string, error) {
  b := make([]byte, 512)
  _, err := rand.Read(b)
  if err != nil {
    return "", err
  }
  return base64.URLEncoding.EncodeToString(b), err
}
