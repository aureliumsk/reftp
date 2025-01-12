package files

import (
	"reftp/db"
	"net/http"
	"github.com/gin-gonic/gin"
	"os"
	"path/filepath"
)

func deleteFile(name string) (bool, error) {
	if err := os.Remove(filepath.Join(Directory, name)); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	_, err := db.Conn.Exec("DELETE FROM files WHERE name = ?", name)
	if err != nil {
		return false, err
	}
	return true, nil
}

func FileDeleteRoute(c *gin.Context) {
	name := c.Param("name")
	found, err := deleteFile(name)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if !found {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	// TODO: invent a response
	c.Redirect(http.StatusFound, "/")
}
