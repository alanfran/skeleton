package main

import (
	"fmt"
	"html/template"
	"strconv"

	"./blog"
	"./user"
	"github.com/gin-gonic/gin"
	"github.com/russross/blackfriday"
)

func blogHomeH(c *gin.Context) {
	posts, err := blogPosts.GetRecentPosts(3)
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
	b, err := blogPosts.GetPost(id)
	if err != nil {
		c.String(404, "Blog post not found.")
		return
	}
	c.JSON(200, b)
}

func postBlogH(c *gin.Context) {
	ctxUser, _ := c.Get("user")
	u := ctxUser.(user.User)

	var b blog.Post
	err := c.Bind(&b)

	if err != nil {
		c.String(500, "Error saving post.")
		return
	}

	b.Author = u.ID
	blogPosts.CreatePost(b)
}

func putBlogH(c *gin.Context) {
	var b blog.Post
	err := c.Bind(&b)
	if err != nil {
		c.String(500, err.Error())
		return
	}
	if b.Title == "" || b.Body == "" {
		c.String(500, "No data to insert.")
		return
	}
	err = blogPosts.PutPost(b)
	if err != nil {
		c.String(500, err.Error())
	}

}

func deleteBlogH(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	err := blogPosts.DelPost(id)
	if err != nil {
		c.String(500, "Error deleting post."+err.Error())
		fmt.Println(err.Error())
		return
	}
	c.String(200, "Deleted post #"+string(id))
}

func markdown(args ...interface{}) template.HTML {
	s := blackfriday.MarkdownCommon([]byte(fmt.Sprintf("%s", args...)))
	return template.HTML(s)
}
