package main

import (
	"github.com/gin-gonic/gin"
)

func registerH(c *gin.Context) {
	// get POST data
	var u User
	c.Bind(&u)

	u, err := users.Create(u)
	if err != nil {
		c.String(500, err.Error())
		return
	}

	//c.String(200, "User "+ u.Name +" registered.")
	c.Redirect(303, "/")
}

func confirmH(c *gin.Context) {
	// take confirmation key
	var ct ConfirmToken
	c.Bind(&ct)

	err := users.ConfirmUser(ct.Token)
	if err != nil {
		c.String(500, err.Error())
		return
	}

	c.Redirect(303, "/")
}

func recoverH(c *gin.Context) {
	var rt RecoverToken
	c.Bind(&rt)

	if rt.Token == "" {
		c.String(400, "No recover token.")
		return
	}

	_, err := users.RecoverUser(rt.Token)
	if err != nil {
		c.String(400, err.Error())
		return
	}

	c.Redirect(303, "/")
}
