package db

import (
	"log"
	"os"

	"github.com/jfilipedias/codepix-go/domain/model"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func ConnectDB() *gorm.DB {
	var dns string
	var db *gorm.DB
	var err error
	var newLogger logger.Interface

	env := os.Getenv("env")

	if os.Getenv("debug") == "true" {
		newLogger = logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{})
	}

	if env != "test" {
		dns = os.Getenv("dns")
		db, err = gorm.Open(postgres.Open(dns), &gorm.Config{
			Logger: newLogger,
		})
	} else {
		dns = os.Getenv("dnsTest")
		db, err = gorm.Open(sqlite.Open(dns), &gorm.Config{
			Logger: newLogger,
		})
	}

	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
		panic(err)
	}

	if os.Getenv("AutoMigrateDb") == "true" {
		db.AutoMigrate(&model.Bank{}, &model.Account{}, &model.PixKey{}, &model.Transaction{})
	}

	return db
}
