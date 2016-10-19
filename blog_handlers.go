package main

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

func blogHomeH(c *gin.Context) {
	posts, err := blog.GetRecentPosts(3)
	if err != nil {
		fmt.Println(err)
		c.String(500, "Error retrieving blog posts.")
		return
	}

	c.HTML(200, "blog", gin.H{
		"_csrf": csrf.GetToken(c),
		"Posts": posts,
	})
}

func getBlogH(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	b, err := blog.GetPost(id)
	if err != nil {
		c.String(404, "Blog post not found.")
		return
	}
	c.JSON(200, b)
}

func postBlogH(c *gin.Context) {
	user, _ := c.Get("user")
	u := user.(User)

	if !u.Admin {
		c.String(500, "You are not authorized to make blog posts.")
	}

	var b BlogPost
	err := c.Bind(&b)
	if err != nil {
		c.String(500, "Error saving post.")
		return
	}
	blog.CreatePost(b)
}

func putBlogH(c *gin.Context) {
	var b BlogPost
	err := c.Bind(&b)
	if err != nil {
		c.String(500, "Error updating post.")
		return
	}
	blog.PutPost(b)
}

func deleteBlogH(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	err := blog.DelPost(id)
	if err != nil {
		c.AbortWithStatus(500)
		return
	}
	c.String(200, "Deleted post #"+string(id))
}
