package files

import (
	"reftp/db"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"mime/multipart"
	"fmt"
	"os"
	"path/filepath"
)

type inputFile struct{
	name string
	Desc string `form:"desc"`
	File *multipart.FileHeader `form:"file" binding:"required"`
}

func createFile(c *gin.Context, f inputFile) error {
	path := filepath.Join(Directory, f.name)
	exists := true
	if _, err := os.Stat(path); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		exists = false
	}
	if err := c.SaveUploadedFile(f.File, path); err != nil {
		return err
	}
	if !exists {
		_, err := db.Conn.Exec("INSERT INTO files (name, description) VALUES (?, ?)", f.name, f.Desc)
		if err != nil {
			os.Remove(path)
			return err
		}
		return nil
	}
	_, err := db.Conn.Exec(
		"UPDATE files SET (description, modifiedAt) = (?, datetime()) WHERE name = ?",
		f.Desc, f.name,
	)
	return err
}

func FileModifyRoute(c *gin.Context) {
	var f inputFile
	name := c.Param("name")
	if err := c.MustBindWith(&f, binding.FormMultipart); err != nil {
		return
	}
	f.name = name
	if err := createFile(c, f); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Redirect(http.StatusFound, fmt.Sprintf("/files/%s", name))
}
