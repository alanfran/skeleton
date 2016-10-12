package main

import (
  "time"
  "errors"
  "strings"
  "golang.org/x/crypto/bcrypt"
  "gopkg.in/pg.v4"
  "crypto/rand"
  "encoding/base64"

  //"fmt"
  //"strconv"
)


var (
  bcryptCost = 13
  nyi = errors.New("Not yet implemented.")

)

type User struct {
  ID             int
	Name           string
	Email          string
	Password       string

  Created        time.Time
  // email confirmed
	Confirmed      bool

  // login lock
	AttemptNumber  int64
	AttemptTime    time.Time `sql:",null"`
	Locked         time.Time `sql:",null"`
}

type ConfirmToken struct {
  UserID          int       ``
  Token           string
}

type RecoverToken struct {
  UserID          int
  Token           string
  Expires         time.Time
}

type UserStore struct {
  db   *pg.DB
  mailer *Mailer
}

func NewUserStore(database *pg.DB, mailer *Mailer) *UserStore {
  // initialize database
  _, err := database.Exec(`DROP TABLE users; CREATE TABLE IF NOT EXISTS users (
    id serial PRIMARY KEY,
    name text NOT NULL,
    email text NOT NULL,
    password text NOT NULL,
    created timestamp NOT NULL,
    confirmed boolean NOT NULL,
    attempt_number bigint,
    attempt_time timestamp,
    locked timestamp )`)
  if err != nil {
    panic("Error initializing user table.")
  }

  _, err = database.Exec(`CREATE TABLE IF NOT EXISTS confirm_tokens(
    user_id   int  UNIQUE NOT NULL,
    token text UNIQUE NOT NULL )`)
  if err != nil {
    panic("Error initializing confirmation token table.")
  }

  _, err = database.Exec(`CREATE TABLE IF NOT EXISTS recover_tokens(
    user_id     int  UNIQUE NOT NULL,
    token   text UNIQUE NOT NULL,
    expires timestamp   NOT NULL )`)
  if err != nil {
    panic("Error initializing recovery token table.")
  }

  // initialize the memory store
  return &UserStore{ database, mailer }
}

func (s UserStore) Create(u User) (User, error) {
  // name unique
  _, err := s.GetByName(u.Name)
  if err == nil {     // successful query, err != ErrNoRows
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

func (s UserStore) Put(u User) error {
  err := s.db.Update(&u)

  return err
}

func (s UserStore) Get(id int) (u User, err error) {
  u.ID = id
  err = s.db.Select(&u)

  return u, err
}

func (s UserStore) GetByName(n string) (u User, err error) {
  err = s.db.Model(&u).Where("name = ? OR email = ?", n, n).Select()

  return u, err
}

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

func (s UserStore) Validate(name, pass string) error {
  u, err := s.GetByName(name)
  if err != nil {
    return err
  }
  err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(pass))
  if err != nil {
    // insert login rate limit

    return err
  }

  return err
}

func (s UserStore) ConfirmUser(token string) error {
	// get user by ConfirmToken
  var ct ConfirmToken
  var u User

  err := s.db.Model(&ct).Where("token = ?", token).Select()
  if err != nil {
    //fmt.Println("Token not found.")
    return err
  }

  u, err = s.Get(ct.UserID)
  if err != nil {
    //fmt.Println("User not found.")
    return err
  }

  // set confirmed
  u.Confirmed = true
  _, err = s.db.Model(&u).Column("confirmed").Update()
  if err != nil {
    return err
  }

  // delete token
  _, err = s.db.Model(&ct).Where("token = ?", token).Delete()

	return err
}

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

func (s UserStore) RecoverUser(token string) (User, error) {
	var u User
  var rt RecoverToken

  err := s.db.Model(&rt).Where("token = ?", token).Select()
  if err != nil {
    return u, err
  }

  if rt.Expires.Before(time.Now()) {
    _, err = s.db.Model(&rt).Where("token = ?", token).Delete()
    if err != nil {
      return u, err
    }
    return u, errors.New("The recover token has expired.")
  }

  u, err = s.Get(rt.UserID)
  if err != nil {
    return u, err
  }

  _, err = s.db.Model(&rt).Where("token = ?", token).Delete()
  if err != nil {
    return u, err
  }

	return u, err
}

func GenerateRandomToken() (string, error) {
  b := make([]byte, 32)
  _, err := rand.Read(b)
  if err != nil {
    return "", err
  }
  return base64.URLEncoding.EncodeToString(b), err
}
