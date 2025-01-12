package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"path/filepath"
	"reftp/middleware"
	"reftp/db"
	"reftp/files"
)

var databasePath = flag.String("db", "prod.db", "database path")
var templatePattern = flag.String("templates", "./pages/*", "template pattern")
var doCreateTables = flag.Bool("create", false, "create tables")

func main() {
	flag.StringVar(&files.Directory, "dir", "", "storage directory")
	flag.Parse()
	if files.Directory == "" {
		dir, err := os.Getwd()
		if err != nil {
			log.Fatalf("can't get wd: %v\n", err)
		}
		files.Directory = filepath.Join(dir, "storage")
	} else {
		path, err := filepath.Abs(files.Directory)
		if err != nil {
			log.Fatalf("can't get abs of %s: %v\n", files.Directory, err)
		}
		files.Directory = path
	}
	if err := db.Init(*databasePath); err != nil {
		log.Fatalf("can't init the database: %v\n", err)
	}
	if *doCreateTables {
		if err := db.CreateTables(); err != nil {
			log.Fatalf("can't create tables: %v\n", err)
		}
		log.Println("succesfully dropped & created all of the tables")
	}
	router := gin.Default()
	router.LoadHTMLGlob(*templatePattern)
	router.Use(middleware.StatusResponse)
	router.NoRoute(middleware.NotFound)
	router.GET("/", files.Index)
	router.GET("/files/:name", files.FilePageRoute)
	router.GET("/files/download/:name", files.FileServeRoute)
	router.POST("/files/:name", files.FileModifyRoute)
	router.DELETE("/files/:name", files.FileDeleteRoute)
	router.Run()
}
