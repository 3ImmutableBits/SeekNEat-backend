package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth"
	"gorm.io/gorm"
	"github.com/3ImmutableBits/SeekNEat-backend/config"
)

var tokenAuth *jwtauth.JWTAuth

var db *gorm.DB

func AddRoutes(appRouter chi.Router, dbLocal *gorm.DB) {
	db = dbLocal

	r := chi.NewRouter()

	tokenAuth = jwtauth.New("HS256", []byte(config.SecretKey), nil)

	r.Post("/login", loginHandler)
	r.Post("/register", registerHandler)

	appRouter.Mount("/api", r)
}
