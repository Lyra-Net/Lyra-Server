package middleware

import (
	"fmt"
	"identity-service/redisconn"
	"identity-service/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		token := strings.TrimPrefix(auth, "Bearer ")
		claims, err := utils.ParseAccessToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		userID := uint(claims["user_id"].(float64))
		jti := claims["jti"].(string)

		key := fmt.Sprintf("blacklist:access:%s", jti)
		_, err = redisconn.Client.Get(redisconn.Ctx, key).Result()
		if err != redis.Nil {

			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token has been revoked"})
			return
		}

		iat := int64(claims["iat"].(float64))
		changePassAt, err := redisconn.GetChangePassAt(userID)
		if err != nil && err != redis.Nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
			return
		}
		if err == nil && iat < changePassAt {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token has been revoked due to password change"})
			return
		}

		c.Set("user_id", userID)
		c.Next()
	}
}
