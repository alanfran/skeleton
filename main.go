package main

import (
	"github.com/gin-gonic/contrib/sessions"

	"./auth"
	"./blog"
	"./user"
	"gopkg.in/pg.v4"
)

var (
	db *pg.DB

	blogPosts blog.Storer
	users     user.Storer
	auths     auth.Storer
	mailer    Mailer

	cookieStore sessions.CookieStore

	sessionCookieName = getSessionCookieName()
	cookieSecret      = getCookieSecret()
	csrfSecret        = getCsrfSecret()

	dbAddr         = getDbAddr()
	dbUser         = getDbUser()
	dbPassword     = getDbPass()
	dbDatabase     = getDbDatabase()
	dbTestDatabase = getDbTestDatabase()

	appName = getAppName()
	appURL  = getAppURL()
)

func main() {
	// connect to DB
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

	// init data stores
	mailer = NewMockMailer()
	users = user.NewPgStore(db)
	auths = auth.NewPgStore(db)
	blogPosts = blog.NewPgStore(db)

	cookieStore = sessions.NewCookieStore([]byte(cookieSecret))

	// routes
	r := initRoutes()

	// run TLS
	r.Run(":8080")
}
