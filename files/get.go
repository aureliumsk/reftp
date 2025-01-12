package files

import (
	"reftp/db"
	"github.com/gin-gonic/gin"
	"path/filepath"
	"strings"
	"os"
	"net/http"
	"mime"
	"time"
	"database/sql"
)

type File struct{
	Name string
	Desc string
	CreatedAt time.Time
	ModifiedAt sql.NullTime
}

var Directory string

func getFile(name string) (*File, error) {
	var file File
	row := db.Conn.QueryRow("SELECT description, createdAt, modifiedAt FROM files WHERE name = ?", name)
	if err := row.Scan(&file.Desc, &file.CreatedAt, &file.ModifiedAt); err != nil {
		// TODO: Do this properly
		return nil, err
	}
	file.Name = name
	file.CreatedAt = file.CreatedAt.Local()
	if file.ModifiedAt.Valid {
		file.ModifiedAt.Time = file.ModifiedAt.Time.Local()
	}
	return &file, nil
}
func guessMime(file *os.File) string {
	mimeType := mime.TypeByExtension(filepath.Ext(file.Name()))
	if mimeType == "" {
		return "application/octet-stream"
	}
	return mimeType
}
func FileServeRoute(c *gin.Context) {
	name := c.Param("name")
	path := filepath.Join(Directory, name)
	if !strings.HasPrefix(path, Directory) {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer file.Close()
	// Handle error (it probably shouldn't occur)
	stat, err := file.Stat()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.DataFromReader(http.StatusOK, stat.Size(), guessMime(file), file, nil)
}
func FilePageRoute(c *gin.Context) {
	file, err := getFile(c.Param("name"))
	if err != nil {
		if err == sql.ErrNoRows {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.HTML(http.StatusOK, "file.html", file)
}



