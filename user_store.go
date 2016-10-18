package main

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gopkg.in/pg.v4"
	//"fmt"
	//"strconv"
)

var (
	bcryptCost = 13
	errNyi     = errors.New("Not yet implemented.")
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

// UserStore stores a reference to the database and the Mailer and provides user CRUD methods.
type UserStore struct {
	db     *pg.DB
	mailer *Mailer
}

// NewUserStore ensures the users, confirm_tokens, and recover_tokens tables exist.
// Then it returns an initialized UserStore.
func NewUserStore(database *pg.DB, mailer *Mailer) *UserStore {
	// initialize database
	_, err := database.Exec(`DROP TABLE users; CREATE TABLE IF NOT EXISTS users (
    id               serial PRIMARY KEY,
    name             text NOT NULL,
    email            text NOT NULL,
    password         text NOT NULL,
    admin            boolean,
    created          timestamp NOT NULL,
    confirmed        boolean NOT NULL,
    attempt_number   bigint,
    attempt_time     timestamp,
    locked           timestamp )`)
	if err != nil {
		panic("Error initializing user table.")
	}

	_, err = database.Exec(`CREATE TABLE IF NOT EXISTS confirm_tokens(
    user_id   int   UNIQUE NOT NULL REFERENCES users (id),
    token     text  PRIMARY KEY )`)
	if err != nil {
		panic("Error initializing confirmation token table.")
	}

	_, err = database.Exec(`CREATE TABLE IF NOT EXISTS recover_tokens(
    user_id   int         UNIQUE NOT NULL REFERENCES users (id),
    token     text        PRIMARY KEY,
    expires   timestamp   NOT NULL )`)
	if err != nil {
		panic("Error initializing recovery token table.")
	}

	// initialize the memory store
	return &UserStore{database, mailer}
}

// Create inserts the provided User into the database.
func (s UserStore) Create(u User) (User, error) {
	// name unique
	_, err := s.GetByName(u.Name)
	if err == nil { // successful query, err != ErrNoRows
		return u, errors.New("That name is taken.")
	}

	// email address exists and is unique
	if !(strings.Contains(u.Email, "@") && strings.Contains(u.Email, ".")) {
		return u, errors.New("Please enter a valid email address.")
	}
	_, err = s.GetByName(u.Email)
	if err == nil {
		return u, errors.New("That email address is already being used.")
	}

	// password > 6 char
	if len(u.Password) < 6 {
		return u, errors.New("Password must have 6 or more characters.")
	}

	// securely hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcryptCost)
	if err != nil {
		return u, errors.New("Error hashing password.")
	}
	u.Password = string(hash)

	u.Created = time.Now()

	// insert into db
	err = s.db.Create(&u)
	if err != nil {
		return u, err
	}

	// generate confirmation token
	confirm := ConfirmToken{UserID: u.ID}
	confirm.Token, err = GenerateRandomToken()
	if err != nil {
		s.db.Delete(&u)
		return u, errors.New("Error generating confirmation token.")
	}
	err = s.db.Create(&confirm)
	if err != nil {
		s.db.Delete(&u)
		return u, err
	}

	// email
	s.mailer.SendConfirmation(u.Email, confirm.Token)
	return u, err
}

// Put updates the provided user in the database.
func (s UserStore) Put(u User) error {
	err := s.db.Update(&u)

	return err
}

// Get returns the user with the provided ID, or an error.
func (s UserStore) Get(id int) (u User, err error) {
	u.ID = id
	err = s.db.Select(&u)

	return u, err
}

// GetByName returns the User with the provided name/email, or an error.
func (s UserStore) GetByName(n string) (u User, err error) {
	err = s.db.Model(&u).Where("name = ? OR email = ?", n, n).Select()

	return u, err
}

// Del deletes the user with the provided ID from the database, or returns an error.
func (s UserStore) Del(id int) error {
	var u User
	u.ID = id
	err := s.db.Delete(&u)
	if err != nil {
		return err
	}
	// delete tokens associated with uid
	s.db.Model(&ConfirmToken{}).Where("user_id = ?", id).Delete()
	s.db.Model(&RecoverToken{}).Where("user_id = ?", id).Delete()
	return err
}

// Validate checks if the provided name/email and password are correct.
func (s UserStore) Validate(name, pass string) (User, error) {
	u, err := s.GetByName(name)
	if err != nil {
		return u, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(pass))
	if err != nil {
		// insert login rate limit
		return u, errors.New("Incorrect password.")
	}

	return u, err
}

// ConfirmUser consumes the provided token and confirms the matching user, or returns an error.
func (s UserStore) ConfirmUser(token string) error {
	// get user by ConfirmToken
	var ct ConfirmToken
	var u User

	err := s.db.Model(&ct).Where("token = ?", token).Select()
	if err != nil {
		return errors.New("Invalid confirmation token.")
	}

	u, err = s.Get(ct.UserID)
	if err != nil {
		return errors.New("User with matching confirmation token not found.")
	}

	// set confirmed
	u.Confirmed = true
	_, err = s.db.Model(&u).Column("confirmed").Update()
	if err != nil {
		return errors.New("Error confirming user.")
	}

	// delete token
	_, err = s.db.Model(&ct).Where("token = ?", token).Delete()

	return err
}

// NewRecover generates a new recovery token and inserts it into the database, or returns an error.
func (s UserStore) NewRecover(uid int) (RecoverToken, error) {
	var rt RecoverToken

	token, err := GenerateRandomToken()
	if err != nil {
		return rt, err
	}
	rt.Token = token
	rt.UserID = uid
	rt.Expires = time.Now().AddDate(0, 0, 1)
	err = s.db.Create(&rt)

	return rt, err
}

// RecoverUser consumes the provided recovery token.
func (s UserStore) RecoverUser(token string) (User, error) {
	var u User
	var rt RecoverToken

	err := s.db.Model(&rt).Where("token = ?", token).Select()
	if err != nil {
		return u, errors.New("Invalid recover token.")
	}

	if rt.Expires.Before(time.Now()) {
		s.db.Model(&rt).Where("token = ?", token).Delete()

		return u, errors.New("The recover token has expired.")
	}

	u, err = s.Get(rt.UserID)
	if err != nil {
		return u, errors.New("User not found.")
	}

	_, err = s.db.Model(&rt).Where("token = ?", token).Delete()
	if err != nil {
		return u, errors.New("Error deleting used recover token.")
	}

	return u, err
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
