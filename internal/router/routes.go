package router

import (
	"net/http"

	"github.com/calmestend/mercado_lobito/internal/api"
	"github.com/calmestend/mercado_lobito/internal/auth"
	"github.com/calmestend/mercado_lobito/internal/handlers"
)

func Init() {
	// Render HTML
	http.HandleFunc("/", handlers.Home)
	http.HandleFunc("/profile", auth.AuthMiddleware(handlers.Profile))
	http.HandleFunc("/profile/config", auth.AuthMiddleware(handlers.Settings))
	http.HandleFunc("/organization", auth.AuthMiddleware(handlers.Organization))

	// Only Available if you have an organization
	http.HandleFunc("/organization/products", auth.AuthMiddleware(handlers.OrganizationProducts))
	http.HandleFunc("/organization/passport", auth.AuthMiddleware(handlers.OrganizationPassport)) // @TODO: Add download pdfs support via multipart-form

	// Auth
	http.HandleFunc("/auth/login", handlers.Login)
	http.HandleFunc("/auth/register", handlers.Register)

	// Expose Upload directory
	http.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads"))))

	// API
	http.HandleFunc("/auth/signin", auth.Signin)
	http.HandleFunc("/auth/signup", auth.Signup)
	http.HandleFunc("/auth/logout", auth.Logout)
	http.HandleFunc("/api/profile/config", api.ProfileConfig)
	http.HandleFunc("/api/business/collaborators", api.BusinessCollaborators)
	http.HandleFunc("/api/products", api.Products)

	http.ListenAndServe(":3030", nil)
}
