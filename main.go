package main

import (
	"log"
	"net/http"

	"github.com/3ImmutableBits/SeekNEat-backend/api"
	"github.com/3ImmutableBits/SeekNEat-backend/config"
	"github.com/3ImmutableBits/SeekNEat-backend/models"
	"github.com/go-chi/chi/v5"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(sqlite.Open(config.DBFile), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database: ", err)
	}
	db.AutoMigrate(&models.User{}, &models.Meal{})

	r := chi.NewRouter()

	api.AddRoutes(r, db)

	log.Println("Server started")
	err = http.ListenAndServe(config.Port, r)
	if err != nil {
		log.Fatal(err)
	}
}
