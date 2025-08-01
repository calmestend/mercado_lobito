package auth

import (
	"database/sql"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/calmestend/mercado_lobito/internal/components"
	"github.com/calmestend/mercado_lobito/internal/db"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var dbConn *sql.DB

const uploadDir = "./uploads"

func SetDBConnection(database *sql.DB) {
	dbConn = database
}

func CreateSession(userID int) (string, error) {
	uuid := uuid.NewString()

	session := db.Session{
		UUID:   uuid,
		UserID: userID,
	}

	err := session.Set(dbConn)
	if err != nil {
		return "", err
	}

	return uuid, nil
}

func GetSession(uuid string) (*db.Session, error) {
	session := db.Session{
		UUID: uuid,
	}

	err := session.Get(dbConn)
	if err != nil {
		return nil, errors.New("session not found")
	}

	return &session, nil
}

func DeleteSession(uuid string) error {
	session := db.Session{
		UUID: uuid,
	}

	return session.Delete(dbConn)
}

func IsAuthenticated(r *http.Request) bool {
	cookie, err := r.Cookie("session")
	if err != nil {
		return false
	}

	_, err = GetSession(cookie.Value)
	return err == nil
}

func VerifyPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Signin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		component := components.LoginResponse(false, "Error processing form")
		component.Render(r.Context(), w)
		return
	}

	studentID := r.FormValue("student_id")
	password := r.FormValue("password")

	if studentID == "" || password == "" {
		component := components.LoginResponse(false, "Student ID and password are required")
		component.Render(r.Context(), w)
		return
	}

	s := db.Student{
		ID: studentID,
	}

	err = s.GetByID(dbConn)
	if err != nil {
		component := components.LoginResponse(false, "Invalid Credentials")
		component.Render(r.Context(), w)
		return
	}

	u := db.User{
		ID: s.UserID,
	}

	err = u.GetByID(dbConn)
	if err != nil {
		component := components.LoginResponse(false, "Invalid Credentials")
		component.Render(r.Context(), w)
		return
	}

	if !VerifyPassword(u.Hash, password) {
		component := components.LoginResponse(false, "Invalid Credentials")
		component.Render(r.Context(), w)
		return
	}

	token, err := CreateSession(u.ID)
	if err != nil {
		component := components.LoginResponse(false, "Error creating session")
		component.Render(r.Context(), w)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: false,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	})

	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusOK)
}

func Signup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(5 << 20) // 5 MB
	if err != nil {
		component := components.SignupResponse(false, "Error processing form")
		component.Render(r.Context(), w)
		return
	}

	defer r.MultipartForm.RemoveAll()
	defer r.Body.Close()

	u := db.User{
		MiddleNames:     r.FormValue("middle_names"),
		PaternalSurname: r.FormValue("paternal_surname"),
		MaternalSurname: r.FormValue("maternal_surname"),
		Email:           r.FormValue("email"),
	}

	s := db.Student{
		ID:         r.FormValue("student_id"),
		Grade:      "",
		ClassGroup: "",
	}

	password := r.FormValue("password")
	confirmPassword := r.FormValue("confirm_password")

	// --- File Upload Logic ---
	file, _, err := r.FormFile("file") // "file" is the name attribute from your <input type="file">
	if err != nil {
		if err == http.ErrMissingFile {
			component := components.SignupResponse(false, "Profile photo is required")
			component.Render(r.Context(), w)
			return
		} else {
			component := components.SignupResponse(false, "Error retrieving profile photo: "+err.Error())
			component.Render(r.Context(), w)
			return
		}
	} else {
		defer file.Close()

		// Ensure the upload directory exists
		if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
			err = os.MkdirAll(uploadDir, 0755) // Create with read/write/execute permissions for owner, read/execute for group/others
			if err != nil {
				component := components.SignupResponse(false, "Error creating upload directory: "+err.Error())
				component.Render(r.Context(), w)
				return
			}
		}

		// Use the student ID for the filename
		filename := s.ID + ".jpg"

		// Construct the full path to save the file
		filePath := filepath.Join(uploadDir, filename)

		// Create the new file on the server
		dst, err := os.Create(filePath)
		if err != nil {
			component := components.SignupResponse(false, "Error creating file on server: "+err.Error())
			component.Render(r.Context(), w)
			return
		}
		defer dst.Close() // Close the destination file

		// Copy the uploaded file's content to the new file
		if _, err := io.Copy(dst, file); err != nil {
			component := components.SignupResponse(false, "Error saving profile photo: "+err.Error())
			component.Render(r.Context(), w)
			return
		}
	}

	// Empty values validation
	if s.ID == "" || u.MiddleNames == "" || u.PaternalSurname == "" || u.MaternalSurname == "" || u.Email == "" || password == "" || confirmPassword == "" {
		component := components.SignupResponse(false, "All fields are required")
		component.Render(r.Context(), w)
		return
	}

	if password != confirmPassword {
		component := components.SignupResponse(false, "Passwords don't match")
		component.Render(r.Context(), w)
		return
	}

	// Create User
	err = u.GetByEmail(dbConn)
	if err != nil && !strings.Contains(err.Error(), "user not found") {
		component := components.SignupResponse(false, "Error creating user")
		component.Render(r.Context(), w)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		component := components.SignupResponse(false, "Error creating user")
		component.Render(r.Context(), w)
		return
	}

	u.Hash = string(hashedPassword)

	err = u.Set(dbConn)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			component := components.SignupResponse(false, "User already exists")
			component.Render(r.Context(), w)
		} else {
			component := components.SignupResponse(false, "Error creating user")
			component.Render(r.Context(), w)
		}
		return
	}

	s.UserID = u.ID
	err = s.Set(dbConn)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			component := components.SignupResponse(false, "Student already exists")
			component.Render(r.Context(), w)
		} else {
			component := components.SignupResponse(false, "Error creating student")
			component.Render(r.Context(), w)
		}
		return
	}

	// Create session
	token, err := CreateSession(u.ID)
	if err != nil {
		component := components.SignupResponse(false, "User created but error logging in")
		component.Render(r.Context(), w)
		return
	}

	// Set cookies
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: false,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	})

	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusOK)
}

func GetSessionFromRequest(r *http.Request) (*db.Session, error) {
	cookie, err := r.Cookie("session")
	if err != nil {
		return nil, errors.New("no session cookie found")
	}

	return GetSession(cookie.Value)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err == nil {
		DeleteSession(cookie.Value)
	}

	// Delete Cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: false,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	})

	http.Redirect(w, r, "/auth/login", http.StatusFound)
}

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := GetSessionFromRequest(r)
		if err != nil {
			http.Redirect(w, r, "/auth/login", http.StatusFound)
			return
		}

		_, err = GetSession(session.UUID)
		if err != nil {
			http.SetCookie(w, &http.Cookie{
				Name:     "session",
				Value:    "",
				Expires:  time.Now().Add(-time.Hour),
				HttpOnly: false,
				Secure:   false,
				SameSite: http.SameSiteStrictMode,
				Path:     "/",
			})
			http.Redirect(w, r, "/auth/login", http.StatusFound)
			return
		}

		next(w, r)
	}
}
