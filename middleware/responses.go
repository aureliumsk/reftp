package middleware

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"time"
)

func NotFound(c *gin.Context) {
	c.HTML(http.StatusNotFound, "notfound.html", nil)
}
func InternalServerError(c *gin.Context) {
	c.HTML(
		http.StatusInternalServerError,
		"internalerror.html",
		gin.H{"time": time.Now().Format("2006/01/02 - 15:04:05")},
	)
}
