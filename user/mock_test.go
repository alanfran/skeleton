package user

import (
	"testing"
)

var (
	mock Storer
)

func init() {
	mock = NewMockStore()
}

func TestMockCreateAndConfirmUser(t *testing.T) {
	// Create User
	u, err := mock.Create(User{
		Name:     "Test User",
		Password: "Test Pass",
		Email:    "test@email.com",
	})
	if err != nil {
		t.Error(err)
	}

	ct, err := mock.CreateConfirmationToken(u.ID)
	if err != nil {
		t.Error(err)
	}

	err = mock.ConfirmUser(ct.Token)
	if err != nil {
		t.Error(err)
	}

	// ensure user is now confirmed
	u2, _ := mock.Get(u.ID)
	if u2.Confirmed != true {
		t.Error("User not confirmed.")
	}

	mock.Del(u.ID)
}

func TestMockValidateUser(t *testing.T) {
	p := "Test Pass"
	u, err := mock.Create(User{
		Name:     "Test User",
		Password: p,
		Email:    "test@email.com",
	})
	if err != nil {
		t.Error(err)
	}
	defer mock.Del(u.ID)

	_, err = mock.Validate(u.Name, p)
	if err != nil {
		t.Error(err)
	}
}

// test recover
func TestMockRecoverUser(t *testing.T) {
	u, err := mock.Create(User{
		Name:     "Recover Test User",
		Password: "1234567",
		Email:    "recover@email.com",
	})
	if err != nil {
		t.Error("Error creating user in recover test.")
		t.Error(err)
		return
	}

	rt, err := mock.NewRecover(u.ID)
	if err != nil {
		t.Error("Error creating recovery token.")
		t.Error(err)
		return
	}

	u2, err := mock.RecoverUser(rt.Token)
	if err != nil {
		t.Error("Error looking up recovery token.")
		t.Error(err)
		return
	}

	// delete user
	err = mock.Del(u2.ID)
	if err != nil {
		t.Error("Error deleting user.")
		t.Error(err)
	}
}
