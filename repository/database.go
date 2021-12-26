package repository

import (
	"fmt"
	"gorm.io/gorm"
	"log"
)
import "gorm.io/driver/postgres"

var DB *gorm.DB

func ConnectDb(host string, user string, password string, dbName string, port string) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", host, user, password, dbName, port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}
	err = db.AutoMigrate(&Chat{}, &User{}, &Message{})
	if err != nil {
		log.Fatalf("error on auto migration: %v", err)
	}
	DB = db
	log.Printf("successfully connected to database")
}
