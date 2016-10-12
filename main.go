package main

import (
  //"net/http"
  //"html/template"
  //"encoding/base64"
  //"os"
  //"path/filepath"
  //"time"
  //"log"

  "github.com/gin-gonic/gin"
  "github.com/gin-gonic/contrib/sessions"
  "github.com/utrack/gin-csrf"

  "gopkg.in/pg.v4"
)

var (
  db *pg.DB

  blog *BlogStore
  users *UserStore
  auth *AuthStore
  mailer *Mailer

  cookieStore sessions.CookieStore

  CookieSecret = "Secret Used To Authenticate Cookies"
  CsrfSecret = "Insert Secret Here"
)

func main() {
  r := initRoutes()

  r.LoadHTMLGlob("views/*")

  // connect to DB
  db := pg.Connect(&pg.Options{
    User: "postgres",
    Password: "postgres",
    Database: "postgres",
  })
  // verify connection
  _, err := db.Exec(`SELECT 1`)
  if err != nil {
    panic("Error connecting to the database.")
  }

  // init data stores
  mailer = NewMailer()
  users = NewUserStore(db, mailer)
  blog = NewBlogStore(db)

  cookieStore = sessions.NewCookieStore([]byte(CookieSecret))

  // middleware
  r.Use(sessions.Sessions("session", cookieStore))
  r.Use(csrf.Middleware(csrf.Options{
        Secret: CsrfSecret,
        ErrorFunc: func(c *gin.Context){
            c.String(400, "CSRF token mismatch")
            c.Abort()
        },
    }))
  r.Use(secureOptions())

  // run TLS
  r.Run(":8080")
}
