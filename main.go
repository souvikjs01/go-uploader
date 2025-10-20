package main

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.POST("/upload", handleUpload)

	r.Run(":8000")

}

func handleUpload(c *gin.Context) {
	form, err := c.MultipartForm()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	files := form.File["files"]

	for key, file := range files {
		err := saveFile(file, key)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "file uploaded successfully",
	})

}

func saveFile(file *multipart.FileHeader, key int) error {
	ext := strings.Split(file.Filename, ".")
	if ext[1] != "png" && ext[1] != "jpg" {
		return fmt.Errorf("pass valid file format png or jpg")
	}

	src, err := file.Open()

	if err != nil {
		return err
	}

	defer src.Close()

	destPath := filepath.Join("./dest", ext[0]+strconv.Itoa(key)+"."+ext[1])

	dst, err := os.Create(destPath)
	if err != nil {
		return nil
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)

	if err != nil {
		return err
	}

	return nil

}
