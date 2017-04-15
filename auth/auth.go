package auth

import (
	"crypto/rand"
	"encoding/base64"
	"time"
)

// Auth stores password-based Authentication.
type Auth struct {
	UserID  int
	Key     string
	IP      string
	Created time.Time
	Expires time.Time
}

// Oauth stores OAuth credentials.
type Oauth struct {
	UserID   int
	Provider string
	Token    string
	Refresh  string
	Expires  time.Time
}

// Storer defines the behavior of an Auth store.
type Storer interface {
	Create(userid int, ip string) (Auth, error)
	Get(key string) (Auth, error)
	Del(key string) error
}

// GenerateAuthKey returns a base64-encoded 4096-bit key, or an error.
func GenerateAuthKey() (string, error) {
	b := make([]byte, 512)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), err
}
