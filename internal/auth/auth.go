package auth

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/calmestend/mercado_lobito/internal/components"
	"github.com/calmestend/mercado_lobito/internal/db"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var dbConn *sql.DB

// SetDBConnection establece la conexión a la base de datos
func SetDBConnection(database *sql.DB) {
	dbConn = database
}

func generateSecureToken() string {
	return uuid.NewString()
}

// CreateSession crea una nueva sesión en la base de datos
func CreateSession(userID int) (string, error) {
	uuid := generateSecureToken()

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

// GetSession obtiene una sesión válida desde la base de datos
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

// DeleteSession elimina una sesión de la base de datos
func DeleteSession(uuid string) error {
	session := db.Session{
		UUID: uuid,
	}

	return session.Delete(dbConn)
}

// IsAuthenticated verifica si un token es válido
func IsAuthenticated(r *http.Request) bool {
	cookie, err := r.Cookie("session")
	if err != nil {
		return false
	}

	_, err = GetSession(cookie.Value)
	return err == nil
}

// GetUserByEmail obtiene un usuario por email
func GetUserByEmail(email string) (*db.User, error) {
	stmt := `
		SELECT id, middle_names, paternal_surname, maternal_surname, personal_id, email, hash
		FROM users
		WHERE email = ?
	`
	row := dbConn.QueryRow(stmt, email)

	var user db.User
	err := row.Scan(&user.ID, &user.MiddleNames, &user.PaternalSurname, &user.MaternalSurname, &user.PersonalID, &user.Email, &user.Hash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

// CreateUser crea un nuevo usuario en la base de datos
func CreateUser(email, password string) (*db.User, error) {
	// Verificar si el usuario ya existe
	_, err := GetUserByEmail(email)
	if err == nil {
		return nil, errors.New("user already exists")
	}

	// Hash de la contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := db.User{
		MiddleNames:     "", // Campo vacío
		PaternalSurname: "", // Campo vacío
		MaternalSurname: "", // Campo vacío
		PersonalID:      0,  // Campo vacío/null
		Email:           email,
		Hash:            string(hashedPassword),
	}

	err = user.Set(dbConn)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// VerifyPassword verifica si la contraseña es correcta
func VerifyPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Signin maneja el inicio de sesión
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

	email := r.FormValue("email")
	password := r.FormValue("password")

	// Validar campos vacíos
	if email == "" || password == "" {
		component := components.LoginResponse(false, "Email and password are required")
		component.Render(r.Context(), w)
		return
	}

	// Validar formato de email básico
	if !strings.Contains(email, "@") {
		component := components.LoginResponse(false, "Please enter a valid email")
		component.Render(r.Context(), w)
		return
	}

	// Buscar usuario en la base de datos
	user, err := GetUserByEmail(email)
	if err != nil {
		component := components.LoginResponse(false, "Invalid credentials")
		component.Render(r.Context(), w)
		return
	}

	// Verificar contraseña
	if !VerifyPassword(user.Hash, password) {
		component := components.LoginResponse(false, "Invalid credentials")
		component.Render(r.Context(), w)
		return
	}

	// Crear sesión
	token, err := CreateSession(user.ID)
	if err != nil {
		component := components.LoginResponse(false, "Error creating session")
		component.Render(r.Context(), w)
		return
	}

	// Establecer cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: false,
		Secure:   false, // Cambiar a true en producción con HTTPS
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	})

	// Redirigir usando HTMX
	w.Header().Set("HX-Redirect", "/dashboard")
	w.WriteHeader(http.StatusOK)
}

// Signup maneja el registro de usuarios
func Signup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		component := components.SignupResponse(false, "Error processing form")
		component.Render(r.Context(), w)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")
	confirmPassword := r.FormValue("confirm_password")

	// Validar campos vacíos
	if email == "" || password == "" || confirmPassword == "" {
		component := components.SignupResponse(false, "All fields are required")
		component.Render(r.Context(), w)
		return
	}

	// Validar que las contraseñas coincidan
	if password != confirmPassword {
		component := components.SignupResponse(false, "Passwords don't match")
		component.Render(r.Context(), w)
		return
	}

	// Crear usuario
	user, err := CreateUser(email, password)
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

	// Crear sesión automáticamente después del registro
	token, err := CreateSession(user.ID)
	if err != nil {
		component := components.SignupResponse(false, "User created but error logging in")
		component.Render(r.Context(), w)
		return
	}

	// Establecer cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: false,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	})

	// Redirigir usando HTMX
	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusOK)
}

// GetSessionFromRequest obtiene la sesión desde la request
func GetSessionFromRequest(r *http.Request) (*db.Session, error) {
	cookie, err := r.Cookie("session")
	if err != nil {
		return nil, errors.New("no session cookie found")
	}

	return GetSession(cookie.Value)
}

// GetUserFromRequest obtiene el usuario desde la request
func GetUserFromRequest(r *http.Request) (*db.User, error) {
	session, err := GetSessionFromRequest(r)
	if err != nil {
		return nil, err
	}

	user := db.User{ID: session.UserID}
	err = user.Get(dbConn)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Logout maneja el cierre de sesión
func Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err == nil {
		// Eliminar sesión de la base de datos
		DeleteSession(cookie.Value)
	}

	// Eliminar cookie del navegador
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: false,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	})

	// Redirigir usando HTMX
	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusOK)
}

// AuthMiddleware middleware para proteger rutas
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := GetSessionFromRequest(r)
		if err != nil {
			http.Redirect(w, r, "/auth/login", http.StatusFound)
			return
		}

		// Verificar que la sesión aún existe en la base de datos
		_, err = GetSession(session.UUID)
		if err != nil {
			// Limpiar cookie inválida
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
