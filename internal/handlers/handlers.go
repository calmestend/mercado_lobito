package handlers

import (
	"net/http"

	"github.com/calmestend/mercado_lobito/internal/auth"
	"github.com/calmestend/mercado_lobito/internal/components"
	"github.com/calmestend/mercado_lobito/internal/views"
)

func Home(w http.ResponseWriter, r *http.Request) {
	isAuth := auth.IsAuthenticated(r)

	homeComponent := views.Home()
	page := views.Index(homeComponent, isAuth)
	page.Render(r.Context(), w)
}

func Login(w http.ResponseWriter, r *http.Request) {
	if auth.IsAuthenticated(r) {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	loginComponent := components.Login()
	page := views.Index(loginComponent, false)
	page.Render(r.Context(), w)
}

func Register(w http.ResponseWriter, r *http.Request) {
	if auth.IsAuthenticated(r) {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	signupComponent := components.Signup()
	page := views.Index(signupComponent, false)
	page.Render(r.Context(), w)
}
