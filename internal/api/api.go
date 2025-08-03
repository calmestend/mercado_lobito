package api

import (
	"fmt"
	"net/http"

	"github.com/calmestend/mercado_lobito/internal/auth"
	"github.com/calmestend/mercado_lobito/internal/db"
)

func ProfileConfig(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	dbConn := db.Init()

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

	name := r.FormValue("business_name")
	businessType := r.FormValue("business_type")
	description := r.FormValue("description")

	business := db.Business{OwnerID: student.ID}
	if err := business.GetByOwnerID(dbConn); err == nil {
		business.Name = name
		business.Type = businessType
		business.Description = description
		if err := business.Update(dbConn); err != nil {
			http.Error(w, "Error updating business", http.StatusInternalServerError)
			return
		}
	} else {
		business = db.Business{
			Name:        name,
			Type:        businessType,
			Description: description,
			OwnerID:     student.ID,
		}
		if err := business.Set(dbConn); err != nil {
			http.Error(w, "Error creating business", http.StatusInternalServerError)
			return
		}
	}

	fmt.Fprintf(w, `<p style="color:green;">Datos guardados correctamente.</p>`)
}

func Organization(w http.ResponseWriter, r *http.Request) {

}

