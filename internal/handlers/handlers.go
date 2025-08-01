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

func Profile(w http.ResponseWriter, r *http.Request) {
	isAuth := auth.IsAuthenticated(r)

	profileComponent := views.Profile()
	page := views.Index(profileComponent, isAuth)
	page.Render(r.Context(), w)
}

func Organization(w http.ResponseWriter, r *http.Request) {
	isAuth := auth.IsAuthenticated(r)

	organizationComponent := views.Organization()
	page := views.Index(organizationComponent, isAuth)
	page.Render(r.Context(), w)
}

func OrganizationPassport(w http.ResponseWriter, r *http.Request) {
	isAuth := auth.IsAuthenticated(r)

	organizationPassportComponent := views.OrganizationPassport()
	page := views.Index(organizationPassportComponent, isAuth)
	page.Render(r.Context(), w)
}

func OrganizationProducts(w http.ResponseWriter, r *http.Request) {
	isAuth := auth.IsAuthenticated(r)

	productsComponent := views.Products()
	page := views.Index(productsComponent, isAuth)
	page.Render(r.Context(), w)
}

func Settings(w http.ResponseWriter, r *http.Request) {
	isAuth := auth.IsAuthenticated(r)

	settingsComponent := views.Settings()
	page := views.Index(settingsComponent, isAuth)
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
