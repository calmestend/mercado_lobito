package router

import (
	"net/http"

	"github.com/calmestend/mercado_lobito/internal/auth"
	"github.com/calmestend/mercado_lobito/internal/handlers"
)

func Init() {
	// Render HTML
	http.HandleFunc("/", handlers.Home)
	http.HandleFunc("/dashboard", auth.AuthMiddleware(handlers.Home))
	http.HandleFunc("/profile", auth.AuthMiddleware(handlers.Home))

	http.HandleFunc("/auth/login", handlers.Login)
	http.HandleFunc("/auth/register", handlers.Register)

	// API
	http.HandleFunc("/auth/signin", auth.Signin)
	http.HandleFunc("/auth/signup", auth.Signup)
	http.HandleFunc("/auth/logout", auth.Logout)
	http.ListenAndServe(":3030", nil)
}
