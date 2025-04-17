package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func EnsureSession() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID, err := c.Cookie("session_id")
		if err != nil || sessionID == "" {
			newSessionID := uuid.NewString()
			c.SetCookie("session_id", newSessionID, 60*60*24*7, "/", "", false, true) // 1 week
			c.Set("session_id", newSessionID)
		} else {
			c.Set("session_id", sessionID)
		}
		c.Next()
	}
}
