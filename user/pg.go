package user

import (
	"errors"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gopkg.in/pg.v4"
)

var (
	bcryptCost = 13
	errNyi     = errors.New("Not yet implemented.")
)

// PgStore stores a reference to the database and the Mailer and provides user CRUD methods.
type PgStore struct {
	db *pg.DB
}

// NewPgStore ensures the users, confirm_tokens, and recover_tokens tables exist.
// Then it returns an initialized PgStore.
func NewPgStore(database *pg.DB) *PgStore {
	// initialize database
	_, err := database.Exec(`CREATE TABLE IF NOT EXISTS users (
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
	return &PgStore{db: database}
}

// Create inserts the provided User into the database.
func (s PgStore) Create(u User) (User, error) {
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

	return u, err
}

// Put updates the provided user in the database.
func (s PgStore) Put(u User) error {
	err := s.db.Update(&u)

	return err
}

// Get returns the user with the provided ID, or an error.
func (s PgStore) Get(id int) (u User, err error) {
	u.ID = id
	err = s.db.Select(&u)

	return u, err
}

// GetByName returns the User with the provided name/email, or an error.
func (s PgStore) GetByName(n string) (u User, err error) {
	err = s.db.Model(&u).Where("name = ? OR email = ?", n, n).Select()

	return u, err
}

// Del deletes the user with the provided ID from the database, or returns an error.
func (s PgStore) Del(id int) error {
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
func (s PgStore) Validate(name, pass string) (User, error) {
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

// CreateConfirmationToken creates a confirmation token for a user and sets their account as unconfirmed.
func (s PgStore) CreateConfirmationToken(userid int) (ConfirmToken, error) {
	// generate confirmation token
	var confirm ConfirmToken
	t, err := GenerateRandomToken()
	if err != nil {
		return confirm, errors.New("Error generating confirmation token.")
	}
	confirm.UserID = userid
	confirm.Token = t
	err = s.db.Create(&confirm)

	return confirm, err
}

// ConfirmUser consumes the provided token and confirms the matching user, or returns an error.
func (s PgStore) ConfirmUser(token string) error {
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
func (s PgStore) NewRecover(uid int) (RecoverToken, error) {
	var rt RecoverToken

	token, err := GenerateRandomToken()
	if err != nil {
		return rt, err
	}
	rt.Token = token
	rt.UserID = uid
	rt.Expires = time.Now().AddDate(0, 0, 1)
	err = s.db.Create(&rt)
	if err != nil {
		return rt, errors.New("Error creating recovery token.")
	}

	return rt, err
}

// RecoverUser consumes the provided recovery token.
func (s PgStore) RecoverUser(token string) (User, error) {
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
