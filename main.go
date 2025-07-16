package main

import (
	"go-crud-oapi/config"
	"go-crud-oapi/internal/controller"
	"go-crud-oapi/internal/db"
	"go-crud-oapi/internal/model"
	"go-crud-oapi/internal/repository"
	"go-crud-oapi/internal/router"
	"go-crud-oapi/internal/service"
	"go-crud-oapi/pkg/logger"
	"log"
	"net/http"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func main() {
	logger.Init()

	cfg := config.Load()
	dbConn := db.Init(cfg)
	seedAdminUser(dbConn)

	repo := repository.NewUserRepository(dbConn)
	svc := service.NewUserService(repo)
	userController := controller.NewUserController(svc)
	authController := controller.NewAuthController(repo)

	// Inject all controllers to router
	r := router.NewRouter(userController, authController)

	port := cfg.ServerPort
	//port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // fallback default
	}

	log.Printf("🚀 Server running on : %s", port)
	http.ListenAndServe(":"+port, r)

}

func seedAdminUser(db *gorm.DB) {
	var count int64
	db.Model(&model.User{}).Where("email = ?", "admin@example.com").Count(&count)
	if count == 0 {
		hash, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		db.Create(&model.User{
			Name:     "Admin User",
			Email:    "admin@example.com",
			Password: string(hash),
			Role:     "admin",
		})
	}
}
