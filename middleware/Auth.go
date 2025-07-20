package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"recipes-api/model"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("AuthMiddleware called for %s", c.Request.URL.Path)
		tokenValue := c.GetHeader("Authorization")
		claims := &model.Claims{}
		tkn, err := jwt.ParseWithClaims(tokenValue, claims,
			func(token *jwt.Token) (interface{}, error) {
				return []byte(os.Getenv("JWT_SECRET")), nil
			})
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		if tkn == nil || !tkn.Valid {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		c.Next()
	}
}
