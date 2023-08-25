package main

import (
	"fmt"
	"github.com/wpcodevo/two_factor_golang/models"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/wpcodevo/two_factor_golang/controllers"
	"github.com/wpcodevo/two_factor_golang/routes"
	"gorm.io/driver/mysql" // Menggunakan driver MySQL
	"gorm.io/gorm"
)

var (
	DB                  *gorm.DB
	server              *gin.Engine
	AuthController      controllers.AuthController
	AuthRouteController routes.AuthRouteController
)

func init() {
	var err error
	dsn := "rafza:GantengR28FBunga@tcp(127.0.0.1:3306)/otp_db?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	DB.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal("Failed to connect to the Database")
	}
	fmt.Println("🚀 Connected Successfully to the Database")

	AuthController = controllers.NewAuthController(DB)
	AuthRouteController = routes.NewAuthRouteController(AuthController)

	server = gin.Default()
}

func main() {
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:3000"}
	corsConfig.AllowCredentials = true

	server.Use(cors.New(corsConfig))

	router := server.Group("/api")
	router.GET("/healthchecker", func(ctx *gin.Context) {
		message := "Welcome to Two-Factor Authentication with Golang"
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": message})
	})

	AuthRouteController.AuthRoute(router)
	log.Fatal(server.Run(":8000"))
}
