package main

import (
	"os"

	"github.com/gin-gonic/gin"
)

func getSessionCookieName() string {
	name := os.Getenv("SESSION_COOKIE_NAME")
	if name == "" {
		if gin.Mode() == gin.ReleaseMode {
			panic("Please set the SESSION_COOKIE_NAME environment variable.")
		}
		name = "skeleton"
	}
	return name
}

func getCookieSecret() string {
	secret := os.Getenv("COOKIE_SECRET")
	if secret == "" {
		if gin.Mode() == gin.ReleaseMode {
			panic("Please set the COOKIE_SECRET environment variable.")
		}
		secret = "change this secret used to authenticate cookies"
	}
	return secret
}

func getCsrfSecret() string {
	secret := os.Getenv("CSRF_SECRET")
	if secret == "" {
		if gin.Mode() == gin.ReleaseMode {
			panic("Please set the CSRF_SECRET environment variable.")
		}
		secret = "change this csrf secret"
	}
	return secret
}

func getDbAddr() string {
	addr := os.Getenv("DB_ADDR")
	if addr == "" {
		addr = "localhost:5432"
	}
	return addr
}

func getDbUser() string {
	user := os.Getenv("DB_USER")
	if user == "" {
		if gin.Mode() == gin.ReleaseMode {
			panic("Please set the DB_USER environment variable.")
		}
		user = "postgres"
	}
	return user
}

func getDbPass() string {
	pass := os.Getenv("DB_PASSWORD")
	if pass == "" {
		if gin.Mode() == gin.ReleaseMode {
			panic("Please set the DB_PASSWORD environment variable.")
		}
		pass = "postgres"
	}
	return pass
}

func getDbDatabase() string {
	dbname := os.Getenv("DB_DATABASE")
	if dbname == "" {
		if gin.Mode() == gin.ReleaseMode {
			panic("Please set the DB_DATABASE environment variable.")
		}
		dbname = "postgres"
	}
	return dbname
}

func getDbTestDatabase() string {
	tdb := os.Getenv("DB_TEST_DATABASE")
	if tdb == "" {
		tdb = "test"
	}
	return tdb
}

func getAppName() string {
	app := os.Getenv("APP_NAME")
	if app == "" {
		if gin.Mode() == gin.ReleaseMode {
			panic("Please set the APP_NAME environment variable.")
		}
		app = "skeleton"
	}
	return app
}

func getAppURL() string {
	url := os.Getenv("APP_URL")
	if url == "" {
		if gin.Mode() == gin.ReleaseMode {
			panic("Please set the APP_URL environment variable.")
		}
		url = "localhost:8080"
	}
	return url
}
