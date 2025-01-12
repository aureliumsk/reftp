package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func StatusResponse(c *gin.Context) {
	c.Next()
	switch status := c.Writer.Status(); status {
	case http.StatusInternalServerError:
		InternalServerError(c)
	case http.StatusNotFound:
		NotFound(c)
	}
}
