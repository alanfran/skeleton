package main

import (
	"testing"
	//"fmt"
	//"strconv"
	"gopkg.in/pg.v4"
)

func init() {
	if db == nil {
		db = pg.Connect(&pg.Options{
			User:     dbUser,
			Password: dbPassword,
			Database: dbTestDatabase,
		})
		// verify connection
		_, err := db.Exec(`SELECT 1`)
		if err != nil {
			panic("Error connecting to the database.")
		}
	}

	users = NewUserStore(db, NewMailer())
}

func TestCreateAndConfirmUser(t *testing.T) {
	// Create User
	u, err := users.Create(User{
		Name:     "Test User",
		Password: "Test Pass",
		Email:    "test@email.com",
	})
	if err != nil {
		t.Error(err)
	}

	// test confirmation
	var ct ConfirmToken
	err = db.Model(&ct).Where("user_id = ?", u.ID).Select()
	if err != nil {
		t.Error("Error retrieving confirmation token.")
	}

	err = users.ConfirmUser(ct.Token)
	if err != nil {
		t.Error(err)
	}

	// ensure user is now confirmed
	u2, _ := users.Get(u.ID)
	if u2.Confirmed != true {
		t.Error("User not confirmed.")
	}

	users.Del(u.ID)
}

func TestValidateUser(t *testing.T) {
	p := "Test Pass"
	u, err := users.Create(User{
		Name:     "Test User",
		Password: p,
		Email:    "test@email.com",
	})
	if err != nil {
		t.Error(err)
	}
	defer users.Del(u.ID)
	defer db.Model(&ConfirmToken{}).Where("user_id = ?", u.ID).Delete()

	_, err = users.Validate(u.Name, p)
	if err != nil {
		t.Error(err)
	}
}

// test recover
func TestRecoverUser(t *testing.T) {
	u, err := users.Create(User{
		Name:     "Recover Test User",
		Password: "1234567",
		Email:    "recover@email.com",
	})
	if err != nil {
		t.Error("Error creating user in recover test.")
		t.Error(err)
		return
	}

	rt, err := users.NewRecover(u.ID)
	if err != nil {
		t.Error("Error creating recovery token.")
		t.Error(err)
		return
	}

	u2, err := users.RecoverUser(rt.Token)
	if err != nil {
		t.Error("Error looking up recovery token.")
		t.Error(err)
		return
	}

	// make sure the token has been consumed
	var ct2 RecoverToken
	err = users.db.Model(&ct2).Where("token = ?", rt.Token).Select()
	if err == nil {
		t.Error("The recovery token has not been deleted after use.")
		t.Error(err)
	}

	var confirm ConfirmToken
	_, _ = users.db.Model(&confirm).Where("user_id = ?", u2.ID).Delete()

	// delete user
	err = users.Del(u2.ID)
	if err != nil {
		t.Error("Error deleting user.")
		t.Error(err)
	}
}
