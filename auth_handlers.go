package main

import (
	"github.com/alanfran/skeleton/user"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

var (
	// AuthKey is the name of the Auth cookie.
	AuthKey = "auth"
)

func (app *App) loginH(c *gin.Context) {
	var u user.User
	c.Bind(&u)

	u, err := app.users.Validate(u.Name, u.Password)
	if err != nil {
		c.String(500, "Invalid login credentials.")
		return
	}

	// create auth token in database
	a, err := app.auths.Create(u.ID, c.ClientIP())
	if err != nil {
		c.String(500, "Error creating session.")
		return
	}

	// set cookie
	session := sessions.Default(c)
	session.Set(AuthKey, a.Key)
	session.Save()

	c.String(200, "Success.")
}

// authProtect middleware ensures user is logged in
func (app *App) logoutH(c *gin.Context) {
	// get auth key
	key, exists := c.Get(AuthKey)
	if !exists {
		c.String(500, "You are not logged in.")
		return
	}

	// delete auth from database
	err := app.auths.Del(key.(string))
	if err != nil {
		c.String(500, err.Error())
	}

	session := sessions.Default(c)
	session.Delete(AuthKey)

	// redirect
	c.Redirect(303, "/")
}
