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

	r.GET("/board", forumHomeH)
	r.GET("/board/b/:id", forumBoardH)
	r.GET("/board/t/:tid", forumThreadH)
	r.GET("/board/recent", forumRecentH)
	r.GET("/board/popular", forumPopularH)

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

	// forum endpoints
	r.GET("/api/boards", getBoardsH)
	r.GET("/api/board/:id", getThreadsH)
	r.GET("/api/thread/:id", getPostsH)

	admin.POST("/api/board", postBoardH)
	authed.POST("/api/board/:id", postThreadH)
	authed.POST("/api/thread/:id", postPostH)

	admin.PUT("/api/board/:id", putBoardH)
	authed.PUT("/api/thread/:tid", putThreadH)
	authed.PUT("/api/post/:id", putPostH)

	admin.DELETE("/api/board/:id", deleteBoardH)
	authed.DELETE("/api/thread/:id", deleteThreadH)
	authed.DELETE("/api/post/:id", deletePostH)

	return r
}
