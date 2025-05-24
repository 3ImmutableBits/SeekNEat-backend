package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/3ImmutableBits/SeekNEat-backend/models"
	"github.com/go-chi/jwtauth"
)

type newMealRequest struct {
	Latitude       float64 `json:"latitude"`
	Longitude      float64 `json:"longitude"`
	AvailableSpots uint    `json:"available_spots"`
	Name           string  `json:"name"`
	Description    string  `json:"description"`
	Price          string  `json:"price"`
	Timestamp      int64   `json:"timestamp"`
}

func newMealHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	data := newMealRequest{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		sendError("Invalid json", w)
		return
	}

	if data.Latitude > 90 || data.Latitude < -90 || data.Longitude > 180 || data.Longitude < -180 {
		sendError("Invalid coords", w)
		return
	}
	if data.AvailableSpots == 0 {
		sendError("Availale spots must be greater than 0", w)
		return
	}
	if data.Name == "" {
		sendError("Name cannot be empty", w)
		return
	}
	if data.Description == "" {
		sendError("Description cannot be empty", w)
	}
	if data.Price == "" {
		sendError("Price cannot be empty", w)
	}
	if data.Timestamp < time.Now().Unix() {
		sendError("Time cannot be in the past", w)
	}

	_, claims, _ := jwtauth.FromContext(r.Context())
	meal := models.Meal{
		Location:       models.Coords{Latitude: data.Latitude, Longitude: data.Longitude},
		HostId:         uint(claims["userId"].(float64)),
		AvailableSpots: data.AvailableSpots,
		Name:           data.Name,
		Description:    data.Description,
		Price:          data.Price,
	}
	result := db.Create(&meal)

	w.WriteHeader(http.StatusOK)
	if result.Error != nil {
		sendError("Database error", w)
		return
	}
	json.NewEncoder(w).Encode(map[string]any{
		"success": true,
		"error":   "",
	})
}

type joinMealRequest struct {
	MealId uint `json:"meal_id"`
}

func joinMealHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	data := joinMealRequest{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		sendError("Invalid json", w)
		return
	}

	meal := models.Meal{}
	if err := db.Where("ID = ?", data.MealId).First(&meal).Error; err != nil {
		sendError("Meal doesn't exist", w)
		return
	}

	_, claims, _ := jwtauth.FromContext(r.Context())
	u := models.User{}
	db.Where("ID = ?", claims["userId"]).First(&u)

	result := db.Model(&meal).Association("Clients").Append(&u)
	w.WriteHeader(http.StatusOK)
	if result != nil {
		sendError("Database error", w)
		return
	}
	json.NewEncoder(w).Encode(map[string]any{
		"success": true,
		"error":   "",
	})
}

type fetchMealRequest struct {
	Query string `json:"query"`
}

func fetchMealHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	data := fetchMealRequest{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		sendError("Invalid json", w)
		return
	}

	var meals []models.Meal
	data.Query = fmt.Sprintf("%%%s%%", data.Query)
	queryResult := db.Where("name LIKE ? OR description LIKE ?", data.Query, data.Query).Find(&meals)
	if queryResult.Error != nil {
		sendError("Database error", w)
		return
	}

	var result []map[string]any

	for _, meal := range meals {
		if len(meal.Clients) < int(meal.AvailableSpots) && time.Unix(meal.Timestamp, 0).After(time.Now()) {
			result = append(result, map[string]any{
				"id":             meal.ID,
				"latitude":       meal.Location.Latitude,
				"longitude":      meal.Location.Longitude,
				"host_id":        meal.HostId,
				"timestamp":      meal.Timestamp,
				"price":          meal.Price,
				"name":           meal.Name,
				"description":    meal.Description,
				"spots":          meal.AvailableSpots,
				"occupied_spots": len(meal.Clients),
			})
		}
	}

	json.NewEncoder(w).Encode(map[string]any{
		"success": true,
		"result":  result,
	})
}

type deleteMealRequest struct {
	MealId uint `json:"meal_id"`
}

func deleteMealHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	data := deleteMealRequest{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		sendError("Invalid json", w)
		return
	}

	meal := models.Meal{}
	if err := db.Where("ID = ?", data.MealId).First(&meal).Error; err != nil {
		sendError("Meal doesn't exist", w)
		return
	}

	_, claims, _ := jwtauth.FromContext(r.Context())
	u := models.User{}
	db.Where("ID = ?", claims["userId"]).First(&u)

	result := db.Delete(&meal)
	w.WriteHeader(http.StatusOK)
	if result.Error != nil {
		sendError("Database error", w)
		return
	}
	json.NewEncoder(w).Encode(map[string]any{
		"success": true,
		"error":   "",
	})
}
