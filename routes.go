package main

import (
	"fmt"
	"html/template"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

func initRoutes() *gin.Engine {
	// logging + Recovery
	r := gin.Default()

	r.RedirectTrailingSlash = true

	//r.LoadHTMLGlob("views/*")

	var funcMap = template.FuncMap{
		"markdown": markdown,
	}

	if tmpl, err := template.New("blog").Funcs(funcMap).ParseGlob("views/*"); err == nil {
		r.SetHTMLTemplate(tmpl)
	} else {
		panic(err)
	}

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

	// API
	authed := r.Group("/", authProtect)
	admin := authed.Group("/", adminProtect)

	// blog endpoints
	r.GET("/api/blog/:id", getBlogH)
	admin.POST("/api/blog", postBlogH)
	admin.PUT("/api/blog/:id", putBlogH)
	admin.DELETE("/api/blog/:id", deleteBlogH)

	// user endpoints
	r.POST("/api/register", registerH)
	r.GET("/api/confirm", confirmH)
	//r.GET("/api/recover", recoverH)

	// authentication endpoints
	r.POST("/api/login", loginH)
	authed.POST("/api/logout", logoutH)

	return r
}
