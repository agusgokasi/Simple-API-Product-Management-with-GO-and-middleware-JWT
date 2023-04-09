package database

import (
	"eleventh-learn/models"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var (
	host     = "localhost"
	user     = "postgres"
	password = ""
	port     = "5432"
	dbname   = "db_go_sql"
	db       *gorm.DB
	err      error
)

func CreateAdmin() {
	// Check if admin user already exists
	var admin models.User
	if err := db.Where("email = ?", "admin@mail.com").First(&admin).Error; err == nil {
		// Admin user already exists, do nothing
		return
	}

	// Create admin user
	admin = models.User{
		FullName: "Admin",
		Email:    "admin@mail.com",
		Password: "123456",
		IsAdmin:  true,
	}

	if err := db.Create(&admin).Error; err != nil {
		log.Fatalf("failed to seed admin user: %v", err)
	}
}

func StartDB() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, password, dbname, port)

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "db_go_sql.",
			SingularTable: false,
		},
	})
	if err != nil {
		log.Fatal("error connecting to db", err)
	}

	log.Println("successfully connected to")
	db.Debug().AutoMigrate(models.User{}, models.Product{})
	CreateAdmin()
}

func GetDB() *gorm.DB {
	return db
}
