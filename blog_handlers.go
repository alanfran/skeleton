package main

import (
  "fmt"
  "strconv"
  "github.com/gin-gonic/gin"
)

func blogHomeH(c *gin.Context) {
  posts, err := blog.GetRecentPosts(3)
  if err != nil {
    fmt.Println(err)
    c.AbortWithStatus(500)
  }

  c.JSON(200, posts)
}

func getBlogH(c *gin.Context) {
  id, _ := strconv.Atoi(c.Param("id"))
  b, err := blog.GetPost(id)
  if err != nil {
    c.AbortWithStatus(500)
  }
  c.JSON(200, b)
}

func postBlogH(c *gin.Context) {
  var b BlogPost
  err := c.Bind(&b)
  if err != nil {
    c.AbortWithStatus(500)
  }
  blog.CreatePost(b)
}

func putBlogH(c *gin.Context) {
  var b BlogPost
  err := c.Bind(&b)
  if err != nil {
    c.AbortWithStatus(500)
  }
  blog.PutPost(b)
}

func deleteBlogH(c *gin.Context) {
  id, _ := strconv.Atoi(c.Param("id"))
  err := blog.DelPost(id)
  if err != nil {
    c.AbortWithStatus(500)
  }
  c.String(200, "Deleted post #"+string(id))
}
