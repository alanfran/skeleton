package user

import (
	"crypto/rand"
	"encoding/base64"
	"time"
)

// User stores user accont information.
type User struct {
	ID       int
	Name     string
	Email    string
	Password string

	Admin bool

	Created   time.Time
	Confirmed bool

	// login lock
	AttemptNumber int64
	AttemptTime   time.Time `sql:",null"`
	Locked        time.Time `sql:",null"`
}

// ConfirmToken stores user email confirmation tokens.
type ConfirmToken struct {
	UserID int ``
	Token  string
}

// RecoverToken stores user account recovery tokens.
type RecoverToken struct {
	UserID  int
	Token   string
	Expires time.Time
}

// Storer defines the behavior of a user store.
type Storer interface {
	Create(u User) (User, error)
	Get(id int) (u User, err error)
	GetByName(n string) (u User, err error)
	Put(u User) error
	Del(id int) error

	Validate(name, pass string) (User, error)

	CreateConfirmationToken(userid int) (ConfirmToken, error)
	ConfirmUser(token string) error

	NewRecover(uid int) (RecoverToken, error)
	RecoverUser(token string) (User, error)
}

// GenerateRandomToken generates a base64-encoded 256-bit key used to
// confirm a user's email address, or to reset their password.
func GenerateRandomToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), err
}
