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
	// cache.InitRedis()

	//gin.SetMode(gin.ReleaseMode)
	gin.SetMode(gin.DebugMode)
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
		v1.POST("/upload", controllers.UploadFile)

		patient := v1.Group("/patient")
		{
			patient.POST("/login", auth.LoginHandler)
			patient.POST("/register", auth.PatientRegisterHandler)
			patient.GET("/:email", middleware.AuthenticateJWT(), controllers.GetPatientByEmail)
			patient.PUT("", middleware.AuthenticateJWT(), controllers.UpdatePatient)
			patient.DELETE("/:email", middleware.AuthenticateJWT(), controllers.DeleteBook)

			// v1.GET("/books", middleware.APIKeyAuthMiddleware(), controllers.FindBooks)
			// v1.POST("/books", middleware.APIKeyAuthMiddleware(), middleware.AuthenticateJWT(), controllers.CreateBook)
			// v1.GET("/books/:id", middleware.APIKeyAuthMiddleware(), controllers.FindBook)
			// v1.PUT("/books/:id", middleware.APIKeyAuthMiddleware(), controllers.UpdateBook)
			// v1.DELETE("/books/:id", middleware.APIKeyAuthMiddleware(), controllers.DeleteBook)
		}
		form := patient.Group("/form", middleware.AuthenticateJWT())
		{
			form.POST("", controllers.PatientFormSubmit)
			form.GET("/:form_id", controllers.GetPatientFormByFormID)
		}
	}
	// r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	if err := r.Run(":" + config.Cfg.ServerPort); err != nil {
		log.Fatal(err)
	}
}
