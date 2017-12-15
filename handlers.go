package main

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/prometheus/config"
	"strings"
	"mime"
	"path"
)

func serveStatic(context *gin.Context, fileName string) {
	bytes, err := Asset(strings.Replace(fileName, "/", "", 1))
	if err != nil {
		context.AbortWithError(404, err)
		return
	}
	context.Data(200, mime.TypeByExtension(path.Ext(fileName)), bytes)
}

func staticFileHandler(fileName string) gin.HandlerFunc {
	return func(context *gin.Context) {
		serveStatic(context, fileName)
	}
}

func staticFilePathHandler(context *gin.Context) {
	serveStatic(context, context.Param("file"))
}

func saveConfig(context *gin.Context) {
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
	backupFilePath := backupFile()
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
		"Status":  "Success",
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
		"Status":  "Success",
		"Details": fmt.Sprintf("Config applied. %v", string(bytes)),
	})
}
