package database

import (
	"fmt"
	"log"

	"github.com/TFX0019/api-go-gds/pkg/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := config.GetEnv("DATABASE_URL", "")
	if dsn == "" {
		host := config.GetEnv("DB_HOST", "localhost")
		user := config.GetEnv("DB_USER", "postgres")
		password := config.GetEnv("DB_PASSWORD", "postgres")
		dbname := config.GetEnv("DB_NAME", "postgres")
		port := config.GetEnv("DB_PORT", "5432")
		sslmode := config.GetEnv("DB_SSLMODE", "disable")

		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", host, user, password, dbname, port, sslmode)
	}

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	log.Println("Database connection established")
}
