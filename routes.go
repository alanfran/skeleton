package main

import (
	"fmt"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

func initRoutes() *gin.Engine {
	// logging + Recovery
	r := gin.Default()

	r.RedirectTrailingSlash = true

	// global middleware
	r.Use(sessions.Sessions(sessionCookieName, cookieStore))
	if gin.Mode() == gin.ReleaseMode {
		fmt.Println("Using secure middleware.")
		r.Use(secureOptions())
	}
	r.Use(csrfProtect())

	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index", gin.H{
			"_csrf": csrf.GetToken(c),
			"Auth":  getAuthed(c),
		})
	})
	r.Static("/static", "static")

	r.GET("/blog", blogHomeH)

	// blog endpoints
	r.GET("/api/blog/:id", getBlogH)
	r.POST("/api/blog", postBlogH)
	r.PUT("/api/blog/:id", putBlogH)
	r.DELETE("/api/blog/:id", deleteBlogH)

	// user endpoints
	r.POST("/api/register", registerH)
	r.GET("/api/confirm", confirmH)

	// session endpoints
	r.POST("/api/login", loginH)

	// must be logged in
	a := r.Group("/", authProtect)

	a.POST("/api/logout", logoutH)

	return r
}
