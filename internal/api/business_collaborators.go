package api

import (
	"log"
	"net/http"
	"strconv"

	"github.com/calmestend/mercado_lobito/internal/auth"
	"github.com/calmestend/mercado_lobito/internal/db"
	"github.com/calmestend/mercado_lobito/internal/views"
)

func BusinessCollaborators(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getAllCollaboratorsByBusiness(w, r)
	case http.MethodPost:
		createCollaborator(w, r)
	case http.MethodPatch:
		updateCollaborator(w, r)
	case http.MethodDelete:
		deleteCollaborator(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func getAllCollaboratorsByBusiness(w http.ResponseWriter, r *http.Request) {
	dbConn := db.Init()
	defer dbConn.Close()
	sess, err := auth.GetSessionFromRequest(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	student := db.Student{UserID: sess.UserID}
	if err := student.GetByUserID(dbConn); err != nil {
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	}
	business := db.Business{OwnerID: student.ID}
	if err := business.GetByOwnerID(dbConn); err != nil {
		http.Error(w, "Business not found", http.StatusNotFound)
		return
	}

	collabs, err := business.GetCollaboratorsByBusinessID(dbConn)
	log.Print(collabs)
	if err != nil {
		http.Error(w, "Error retrieving collaborators", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	views.CollaboratorsList(collabs).Render(r.Context(), w)
}

func createCollaborator(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}
	dbConn := db.Init()
	defer dbConn.Close()
	sess, err := auth.GetSessionFromRequest(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	student := db.Student{UserID: sess.UserID}
	if err := student.GetByUserID(dbConn); err != nil {
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	}
	business := db.Business{OwnerID: student.ID}
	if err := business.GetByOwnerID(dbConn); err != nil {
		http.Error(w, "Business not found", http.StatusNotFound)
		return
	}

	file, _, err := r.FormFile("file")
	if err == nil {
		_ = file.Close()
	}

	isIntern := r.FormValue("isIntern")
	middleNames := r.FormValue("middle_names")
	paternal := r.FormValue("paternal_surname")
	maternal := r.FormValue("maternal_surname")
	email := r.FormValue("email")

	u := db.User{
		MiddleNames:     middleNames,
		PaternalSurname: paternal,
		MaternalSurname: maternal,
		Email:           email,
	}

	if isIntern == "false" {
		u.PersonalID = r.FormValue("personal_id")
	}

	if err := u.Set(dbConn); err != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	if isIntern == "true" {
		s := db.Student{Grade: r.FormValue("grade"), ClassGroup: r.FormValue("class_group"), UserID: u.ID, ID: r.FormValue("student_id")}
		if err := s.Set(dbConn); err != nil {
			http.Error(w, "Error creating student", http.StatusInternalServerError)
			return
		}
	}

	bc := db.BusinessCollaborator{BusinessID: business.ID, CollaboratorID: u.ID}
	if err := bc.Set(dbConn); err != nil {
		http.Error(w, "Error linking collaborator", http.StatusInternalServerError)
		return
	}

	getAllCollaboratorsByBusiness(w, r)
}

func updateCollaborator(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad form", http.StatusBadRequest)
		return
	}
	dbConn := db.Init()
	defer dbConn.Close()
	sess, err := auth.GetSessionFromRequest(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	student := db.Student{UserID: sess.UserID}
	if err := student.GetByUserID(dbConn); err != nil {
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	}
	business := db.Business{OwnerID: student.ID}
	if err := business.GetByOwnerID(dbConn); err != nil {
		http.Error(w, "Business not found", http.StatusNotFound)
		return
	}

	idVal := r.FormValue("collaborator_id")
	idInt, err := strconv.Atoi(idVal)
	if err != nil {
		http.Error(w, "Bad collaborator ID", http.StatusBadRequest)
		return
	}

	bc := db.BusinessCollaborator{BusinessID: business.ID, CollaboratorID: idInt}
	if err := bc.GetByBusinessAndCollaborator(dbConn); err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	user := db.User{ID: idInt,
		MiddleNames:     r.FormValue("middle_names"),
		PaternalSurname: r.FormValue("paternal_surname"),
		MaternalSurname: r.FormValue("maternal_surname"),
		Email:           r.FormValue("email"),
	}
	if err := user.Update(dbConn); err != nil {
		http.Error(w, "Error updating user", http.StatusInternalServerError)
		return
	}

	if r.FormValue("isIntern") == "true" {
		s := db.Student{UserID: idInt}
		if err := s.GetByUserID(dbConn); err == nil {
			s.Grade = r.FormValue("grade")
			s.ClassGroup = r.FormValue("class_group")
			if err := s.Update(dbConn); err != nil {
				http.Error(w, "Error updating student", http.StatusInternalServerError)
				return
			}
		}
	} else {
		user.PersonalID = r.FormValue("personal_id")
		if err := user.Update(dbConn); err != nil {
			http.Error(w, "Error updating external data", http.StatusInternalServerError)
			return
		}
	}

	getAllCollaboratorsByBusiness(w, r)
}

func deleteCollaborator(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad form", http.StatusBadRequest)
		return
	}
	dbConn := db.Init()
	defer dbConn.Close()
	sess, err := auth.GetSessionFromRequest(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	student := db.Student{UserID: sess.UserID}
	if err := student.GetByUserID(dbConn); err != nil {
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	}
	business := db.Business{OwnerID: student.ID}
	if err := business.GetByOwnerID(dbConn); err != nil {
		http.Error(w, "Business not found", http.StatusNotFound)
		return
	}

	idVal := r.FormValue("collaborator_id")
	idInt, err := strconv.Atoi(idVal)
	if err != nil {
		http.Error(w, "Bad id", http.StatusBadRequest)
		return
	}

	bc := db.BusinessCollaborator{BusinessID: business.ID, CollaboratorID: idInt}
	if err := bc.GetByBusinessAndCollaborator(dbConn); err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if err := bc.Delete(dbConn); err != nil {
		http.Error(w, "Error deleting collaborator", http.StatusInternalServerError)
		return
	}

	getAllCollaboratorsByBusiness(w, r)
}

func CollaboratorForm(w http.ResponseWriter, r *http.Request) {
	typ := r.URL.Query().Get("type")
	w.Header().Set("Content-Type", "text/html")
	switch typ {
	case "intern":
		views.InternFields().Render(r.Context(), w)
	case "external":
		views.ExternalFields().Render(r.Context(), w)
	default:
		http.Error(w, "Invalid Type", http.StatusBadRequest)
	}
}
