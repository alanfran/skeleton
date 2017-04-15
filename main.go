package main

import (
	"github.com/gin-gonic/contrib/sessions"

	"./auth"
	"./blog"
	"./user"
)

func main() {
	// init app with data stores

	app := NewApp()
	app.cookieStore = sessions.NewCookieStore([]byte(app.cookieSecret))

	app.blog = blog.NewPgStore(app.db)
	app.users = user.NewPgStore(app.db)
	app.auths = auth.NewPgStore(app.db)

	// routes
	r := app.initRoutes()

	// run TLS
	r.Run(":8080")
}
