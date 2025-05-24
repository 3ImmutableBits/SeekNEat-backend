package api

import (
	"encoding/json"
	"net/http"

	"github.com/3ImmutableBits/SeekNEat-backend/models"
	"github.com/go-chi/jwtauth"
)

type changeUserRequest struct {
	Username string `gorm:"username"`
	Password string `gorm:"password"`
	Email    string `gorm:"email"`
}

func changeUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	data := changeUserRequest{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		sendError("Invalid json", w)
		return
	}

	u := models.User{}
	_, claims, _ := jwtauth.FromContext(r.Context())
	db.Where("ID = ?", claims["userId"]).First(&u)

	flag := false

	if data.Username != "" {
		flag = true
		u.Username = data.Username
	}
	if data.Email != "" {
		flag = true
		u.Email = data.Email
	}
	if data.Password != "" {
		flag = true
		u.SetPassword(data.Password)
	}

	if !flag {
		sendError("All fields are empty", w)
	}

	result := db.Save(&u)

	if result.Error != nil {
		sendError("Database error", w)
	}
	json.NewEncoder(w).Encode(map[string]any{
		"success": true,
		"error":   "",
	})
}
