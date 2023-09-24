package main

import (
	"log"
	"patient/auth"
	"patient/config"
	"patient/controllers"
	"patient/middleware"
	"patient/models"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func init() {
	err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	// gin.SetMode(gin.DebugMode)
	r := gin.Default()
	r.SetTrustedProxies(nil)
	r.Use(gin.Logger())
	if gin.Mode() == gin.ReleaseMode {
		r.Use(middleware.SecurityMiddleware())
		r.Use(middleware.XssMiddleware())
	}
	r.Use(middleware.CorsMiddleware())
	r.Use(middleware.RateLimitMiddleware(rate.Every(1*time.Minute), 60)) // 60 requests per minute
	models.ConnectDatabase()

	v1 := r.Group("/api/v1")
	{
		v1.GET("/", controllers.Healthcheck)
		v1.POST("/file/upload", middleware.AuthenticateJWT(), controllers.UploadFile)
		v1.GET("/file/view/:id", middleware.AuthenticateJWT(), controllers.GetFile)
		v1.GET("/patients", middleware.AuthenticateJWT(), controllers.GetAllPatients)

		patient := v1.Group("/patient")
		{
			patient.POST("/login", auth.LoginHandler)
			patient.POST("/register", auth.PatientRegisterHandler)
			patient.GET("/:email", middleware.AuthenticateJWT(), controllers.GetPatientByEmail)
			patient.PUT("", middleware.AuthenticateJWT(), controllers.UpdatePatient)
			patient.DELETE("/:email", middleware.AuthenticateJWT(), controllers.DeletePatient)

		}
		form := patient.Group("/form", middleware.AuthenticateJWT())
		{
			form.POST("", middleware.AuthenticateJWT(), controllers.PatientFormSubmit)
			form.GET("/:form_id", middleware.AuthenticateJWT(), controllers.GetPatientFormByFormId)
		}
		patient.GET("/forms/:email", middleware.AuthenticateJWT(), controllers.GetPatientAllFormByPatientId)
		patient.GET("/forms", middleware.AuthenticateJWT(), controllers.GetAllForms)

	}
	if err := r.Run(":" + config.Cfg.ServerPort); err != nil {
		log.Fatal(err)
	}
}
