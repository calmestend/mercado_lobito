package handlers

import (
	"fmt"
	"net/http"

	"github.com/calmestend/mercado_lobito/internal/auth"
	"github.com/calmestend/mercado_lobito/internal/components"
	"github.com/calmestend/mercado_lobito/internal/db"
	"github.com/calmestend/mercado_lobito/internal/views"
)

func Home(w http.ResponseWriter, r *http.Request) {
	isAuth := auth.IsAuthenticated(r)

	homeComponent := views.Home()
	page := views.Index(homeComponent, isAuth)
	page.Render(r.Context(), w)
}

func Profile(w http.ResponseWriter, r *http.Request) {
	dbConn := db.Init()
	defer dbConn.Close()

	isAuth := auth.IsAuthenticated(r)

	sess, err := auth.GetSessionFromRequest(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	user := db.User{ID: sess.UserID}
	if err := user.GetByID(dbConn); err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	student := db.Student{UserID: user.ID}
	if err := student.GetByUserID(dbConn); err != nil {
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	}

	fullName := fmt.Sprintf("%s %s %s", user.MiddleNames, user.PaternalSurname, user.MaternalSurname)
	imgSrc := fmt.Sprintf("/uploads/%s.jpg", student.ID)

	profileComponent := views.Profile(fullName, imgSrc)
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
	dbConn := db.Init()
	isAuth := auth.IsAuthenticated(r)

	sess, err := auth.GetSessionFromRequest(r)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	student := db.Student{UserID: sess.UserID}
	if err := student.GetByUserID(dbConn); err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	business := db.Business{OwnerID: student.ID}
	err = business.GetByOwnerID(dbConn)

	settingsComponent := views.Settings(
		fmt.Sprintf("/uploads/%s.jpg", student.ID),
		business.Name,
		business.Type,
		business.Description,
	)
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
