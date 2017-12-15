package main

import (
	"os"
	"fmt"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/gzip"
)

var prometheusPath = os.Getenv("PROMETHEUS_PATH")
var backupDirectory = os.Getenv("BACKUP_DIRECTORY")
var backupFileFormat = "%v/prometheus-%v.yml"
var prometheusConfigFile = fmt.Sprintf("%v/prometheus.yml", prometheusPath)
var prometheusHostName = os.Getenv("PROMETHEUS_HOST_NAME")
var prometheusReloadEndpoint = fmt.Sprintf("http://%v/-/reload", prometheusHostName)

func main() {
	if prometheusPath == "" || prometheusHostName == "" || backupDirectory == "" {
		panic("PROMETHEUS_PATH & PROMETHEUS_HOST_NAME & BACKUP_DIRECTORY should be set")
	}
	setupBackupFolder()
	router := gin.Default()
	router.Use(gzip.Gzip(gzip.DefaultCompression))
	authorized := router.Group("/", gin.BasicAuth(gin.Accounts{
		"admin": "admin123",
	}))
	authorized.GET("/public/*file", staticFilePathHandler)
	authorized.StaticFile("/load", prometheusConfigFile)
	authorized.GET("/favicon.ico", staticFileHandler("favicon.ico"))
	authorized.GET("/", staticFileHandler("editor.html"))
	authorized.POST("/save", saveConfig)
	authorized.POST("/apply", applyConfig)
	router.Run(":8000")
}

func setupBackupFolder() {
	info, err := os.Stat(backupDirectory)
	if err != nil || !info.IsDir() {
		fmt.Println("Backup directory does not exists. Attempting to create...")
		if err := os.Mkdir(backupDirectory, 0755); err != nil {
			panic(fmt.Sprintf("Error creating backup directory. %v", err))
		}
	}
}

func backupFile() string {
	timeStamp := time.Now().UnixNano() / int64(time.Millisecond)
	return fmt.Sprintf(backupDirectory, backupFileFormat, timeStamp)
}

func throwError(message string, err error) gin.H {
	return gin.H{
		"Status":  "Error",
		"Error":   fmt.Sprintf("%s", err),
		"Details": message,
	}
}
