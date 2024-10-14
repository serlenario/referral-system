package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/serlenario/referral-system/internal/config"
	"github.com/serlenario/referral-system/internal/controllers"
	"github.com/serlenario/referral-system/internal/middleware"
	"github.com/serlenario/referral-system/internal/models"
	"github.com/serlenario/referral-system/internal/repositories"
	"github.com/serlenario/referral-system/internal/services"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	_ "github.com/serlenario/referral-system/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	cfg := config.LoadConfig()

	dsn := "host=" + cfg.DBHost + " user=" + cfg.DBUser + " password=" + cfg.DBPassword + " dbname=" + cfg.DBName + " port=" + cfg.DBPort + " sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	if err := db.AutoMigrate(&models.User{}, &models.Referral{}); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	userRepo := repositories.NewUserRepository(db)
	referralRepo := repositories.NewReferralRepository(db)
	userService := services.NewUserService(userRepo, referralRepo, cfg.JWTSecret)
	userController := controllers.NewUserController(userService)

	router := gin.Default()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.POST("/register", userController.Register)
	router.POST("/login", userController.Login)
	router.POST("/register_with_referral", userController.RegisterWithReferral)
	router.GET("/referral_code", userController.GetReferralCodeByEmail)

	authorized := router.Group("/")
	authorized.Use(middleware.JWTMiddleware(cfg.JWTSecret))
	{
		authorized.POST("/referral_code", userController.CreateReferralCode)
		authorized.DELETE("/referral_code", userController.DeleteReferralCode)
		authorized.GET("/referrals", userController.GetReferrals)
	}

	log.Println("Server running on port 8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}
