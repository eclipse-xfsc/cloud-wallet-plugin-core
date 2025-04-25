package core

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func AuthMiddleware(dataFetcher ...dataFetcher) gin.HandlerFunc {
	var oidcProv *oidcProvider
	if len(dataFetcher) > 0 {
		oidcProv = getOidcProvider(dataFetcher[0])
	} else {
		oidcProv = getOidcProvider()
	}
	return func(c *gin.Context) {
		token := getTokenValue(c)
		user, err := oidcProv.GetUser(token)
		if user != nil {
			c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), UserKey, user))
			c.Next()
		} else {
			if err != nil {
				_ = c.Error(err)
			}
			c.JSON(http.StatusForbidden, gin.H{"message": UserNotFound})
			c.Abort()
		}
	}
}

func getTokenValue(c *gin.Context) string {
	token := c.Request.Header.Get("Authorization")
	res, _ := strings.CutPrefix(token, "Bearer ")
	return res
}
