package main

import (
	//"net/http"
	//"html/template"
	//"encoding/base64"
	//"os"
	//"path/filepath"
	//"time"
	//"log"

	"fmt"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"

	"gopkg.in/pg.v4"
)

var (
	db *pg.DB

	blog   *BlogStore
	users  *UserStore
	auth   *AuthStore
	mailer *Mailer

	cookieStore sessions.CookieStore

	SessionCookieName = "session"
	CookieSecret      = "Secret Used To Authenticate Cookies"
	CsrfSecret        = "Insert Secret Here"

	dbUser     = "postgres"
	dbPassword = "postgres"
	dbDatabase = "postgres"
)

func main() {
	// connect to DB
	db = pg.Connect(&pg.Options{
		User:     dbUser,
		Password: dbPassword,
		Database: dbDatabase,
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

	// routes
	r := initRoutes()

	r.LoadHTMLGlob("views/*")

	// global middleware
	r.Use(sessions.Sessions(SessionCookieName, cookieStore))
	if gin.Mode() == gin.ReleaseMode {
		fmt.Println("Using secure middleware.")
		r.Use(secureOptions())
	}

	// run TLS
	r.Run(":8080")
}
