package main

import (
	"fmt"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
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

	r.Use(setAuth)

	// HTML

	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index", buildData(c))
	})
	r.Static("/static", "static")

	r.GET("/blog", blogHomeH)
	//r.GET("/portfolio", portfolioH)
	//r.GET("/contact", contactH)

	// API
	a := r.Group("/", authProtect)

	// blog endpoints
	r.GET("/api/blog/:id", getBlogH)
	a.POST("/api/blog", postBlogH)
	a.PUT("/api/blog/:id", putBlogH)
	a.DELETE("/api/blog/:id", deleteBlogH)

	// user endpoints
	r.POST("/api/register", registerH)
	r.GET("/api/confirm", confirmH)
	//r.GET("/api/recover", recoverH)

	// authentication endpoints
	r.POST("/api/login", loginH)
	a.POST("/api/logout", logoutH)

	return r
}
