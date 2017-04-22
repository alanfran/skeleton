package user

import (
	"errors"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// MockStore implements an in-memory User storer for testing purposes.
type MockStore struct {
	id    map[int]User
	name  map[string]User
	email map[string]User

	confirm map[string]ConfirmToken
	recover map[string]RecoverToken

	idSerial int
}

// NewMockStore returns an initialized MockStore.
func NewMockStore() *MockStore {
	return &MockStore{
		id:      map[int]User{},
		name:    map[string]User{},
		email:   map[string]User{},
		confirm: map[string]ConfirmToken{},
		recover: map[string]RecoverToken{},
	}
}

// Create inserts the provided User into the memory store.
func (s *MockStore) Create(u User) (User, error) {
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

	u.ID = s.idSerial
	s.idSerial++

	// insert into db
	s.Put(u)

	return u, err
}

// Put updates the provided user in the database.
func (s MockStore) Put(u User) error {
	s.id[u.ID] = u
	s.name[u.Name] = u
	s.email[u.Email] = u

	return nil
}

// Get returns the user with the provided ID, or an error.
func (s MockStore) Get(id int) (u User, err error) {
	u, ok := s.id[id]
	if !ok {
		err = errors.New("User not found in memory store.")
	}

	return u, err
}

// GetByName returns the User with the provided name/email, or an error.
func (s MockStore) GetByName(n string) (u User, err error) {
	u, ok := s.name[n]
	if !ok {
		u, ok = s.email[n]
		if !ok {
			err = errors.New("User not found.")
		}
	}
	return u, err
}

// Del deletes the user with the provided ID from the memory store, or returns an error.
func (s MockStore) Del(id int) error {
	u, err := s.Get(id)
	delete(s.id, id)
	delete(s.name, u.Name)
	delete(s.email, u.Email)

	// delete tokens associated with uid
	return err
}

// Validate checks if the provided name/email and password are correct.
func (s MockStore) Validate(name, pass string) (User, error) {
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
func (s MockStore) CreateConfirmationToken(userid int) (ConfirmToken, error) {
	// generate confirmation token
	var confirm ConfirmToken
	t, err := GenerateRandomToken()
	if err != nil {
		return confirm, errors.New("Error generating confirmation token.")
	}
	confirm.UserID = userid
	confirm.Token = t

	s.confirm[confirm.Token] = confirm

	return confirm, err
}

// ConfirmUser consumes the provided token and confirms the matching user, or returns an error.
func (s MockStore) ConfirmUser(token string) error {
	// get user by ConfirmToken
	var ct ConfirmToken
	var u User

	ct, ok := s.confirm[token]
	if !ok {
		return errors.New("Invalid confirmation token.")
	}

	u, err := s.Get(ct.UserID)
	if err != nil {
		return errors.New("User with matching confirmation token not found.")
	}

	// set confirmed
	u.Confirmed = true
	err = s.Put(u)
	if err != nil {
		return errors.New("Error confirming user.")
	}

	// delete token
	delete(s.confirm, token)

	return err
}

// NewRecover generates a new recovery token and inserts it into the memory store, or returns an error.
func (s MockStore) NewRecover(uid int) (RecoverToken, error) {
	var rt RecoverToken

	token, err := GenerateRandomToken()
	if err != nil {
		return rt, err
	}
	rt.Token = token
	rt.UserID = uid
	rt.Expires = time.Now().AddDate(0, 0, 1)

	s.recover[rt.Token] = rt

	return rt, err
}

// RecoverUser consumes the provided recovery token.
func (s MockStore) RecoverUser(token string) (User, error) {
	var u User
	var rt RecoverToken

	rt, ok := s.recover[token]
	if !ok {
		return u, errors.New("Invalid recover token.")
	}

	if rt.Expires.Before(time.Now()) {
		delete(s.recover, token)
		return u, errors.New("The recover token has expired.")
	}

	u, err := s.Get(rt.UserID)
	if err != nil {
		return u, errors.New("User not found.")
	}

	delete(s.recover, token)

	return u, err
}
