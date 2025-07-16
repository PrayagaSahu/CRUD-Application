package db

import (
	"database/sql"
	"fmt"
	"log"

	"go-crud-oapi/config"
	"go-crud-oapi/internal/model"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Init initializes the database connection and runs migrations.
func Init(cfg *config.Config) *gorm.DB {

	//dbDriver := cfg.DBDriver
	host := cfg.DBHost
	port := cfg.DBPort
	user := cfg.DBUser
	password := cfg.DBPassword
	dbName := cfg.DBName

	// Step 1: Connect to postgres system DB to check/create target DB
	systemDSN := fmt.Sprintf("host=%s user=%s password=%s dbname=postgres port=%s sslmode=disable", host, user, password, port)
	systemDB, err := sql.Open("postgres", systemDSN)
	if err != nil {
		log.Fatal("❌ System DB connection failed:", err)
	}
	defer systemDB.Close()

	// Step 2: Check if DB exists
	var exists bool
	err = systemDB.QueryRow("SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)", dbName).Scan(&exists)
	if err != nil {
		log.Fatal("❌ Failed to check DB existence:", err)
	}

	// Step 3: Create DB if not exists
	if !exists {
		_, err = systemDB.Exec("CREATE DATABASE " + dbName)
		if err != nil {
			log.Fatalf("❌ Failed to create DB %s: %v", dbName, err)
		}
		log.Printf("✅ Database '%s' created", dbName)
	}

	// Step 4: Connect to target DB using GORM's DSN
	appDSN := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, password, dbName, port)
	gormDB, err := gorm.Open(postgres.Open(appDSN), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ GORM init failed:", err)
	}

	// Step 5: AutoMigrate
	if err := gormDB.AutoMigrate(&model.User{}); err != nil {
		log.Fatal("❌ AutoMigrate failed:", err)
	}

	log.Println("✅ Database connected and migrated")
	return gormDB
}
