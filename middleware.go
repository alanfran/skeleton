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

func authProtect(c *gin.Context) {
	// read cookie
	s := sessions.Default(c)
	key := s.Get(AuthKey)
	// look up session in db
	_, err := auth.Get(key.(string))
	if err != nil {
		c.AbortWithError(403, errors.New("You are not logged in."))
	}
	// set key in context
	c.Set(AuthKey, key)
	c.Next()
}

func csrfProtect() gin.HandlerFunc {
	return csrf.Middleware(csrf.Options{
		Secret: CsrfSecret,
		ErrorFunc: func(c *gin.Context) {
			c.String(400, "CSRF token mismatch.")
			c.Abort()
		},
	})
}
