package auth

import (
	"testing"
)

var (
	mock Storer
)

func init() {
	mock = NewMockStore()
}

func TestMock(t *testing.T) {
	s, err := mock.Create(1337, "127.0.0.1")
	if err != nil {
		t.Error(err)
	}
	s2, err := mock.Get(s.Key)
	if err != nil {
		t.Error(err)
	}
	if s2.UserID != 1337 || s2.IP != "127.0.0.1" {
		t.Error("Retrieved auth does not match input.")
	}
	err = mock.Del(s2.Key)
	if err != nil {
		t.Error(err)
	}
}
