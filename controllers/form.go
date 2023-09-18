package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"patient/constant"
	"patient/models"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
		"msg": "file uploaded successfully",
		"file": gin.H{
			"id":   fieldId,
			"name": fileHeader.Filename,
			"size": fileSize,
		},
	})
}

func GetFile(c *gin.Context) {
	id := c.Param("id")
	if len(id) <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id format"})
		return
	}

	var fileMetadata bson.M
	fileCollection := models.Collection["files"].Database()
	if err := fileCollection.Collection("files.files").FindOne(c, bson.D{primitive.E{Key: "_id", Value: objID}}).Decode(&fileMetadata); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.JSON(http.StatusOK, gin.H{"message": "no data present for this id", "data": nil})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "error on getting data: " + err.Error()})
		return
	}
	var fileBuffer bytes.Buffer
	bucket, err := gridfs.NewBucket(fileCollection, options.GridFSBucket().SetName("files"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	fileStream, err := bucket.OpenDownloadStream(objID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	defer fileStream.Close()
	_, err = io.Copy(c.Writer, fileStream)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error copying file to response: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "successfully",
		"file": gin.H{
			"data": fileBuffer.Bytes(),
		},
	})
}

func PatientFormSubmit(c *gin.Context) {
	var form models.FormRequest
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var newForm models.Form
	if err := copier.Copy(&newForm, &form); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("copier error: %v", err)})
		return
	}

	collection := models.Collection["forms"]
	id, _ := c.Get(constant.PATIENT_ID_CONTEXT)
	email, _ := c.Get(constant.EMAIL_CONTEXT)

	newForm.PatientId = id.(string)
	newForm.PatientEmail = email.(string)
	now := time.Now().Local()
	newForm.CreatedAt = now
	newForm.UpdatedAt = now
	result, err := collection.InsertOne(c, newForm)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("could not save form: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "form submitted successfully",
		"form_id": result.InsertedID,
	})
}

func GetPatientFormByFormId(c *gin.Context) {
	formId := c.Param("form_id")
	if len(formId) <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	objID, err := primitive.ObjectIDFromHex(formId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id format"})
		return
	}

	var result models.Form
	filter := bson.M{"_id": objID}
	collection := models.Collection["forms"]
	if err = collection.FindOne(c, filter).Decode(&result); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.JSON(http.StatusOK, gin.H{"message": "no data found", "data": nil})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error on getting data: " + err.Error()})
		return
	}

	userRole, _ := c.Get(constant.ROLE_CONTEXT)
	patientId, _ := c.Get(constant.PATIENT_ID_CONTEXT)
	if patientId != result.PatientId && userRole != constant.ADMIN_ROLE {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "you have no access to get patient data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}

func GetPatientAllFormByPatientId(c *gin.Context) {
	incomingEmail := c.Param("email")
	if len(incomingEmail) <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	userRole, _ := c.Get(constant.ROLE_CONTEXT)
	patientEmail, _ := c.Get(constant.EMAIL_CONTEXT)
	if incomingEmail != patientEmail && userRole != constant.ADMIN_ROLE {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "you have no access to get patient data"})
		return
	}

	filterJson := c.Query("json")
	filter := models.Filter{
		Page:  1,
		Limit: 10,
	}
	if len(filterJson) > 0 {
		if err := json.Unmarshal([]byte(filterJson), &filter); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid filter request passed"})
			return
		}
	}

	if filter.Page < 1 || filter.Limit < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "filter values must be greater than 0"})
		return
	}

	countOpts := options.Count()
	collection := models.Collection["forms"]
	totalRecords, err := collection.CountDocuments(c, bson.M{"patient_email": incomingEmail}, countOpts)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unable to get record count: " + err.Error()})
		return
	}

	skip := (filter.Page - 1) * filter.Limit
	opts := options.Find()
	opts.SetSort(bson.D{primitive.E{Key: "updated_at", Value: -1}})
	opts.SetLimit(filter.Limit)
	opts.SetSkip(int64(skip))
	var result []models.Form
	cursor, err := collection.Find(c, bson.M{"patient_email": incomingEmail}, opts)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.JSON(http.StatusOK, gin.H{"message": "no data found", "data": nil})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error on getting data: " + err.Error()})
		return
	}

	for cursor.Next(c) {
		var content models.Form
		err := cursor.Decode(&content)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error on getting data: " + err.Error()})
			return
		}

		result = append(result, content)
	}
	metadata := computeMetadata(totalRecords, filter.Page, filter.Limit)
	c.JSON(http.StatusOK, gin.H{
		"data":     result,
		"metadata": metadata,
	})
}

func GetAllForms(c *gin.Context) {
	userRole, _ := c.Get(constant.ROLE_CONTEXT)
	if userRole != constant.ADMIN_ROLE {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "you have no access to get forms"})
		return
	}

	filterJson := c.Query("json")
	filter := models.Filter{
		Page:  1,
		Limit: 10,
	}
	if len(filterJson) > 0 {
		if err := json.Unmarshal([]byte(filterJson), &filter); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid filter request passed"})
			return
		}
	}

	if filter.Page < 1 || filter.Limit < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "filter values must be greater than 0"})
		return
	}

	countOpts := options.Count()
	collection := models.Collection["forms"]
	totalRecords, err := collection.CountDocuments(c, bson.M{}, countOpts)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unable to get record count: " + err.Error()})
		return
	}

	skip := (filter.Page - 1) * filter.Limit
	opts := options.Find()
	opts.SetSort(bson.D{primitive.E{Key: "updated_at", Value: -1}})
	opts.SetLimit(filter.Limit)
	opts.SetSkip(int64(skip))
	var result []models.Form
	cursor, err := collection.Find(c, bson.M{}, opts)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.JSON(http.StatusOK, gin.H{"message": "no data found", "data": nil})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error on getting data: " + err.Error()})
		return
	}

	for cursor.Next(c) {
		var content models.Form
		err := cursor.Decode(&content)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error on getting data: " + err.Error()})
			return
		}

		result = append(result, content)
	}
	metadata := computeMetadata(totalRecords, filter.Page, filter.Limit)
	c.JSON(http.StatusOK, gin.H{
		"data":     result,
		"metadata": metadata,
	})
}
