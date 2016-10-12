package main

import (
  "github.com/gin-gonic/gin"
)

func loginH(c *gin.Context) {
  // if already logged in, send back error msg

  var u User
  c.Bind(&u)

  err := users.Validate(u.Name, u.Password)
  if err != nil {
    c.String(500, "Invalid login credentials.")
  }
}

func logoutH(c *gin.Context) {
  // get auth key

  // if not logged in, send back error msg

  // destroy auth

  // redirect
}
