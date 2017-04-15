package main

import (
	"errors"

	"github.com/alanfran/skeleton/user"

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

// Looks for Auth cookie and adds it and the User to the context.
func (app *App) setAuth(c *gin.Context) {
	// read cookie
	s := sessions.Default(c)
	k := s.Get(AuthKey)
	if k == nil {
		c.Next()
		return
	}
	key := k.(string)
	// look up session in db
	a, err := app.auths.Get(key)
	if err != nil {
		c.Next()
		return
	}
	// set key in context
	c.Set(AuthKey, key)

	// add User to context
	u, err := app.users.Get(a.UserID)
	if err != nil {
		c.Next()
		return
	}

	c.Set("user", u)

	c.Next()
}

func authProtect(c *gin.Context) {
	if !isAuthed(c) {
		c.AbortWithError(403, errors.New("You are not logged in."))
	}
	c.Next()
}

func adminProtect(c *gin.Context) {
	if !isAdmin(c) {
		c.AbortWithError(403, errors.New("Permission denied."))
	}
	c.Next()
}

func (app *App) csrfProtect() gin.HandlerFunc {
	return csrf.Middleware(csrf.Options{
		Secret: app.csrfSecret,
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

	usr := u.(user.User)
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
