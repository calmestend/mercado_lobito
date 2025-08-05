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

	// Expose img directory
	fs := http.FileServer(http.Dir("./internal/img"))
	http.Handle("/img/", http.StripPrefix("/img/", fs))

	// Expose Upload directory
	http.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads"))))

	// API
	http.HandleFunc("/auth/signin", auth.Signin)
	http.HandleFunc("/auth/signup", auth.Signup)
	http.HandleFunc("/auth/logout", auth.Logout)
	http.HandleFunc("/api/profile/config", auth.AuthMiddleware(api.ProfileConfig))
	http.HandleFunc("/api/business/collaborators", auth.AuthMiddleware(api.BusinessCollaborators))
	http.HandleFunc("/api/products/edit/", auth.AuthMiddleware(api.Products))
	http.HandleFunc("/api/products/cancel/", auth.AuthMiddleware(api.Products))
	http.HandleFunc("/api/products", auth.AuthMiddleware(api.Products))

	http.ListenAndServe(":3030", nil)
}
