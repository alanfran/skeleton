package main

import (
  "github.com/gin-gonic/gin"
)

func initRoutes() *gin.Engine {
  // logging + Recovery
  r := gin.Default()

  r.RedirectTrailingSlash = true

  r.GET("/", func(c *gin.Context) {
    c.String(200, "Ok.")
  })

  // blog home page
  r.GET("/api/blog", blogHomeH)

  // blog endpoints
  r.GET("/api/blog/:id", getBlogH)
  r.POST("/api/blog", postBlogH)
  r.PUT("/api/blog/:id", putBlogH)
  r.DELETE("/api/blog/:id", deleteBlogH)

  // user endpoints
  r.POST("/api/register", registerH)

  // session endpoints
  // csrf
  r.POST("/api/login", loginH)

  a := r.Group("/", authProtect)
  // csrf
  a.POST("/api/logout", logoutH)

  return r
}
