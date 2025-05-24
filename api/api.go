package api

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httputil"

	"github.com/3ImmutableBits/SeekNEat-backend/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth"
	"gorm.io/gorm"
)

var tokenAuth *jwtauth.JWTAuth

var db *gorm.DB

func sendError(err string, w http.ResponseWriter) {
	json.NewEncoder(w).Encode(map[string]any{
		"success": false,
		"error":   err,
	})
}

func AddRoutes(appRouter chi.Router, dbLocal *gorm.DB) {
	db = dbLocal

	r := chi.NewRouter()

	tokenAuth = jwtauth.New("HS256", []byte(config.SecretKey), nil)

	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var bodyCopy []byte
			if r.Body != nil {
				bodyCopy, _ = io.ReadAll(r.Body)
				r.Body = io.NopCloser(bytes.NewBuffer(bodyCopy))
			}

			dump, err := httputil.DumpRequest(r, true)
			if err == nil {
				log.Printf("HTTP Request Dump:\n%s\n", dump)
			} else {
				log.Printf("Error dumping request: %v", err)
			}

			// Reset body for downstream handlers
			r.Body = io.NopCloser(bytes.NewBuffer(bodyCopy))
			next.ServeHTTP(w, r)
		})
	})

	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				token, _, err := jwtauth.FromContext(r.Context())

				if err != nil || token == nil {
					json.NewEncoder(w).Encode(map[string]any{
						"success": false,
						"error":   "Unauthenticated",
					})
					return
				}
				next.ServeHTTP(w, r)
			})
		})
		r.Post("/new_meal", newMealHandler)
		r.Post("/join_meal", joinMealHandler)
		r.Post("/fetch_meal", fetchMealHandler)
		r.Post("/delete_meal", deleteMealHandler)
		r.Post("/change_user", changeUserHandler)
		r.Get("/validate_token", func(w http.ResponseWriter, r *http.Request) {
			json.NewEncoder(w).Encode(map[string]any{
				"success": true,
				"error":   "",
			})

		})
	})

	r.Post("/login", loginHandler)
	r.Post("/register", registerHandler)
	appRouter.Mount("/api", r)
}
