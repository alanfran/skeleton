package main

import (
	//"fmt"
	//"log"
	//"net/http"

	//"github.com/justinas/nosurf"

  "github.com/gin-gonic/gin"
  "github.com/gin-gonic/contrib/secure"
)

func secureOptions() gin.HandlerFunc {
  return secure.Secure(secure.Options{
    //AllowedHosts:          []string{"example.com", "ssl.example.com"},
		SSLRedirect:           true,
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

/*func authProtect(c *gin.Context) gin.HandlerFunc {

}
*/
