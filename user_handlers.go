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
    c.AbortWithStatus(500)
  }

  c.String(200, "User "+ u.Name +" registered.")
}

func confirmH(c *gin.Context) {
  // take confirmation key
  var ct ConfirmToken
  c.Bind(&ct)

  err := users.ConfirmUser(ct.Token)
  if err != nil {
    c.AbortWithStatus(500)
  }

  c.String(200, "Your email address has been confirmed.")
}

func recoverH(c *gin.Context) {

}
