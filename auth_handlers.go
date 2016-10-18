package main

import (
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

var (
	// AuthKey is the name of the Auth cookie.
	AuthKey = "auth"
)

func loginH(c *gin.Context) {
	// if already logged in, send back error msg

	var u User
	c.Bind(&u)

	u, err := users.Validate(u.Name, u.Password)
	if err != nil {
		c.String(500, "Invalid login credentials.")
		return
	}

	// create session
	a, err := auth.Create(u.ID, c.ClientIP())
	if err != nil {
		c.String(500, "Error creating session.")
		return
	}

	// set cookie
	session := sessions.Default(c)
	session.Set(AuthKey, a.Key)
	session.Save()

	// redirect
	c.Redirect(303, "/")
}

// authProtect middleware ensures user is logged in
func logoutH(c *gin.Context) {
	// get auth key
	key, exists := c.Get(AuthKey)
	if !exists {
		c.String(500, "You are not logged in.")
		return
	}

	// destroy auth
	err := auth.Del(key.(string))
	if err != nil {
		c.String(500, err.Error())
	}

	session := sessions.Default(c)
	session.Delete(AuthKey)

	// redirect
	c.Redirect(303, "/")
}
