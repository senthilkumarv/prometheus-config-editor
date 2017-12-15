package main

import (
	"github.com/gin-gonic/gin"
	"fmt"
	"net/http"
	"io/ioutil"
	"time"
	"github.com/prometheus/prometheus/config"
	"os"
	"net/url"
	"mime"
	"path/filepath"
	"strings"
)

var prometheusPath = os.Getenv("PROMETHEUS_PATH")
var backupFileFormat = "backup/prometheus-%v.yml"
var prometheusConfigFile = fmt.Sprintf("%v/prometheus.yml", prometheusPath)
var prometheusHostName = os.Getenv("PROMETHEUS_HOST_NAME")
var prometheusReloadEndpoint = fmt.Sprintf("http://%v/-/reload", prometheusHostName)

func main() {
	if prometheusPath == "" || prometheusHostName == "" {
		panic("PROMETHEUS_PATH & PROMETHEUS_HOST_NAME should be set")
	}
	router := gin.Default()
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

func serveStatic(context *gin.Context, fileName string) {
	bytes, err := Asset(strings.Replace(fileName, "/", "", 1))
	if err != nil {
		context.AbortWithError(404, err)
		return
	}
	context.Data(200, mime.TypeByExtension(filepath.Ext(fileName)), bytes)
}

func staticFileHandler(fileName string) gin.HandlerFunc {
	return func(context *gin.Context) {
		serveStatic(context, fileName)
	}
}

func staticFilePathHandler(context *gin.Context) {
	serveStatic(context, context.Param("file"))
}

func throwError(message string, err error) gin.H  {
	return gin.H{
		"Status": "Error",
		"Error": fmt.Sprintf("%s", err),
		"Details": message,
	}
}

func makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func saveConfig(context *gin.Context)  {
	bytes, err := context.GetRawData()
	if err != nil {
		context.JSON(400, throwError("Error retrieving config", err))
		return
	}
	contentBytes, err := ioutil.ReadFile(prometheusConfigFile)
	if err != nil {
		context.JSON(500, throwError("Error reading existing config", err))
		return
	}
	backupFilePath := fmt.Sprintf(backupFileFormat, makeTimestamp())
	err = ioutil.WriteFile(backupFilePath, contentBytes, 0644)
	if err != nil {
		context.JSON(500, throwError("Error saving backup file", err))
		return
	}
	_, err = config.Load(string(bytes))
	if err != nil {
		context.JSON(400, throwError("Invalid config. Please check syntax.", err))
		return
	}
	ioutil.WriteFile(prometheusConfigFile, bytes, 0644)
	context.JSON(200, gin.H{
		"Status": "Success",
		"Details": fmt.Sprintf("Backup file saved at %s", backupFilePath),
	})
}

func applyConfig(context *gin.Context) {
	response, err := http.PostForm(prometheusReloadEndpoint, url.Values{})

	if err != nil {
		context.JSON(400, throwError("Invalid config. Please check syntax.", err))
		return
	}

	bytes, err := ioutil.ReadAll(response.Body)

	if err != nil {
		context.JSON(400, throwError("Invalid config.", err))
		return
	}

	context.JSON(200, gin.H{
		"Status": "Success",
		"Details": fmt.Sprintf("Config applied. %v", string(bytes)),
	})
}