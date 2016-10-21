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
	forum  *ForumStore
	mailer *Mailer

	cookieStore sessions.CookieStore

	sessionCookieName = "session"
	cookieSecret      = "Secret Used To Authenticate Cookies"
	csrfSecret        = "Insert Secret Here"

	dbUser         = "postgres"
	dbPassword     = "postgres"
	dbDatabase     = "postgres"
	dbTestDatabase = "test"

	appName = "pg-skeleton"
	appURL  = "localhost:8080"
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
	auth = NewAuthStore(db)
	blog = NewBlogStore(db)

	cookieStore = sessions.NewCookieStore([]byte(cookieSecret))

	// routes
	r := initRoutes()

	r.LoadHTMLGlob("views/*")

	// run TLS
	r.Run(":8080")
}
