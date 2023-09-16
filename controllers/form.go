package controllers

import (
	"io"
	"net/http"
	"patient/models"
	"regexp"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func UploadFile(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unable to upload file: " + err.Error()})
		return
	}

	fileExtension := regexp.MustCompile(`\.[a-zA-Z0-9]+$`).FindString(fileHeader.Filename)
	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unable to open file: " + err.Error()})
		return
	}

	content, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unable to read file: " + err.Error()})
		return
	}

	fileCollection := models.Collection["files"].Database()
	bucket, err := gridfs.NewBucket(fileCollection, options.GridFSBucket().SetName("files"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	uploadStream, err := bucket.OpenUploadStream(fileHeader.Filename, options.GridFSUpload().SetMetadata(gin.H{"ext": fileExtension}))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	fieldId := uploadStream.FileID
	defer uploadStream.Close()
	fileSize, err := uploadStream.Write(content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "Image uploaded successfully",
		"image": gin.H{
			"id":   fieldId,
			"name": fileHeader.Filename,
			"size": fileSize,
		},
	})
}
