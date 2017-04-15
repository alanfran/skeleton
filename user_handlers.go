package main

import (
	"github.com/gin-gonic/gin"

	"./user"
)

func (app *App) registerH(c *gin.Context) {
	// get POST data
	var u user.User
	c.Bind(&u)

	u, err := app.users.Create(u)
	if err != nil {
		c.String(500, err.Error())
		return
	}

	c.String(200, "User "+u.Name+" registered.")
}

func (app *App) confirmH(c *gin.Context) {
	// take confirmation key
	var ct user.ConfirmToken
	c.Bind(&ct)

	err := app.users.ConfirmUser(ct.Token)
	if err != nil {
		c.String(500, err.Error())
		return
	}

	c.Redirect(303, "/")
}

func (app *App) recoverH(c *gin.Context) {
	var rt user.RecoverToken
	c.Bind(&rt)

	if rt.Token == "" {
		c.String(400, "No recover token.")
		return
	}

	_, err := app.users.RecoverUser(rt.Token)
	if err != nil {
		c.String(400, err.Error())
		return
	}

	c.Redirect(303, "/")
}
