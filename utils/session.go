package utils

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetSessionID returns the session’s UUID, creating and saving one if absent.
func GetSessionID(c *gin.Context) string {
	sess := sessions.Default(c)
	id, ok := sess.Get("session_id").(string)
	if !ok || id == "" {
		// first visit: mint a new UUID
		id = uuid.NewString()
		sess.Set("session_id", id)
		sess.Save()
	}
	return id
}
