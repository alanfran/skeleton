package main

import (
	"fmt"
	"html/template"
	"strconv"

	"github.com/alanfran/skeleton/blog"
	"github.com/alanfran/skeleton/user"
	"github.com/gin-gonic/gin"
	"github.com/russross/blackfriday"
)

func (app *App) blogHomeH(c *gin.Context) {
	posts, err := app.blog.GetRecentPosts(3)
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

func (app *App) getBlogH(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	b, err := app.blog.GetPost(id)
	if err != nil {
		c.String(404, "Blog post not found.")
		return
	}
	c.JSON(200, b)
}

func (app *App) postBlogH(c *gin.Context) {
	ctxUser, _ := c.Get("user")
	u := ctxUser.(user.User)

	var b blog.Post
	err := c.Bind(&b)

	if err != nil {
		c.String(500, "Error saving post.")
		return
	}

	b.Author = u.ID
	app.blog.CreatePost(b)
}

func (app *App) putBlogH(c *gin.Context) {
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
	err = app.blog.PutPost(b)
	if err != nil {
		c.String(500, err.Error())
	}

}

func (app *App) deleteBlogH(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	err := app.blog.DelPost(id)
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
