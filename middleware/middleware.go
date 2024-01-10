package middleware

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/panithee/internship_day2/service"
)

func AuthorizeJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		const BEARER_SCHEMA = "Bearer "
		authHeader := c.GetHeader("Authorization")
		tokenString := authHeader[len(BEARER_SCHEMA):]
		token, err := service.JWTAuthService().ValidateToken(tokenString)

		if err != nil || !token.Valid {
			fmt.Println("Token validation error:", err)
			c.IndentedJSON(401, gin.H{"data": nil, "message": "Token is invalid"})
			c.Abort()
		}

		claims := token.Claims.(jwt.MapClaims)
		fmt.Println(claims)
		c.Set("claims", claims)
		c.Next()

	}
}
