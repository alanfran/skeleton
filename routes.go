package main

import (
  "github.com/gin-gonic/gin"
)

func initRoutes() *gin.Engine {
  // logging + Recovery
  r := gin.Default()

  r.GET("/", func(c *gin.Context) {
    c.String(200, "Ok.")
  })

  // blog home page
  r.GET("/blog", blogHomeH)

  // blog endpoints
  r.GET("/blog/:id", getBlogH)
  r.POST("/blog", postBlogH)
  r.PUT("/blog/:id", putBlogH)
  r.DELETE("/blog/:id", deleteBlogH)

  // user endpoints

  // session endpoints

  return r
}
