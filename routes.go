package main

import (
	"fmt"
	"html/template"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (app *App) initRoutes() *gin.Engine {
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
	r.Use(sessions.Sessions(app.sessionCookieName, app.cookieStore))
	if gin.Mode() == gin.ReleaseMode {
		fmt.Println("Using secure middleware.")
		r.Use(secureOptions())
	}
	r.Use(app.csrfProtect())

	r.Use(app.setAuth)

	// HTML

	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index", buildData(c))
	})
	r.Static("/static", "static")

	r.GET("/blog", app.blogHomeH)

	// API
	authed := r.Group("/", authProtect)
	admin := authed.Group("/", adminProtect)

	// blog endpoints
	r.GET("/api/blog/:id", app.getBlogH)
	admin.POST("/api/blog", app.postBlogH)
	admin.PUT("/api/blog/:id", app.putBlogH)
	admin.DELETE("/api/blog/:id", app.deleteBlogH)

	// user endpoints
	r.POST("/api/register", app.registerH)
	r.GET("/api/confirm", app.confirmH)
	//r.GET("/api/recover", recoverH)

	// authentication endpoints
	r.POST("/api/login", app.loginH)
	authed.POST("/api/logout", app.logoutH)

	return r
}
