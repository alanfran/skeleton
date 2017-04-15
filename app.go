package main

import (
	"os"

	pg "gopkg.in/pg.v4"

	"./auth"
	"./blog"
	"./user"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

// Config stores our application's settings.
type Config struct {
	sessionCookieName string
	cookieSecret      string
	csrfSecret        string

	dbAddr     string
	dbUser     string
	dbPassword string
	dbDatabase string

	appName string
	appURL  string
}

// App stores the components and configuration of the application.
type App struct {
	Config

	db *pg.DB

	blog   blog.Storer
	users  user.Storer
	auths  auth.Storer
	mailer Mailer

	cookieStore sessions.CookieStore
}

// Option facilitates the functional option paradigm. Pass in any number of functions that modify the fields of the *App parameter.
type Option func(*App)

// NewApp initializes the application object. Configuration loads defaults which can be overriden by environment variables which can be overriden by functional options.
func NewApp(options ...Option) *App {
	app := &App{
	//db:          db,
	//	blog:        b,
	//users:       u,
	//	auths:       a,
	//cookieStore: c,
	}

	defaults := map[string]string{
		"SESSION_COOKIE_NAME": "skeleton",
		"COOKIE_SECRET":       "testing cookie secret",
		"CSRF_SECRET":         "testing csrf secret",
		"DB_ADDR":             "localhost:5432",
		"DB_USER":             "postgres",
		"DB_PASSWORD":         "postgres",
		"DB_DATABASE":         "postgres",
		"APP_NAME":            "skeleton",
		"APP_URL":             "localhost:8080",
	}

	// if testing, set test db
	if gin.Mode() == gin.TestMode {
		defaults["DB_DATABASE"] = "test"
	}

	// load any environment variables into defaults
	for k := range defaults {
		env := os.Getenv(k)
		if env != "" {
			defaults[k] = env
		}
	}

	app.sessionCookieName = defaults["SESSION_COOKIE_NAME"]
	app.cookieSecret = defaults["COOKIE_SECRET"]
	app.csrfSecret = defaults["CSRF_SECRET"]
	app.dbAddr = defaults["DB_ADDR"]
	app.dbUser = defaults["DB_USER"]
	app.dbPassword = defaults["DB_PASSWORD"]
	app.dbDatabase = defaults["DB_DATABASE"]
	app.appName = defaults["APP_NAME"]
	app.appURL = defaults["APP_URL"]

	// apply functional options
	for _, v := range options {
		v(app)
	}

	app.db = app.initDatabase()

	return app
}

func (a *App) initDatabase() *pg.DB {
	// connect to DB
	db := pg.Connect(&pg.Options{
		Addr:     a.dbAddr,
		User:     a.dbUser,
		Password: a.dbPassword,
		Database: a.dbDatabase,
	})
	// verify connection
	_, err := db.Exec(`SELECT 1`)
	if err != nil {
		panic("Error connecting to the database.")
	}

	return db
}
