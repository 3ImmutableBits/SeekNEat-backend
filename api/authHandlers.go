package api

import (
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/3ImmutableBits/SeekNEat-backend/models"
)

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	data := loginRequest{}
	err := json.NewDecoder(r.Body).Decode(&data)

	if err != nil {
		json.NewEncoder(w).Encode(map[string]any{
			"success": false,
			"error":   "Invalid json",
		})
		return
	}

	user := models.User{}
	if err := db.Where("email = ? OR username = ?", data.Username, data.Username).First(&user).Error; err == nil &&
		user.CheckPassword(data.Password) {
		_, tokenString, _ := tokenAuth.Encode(map[string]any{
			"userId": user.ID,
		})

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]any{
			"success": true,
			"error":   "",
			"token":   tokenString,
		})
		return
	}
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(map[string]any{
		"success": false,
		"error":   "Invalid credentials",
		"token":   "",
	})
}

type registerRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	data := registerRequest{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]any{
			"success": false,
			"error":   "Invalid json",
		})
		return
	}

	if data.Username == "" {
		json.NewEncoder(w).Encode(map[string]any{
			"success": false,
			"error":   "Username cannot be empty",
		})
		return
	}
	if len(data.Password) < 8 {
		json.NewEncoder(w).Encode(map[string]any{
			"success": false,
			"error":   "Password must have at least 8 characters",
		})
		return
	}
	var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\.[a-zA-Z0-9-.]+$`) // Basic email regex
	if !emailRegex.MatchString(data.Email) {
		json.NewEncoder(w).Encode(map[string]any{
			"success": false,
			"error":   "Invalid email",
		})
		return
	}

	user := models.User{}
	if err := db.Where("email = ?", data.Email).First(&user).Error; err == nil {
		json.NewEncoder(w).Encode(map[string]any{
			"success": false,
			"error":   "Email already belongs to an account",
		})
		return
	}
	if err := db.Where("username = ?", data.Username).First(&user).Error; err == nil {
		json.NewEncoder(w).Encode(map[string]any{
			"success": false,
			"error":   "Username already belongs to an account",
		})
		return
	}

	user = models.User{Username: data.Username, Email: data.Email}
	user.SetPassword(data.Password)

	result := db.Create(&user)
	if result.Error != nil {
		json.NewEncoder(w).Encode(map[string]any{
			"success": false,
			"error":   "Database error",
		})
		return
	}
	json.NewEncoder(w).Encode(map[string]any{
		"success": true,
		"error":   "",
	})
}
