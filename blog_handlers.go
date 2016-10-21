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
		c.String(500, "Error retrieving blog posts.")
		return
	}

	data := buildData(c)
	data["Posts"] = posts
	data["older"] = false
	data["newer"] = false

	c.HTML(200, "blog", data)
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

	var b BlogPost
	err := c.Bind(&b)

	if err != nil {
		c.String(500, "Error saving post.")
		return
	}

	b.Author = u.ID
	blog.CreatePost(b)
}

func putBlogH(c *gin.Context) {
	var b BlogPost
	err := c.Bind(&b)
	if err != nil {
		c.String(500, err.Error())
		return
	}
	if b.Title == "" || b.Body == "" {
		c.String(500, "No data to insert.")
		return
	}
	err = blog.PutPost(b)
	if err != nil {
		c.String(500, err.Error())
	}

}

func deleteBlogH(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	err := blog.DelPost(id)
	if err != nil {
		c.String(500, "Error deleting post."+err.Error())
		fmt.Println(err.Error())
		return
	}
	c.String(200, "Deleted post #"+string(id))
}
