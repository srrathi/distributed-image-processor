package database

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/srrathi/distributed-image-processor/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	Host     string
	Port     string
	Password string
	User     string
	DBName   string
	SSLMode  string
}

func getDatabaseConfig() (*Config, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, err
	}

	config := &Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASSWORD"),
		User:     os.Getenv("DB_USER"),
		DBName:   os.Getenv("DB_DATABASE"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}
	return config, nil
}

func NewConnection() (*gorm.DB, error) {
	config, err := getDatabaseConfig()
	if err != nil {
		return nil, err
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	log.Println("Connected to Postgres")

	err = MigrateDatabases(db)
	if err != nil {
		return db, err
	}

	return db, nil
}

func MigrateDatabases(db *gorm.DB) error {
	err := db.AutoMigrate(&models.JobStatus{})
	if err != nil {
		log.Fatal(err)
		return err
	}

	err = db.AutoMigrate(&models.JobErrors{})
	if err != nil {
		log.Fatal(err)
		return err
	}

	err = db.AutoMigrate(&models.StoreData{})
	if err != nil {
		log.Fatal(err)
		return err
	}

	err = db.AutoMigrate(&models.StoreVisits{})
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}
