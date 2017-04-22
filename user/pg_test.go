package user

import (
	"testing"

	"gopkg.in/pg.v4"
)

var (
	db    *pg.DB
	users Storer

	dbAddr     = "localhost:5432"
	dbUser     = "postgres"
	dbPassword = "postgres"
	dbDatabase = "test"
)

func init() {
	db = pg.Connect(&pg.Options{
		Addr:     dbAddr,
		User:     dbUser,
		Password: dbPassword,
		Database: dbDatabase,
	})
	// verify connection
	_, err := db.Exec(`SELECT 1`)
	if err != nil {
		panic("Error connecting to the database.")
	}

	/*_, err = db.Exec(`DROP TABLE confirm_tokens;
		DROP TABLE recover_tokens;
		DROP TABLE users;`)
	if err != nil {
		panic(err)
	}*/

	users = NewPgStore(db)
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

	ct, err := users.CreateConfirmationToken(u.ID)
	if err != nil {
		t.Error("Error creating confirmation token.")
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

	// delete user
	err = users.Del(u2.ID)
	if err != nil {
		t.Error("Error deleting user.")
		t.Error(err)
	}
}
