package main

import (
	//"fmt"
	//"log"
	//"net/http"

	//"github.com/justinas/nosurf"

	"errors"

	"github.com/gin-gonic/contrib/secure"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

func secureOptions() gin.HandlerFunc {
	return secure.Secure(secure.Options{
		//AllowedHosts:          []string{"example.com", "ssl.example.com"},
		SSLRedirect: true,
		//SSLHost:               "ssl.example.com",
		SSLProxyHeaders:       map[string]string{"X-Forwarded-Proto": "https"},
		STSSeconds:            315360000,
		STSIncludeSubdomains:  true,
		FrameDeny:             true,
		ContentTypeNosniff:    true,
		BrowserXssFilter:      true,
		ContentSecurityPolicy: "default-src 'self'",
	})
}

func setAuth(c *gin.Context) {
	// read cookie
	s := sessions.Default(c)
	key := s.Get(AuthKey).(string)
	// look up session in db
	a, err := auth.Get(key)
	if err != nil {
		c.Next()
		return
	}
	// set key in context
	c.Set(AuthKey, key)

	u, err := users.Get(a.UserID)
	if err != nil {
		c.Next()
		return
	}

	c.Set("user", u)

	c.Next()
}

func authProtect(c *gin.Context) {
	_, exists := c.Get(AuthKey)
	if !exists {
		c.AbortWithError(403, errors.New("You are not logged in."))
	}
	c.Next()
}

func csrfProtect() gin.HandlerFunc {
	return csrf.Middleware(csrf.Options{
		Secret: csrfSecret,
		ErrorFunc: func(c *gin.Context) {
			c.String(400, "CSRF token mismatch.")
			c.Abort()
		},
	})
}

// Helper Functions for building template data
func isAdmin(c *gin.Context) bool {
	u, exists := c.Get("user")
	if !exists {
		return false
	}

	usr := u.(User)
	if !usr.Admin {
		return false
	}

	return true
}

func isAuthed(c *gin.Context) bool {
	_, exists := c.Get(AuthKey)
	if !exists {
		return false
	}
	return true
}

func buildData(c *gin.Context) gin.H {
	return gin.H{
		"Auth":  isAuthed(c),
		"Admin": isAdmin(c),
		"_csrf": csrf.GetToken(c),
	}
}
