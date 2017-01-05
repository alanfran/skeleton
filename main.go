package main

import (
	"github.com/gin-gonic/contrib/sessions"

	"gopkg.in/pg.v4"
)

var (
	db *pg.DB

	blog   *BlogStore
	users  *UserStore
	auth   *AuthStore
	mailer *Mailer

	cookieStore sessions.CookieStore

	sessionCookieName = getSessionCookieName()
	cookieSecret      = getCookieSecret()
	csrfSecret        = getCsrfSecret()

	dbUser         = getDbUser()
	dbPassword     = getDbPass()
	dbDatabase     = getDbDatabase()
	dbTestDatabase = getDbTestDatabase()

	appName = getAppName()
	appURL  = getAppURL()
)

func main() {
	// check environment variables

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
	auth = NewAuthStore(db)
	blog = NewBlogStore(db)

	cookieStore = sessions.NewCookieStore([]byte(cookieSecret))

	// routes
	r := initRoutes()

	// run TLS
	r.Run(":8080")
}
