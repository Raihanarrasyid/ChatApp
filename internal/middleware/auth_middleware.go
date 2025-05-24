package middleware

import (
	"ChatApp/configs"
	"errors"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims

}

func AuthMiddleware(config *configs.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		fmt.Print("Checkpoint 1 \n")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(401, gin.H{"error": "token is required"})
			c.Abort()
			return
		}

		fmt.Print("Checkpoint 2 \n")

		token := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := parseToken(token, config.JwtSecret)
		if err != nil {
			c.JSON(401, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		fmt.Print("Checkpoint 3 \n")
		
		c.Set("user_id", claims.UserID)

		c.Next()
	}
}

func parseToken(tokenString string, secretkey string) (*Claims, error) {

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretkey), nil
	})
	
	fmt.Print("Checkpoint 4 \n")
	if err != nil {
		return nil, errors.New("invalid token, " + "token expired or invalid")
	}
	
	
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token, " + "token expired or invalid")
	}


	return claims, nil
}