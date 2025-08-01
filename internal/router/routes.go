package router

import (
	"net/http"

	"github.com/calmestend/mercado_lobito/internal/auth"
	"github.com/calmestend/mercado_lobito/internal/handlers"
)

// @TODO: Add documentation
// @TODO: Add inputs validation (validator library can be used)
// @TODO: Add stylesheets
// @TODO: Add Products CRUD
// @TODO: Add Business CRUD
// @TODO: Add Students CRUD
// @TODO: Add Users CRUD
func Init() {
	// Render HTML
	http.HandleFunc("/", handlers.Home)
	http.HandleFunc("/profile", auth.AuthMiddleware(handlers.Profile))
	http.HandleFunc("/profile/config", auth.AuthMiddleware(handlers.Settings))
	http.HandleFunc("/organization", auth.AuthMiddleware(handlers.Organization))
	// Only Available if you have an organization
	http.HandleFunc("/organization/products", auth.AuthMiddleware(handlers.OrganizationProducts))
	http.HandleFunc("/organization/passport", auth.AuthMiddleware(handlers.OrganizationPassport))
	// Auth
	http.HandleFunc("/auth/login", handlers.Login)
	http.HandleFunc("/auth/register", handlers.Register)

	// API
	http.HandleFunc("/auth/signin", auth.Signin)
	http.HandleFunc("/auth/signup", auth.Signup)
	http.HandleFunc("/auth/logout", auth.Logout)
	http.ListenAndServe(":3030", nil)
}
