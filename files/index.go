package files

import (
	"reftp/db"
	"net/http"
	"github.com/gin-gonic/gin"
) 

func countFiles() (int, error) {
	var count int
	row := db.Conn.QueryRow("SELECT count(*) FROM files")
	if err := row.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func Index(c *gin.Context) {
	count, err := countFiles()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.HTML(http.StatusOK, "index.html", gin.H{
		"count": count,
	})
}
