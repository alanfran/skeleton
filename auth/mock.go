package auth

import (
	"errors"
	"time"
)

// MockStore provides an in-memory store for testing purposes.
type MockStore struct {
	mem map[string]Auth
}

// Create places an Auth record in the MockStore.
func (s MockStore) Create(userid int, ip string) (Auth, error) {
	var a Auth
	k, err := GenerateAuthKey()
	if err != nil {
		return a, err
	}

	a.UserID = userid
	a.Key = k
	a.IP = ip
	a.Created = time.Now()
	a.Expires = time.Now().AddDate(1, 0, 0)

	s.mem[a.Key] = a
	return a, nil
}

// Get returns the Auth with a matching key.
func (s MockStore) Get(key string) (Auth, error) {
	a, ok := s.mem[key]
	if !ok {
		return a, errors.New("Cannot retrieve Auth: key not found.")
	}

	return a, nil
}

// Del removes the Auth with a matching key, or returns an error if one is not found.
func (s MockStore) Del(key string) error {
	_, ok := s.mem[key]
	if !ok {
		return errors.New("Cannot delete Auth: key not found.")
	}
	delete(s.mem, key)
	return nil
}
