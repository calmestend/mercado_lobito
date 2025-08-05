package api

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/calmestend/mercado_lobito/internal/auth"
	"github.com/calmestend/mercado_lobito/internal/db"
	"github.com/calmestend/mercado_lobito/internal/views"
)

func Products(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/api/products/edit/") {
		editProduct(w, r)
		return
	}
	if strings.HasPrefix(r.URL.Path, "/api/products/cancel/") {
		cancelEditProduct(w, r)
		return
	}

	switch r.Method {
	case http.MethodGet:
		getAllProductsByBusiness(w, r)
	case http.MethodPost:
		createProduct(w, r)
	case http.MethodPatch:
		updateProduct(w, r)
	case http.MethodDelete:
		deleteProduct(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func getAllProductsByBusiness(w http.ResponseWriter, r *http.Request) {
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

	products, err := business.GetProductsByOwnerID(dbConn)
	if err != nil {
		http.Error(w, "Error retrieving products", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	component := views.ProductsTable(products)
	component.Render(r.Context(), w)
}

func createProduct(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
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

	title := r.FormValue("title")
	price := r.FormValue("price")
	stock := r.FormValue("stock")

	priceFloat, err := strconv.ParseFloat(price, 64)
	if err != nil {
		http.Error(w, "Error parsing price", http.StatusBadRequest)
		return
	}

	stockInt, err := strconv.Atoi(stock)
	if err != nil {
		http.Error(w, "Error parsing stock", http.StatusBadRequest)
		return
	}

	product := db.Product{
		Title:      title,
		Price:      priceFloat,
		Stock:      stockInt,
		BusinessID: business.ID,
	}

	err = product.Set(dbConn)
	if err != nil {
		http.Error(w, "Error creating product", http.StatusInternalServerError)
		return
	}

	products, err := business.GetProductsByOwnerID(dbConn)
	if err != nil {
		http.Error(w, "Error retrieving products", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	component := views.ProductsTable(products)
	component.Render(r.Context(), w)
}

func editProduct(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 5 {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	idStr := pathParts[4]
	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	dbConn := db.Init()
	defer dbConn.Close()

	product := db.Product{ID: idInt}
	err = product.GetByID(dbConn)
	if err != nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	component := views.ProductEditRow(product)
	component.Render(r.Context(), w)
}

func cancelEditProduct(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 5 {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	idStr := pathParts[4]
	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	dbConn := db.Init()
	defer dbConn.Close()

	product := db.Product{ID: idInt}
	err = product.GetByID(dbConn)
	if err != nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	component := views.ProductRow(product)
	component.Render(r.Context(), w)
}

func updateProduct(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	dbConn := db.Init()
	defer dbConn.Close()

	// Get session to verify ownership
	sess, err := auth.GetSessionFromRequest(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	id := r.FormValue("id")
	title := r.FormValue("title")
	price := r.FormValue("price")
	stock := r.FormValue("stock")

	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Error parsing product id", http.StatusBadRequest)
		return
	}

	priceFloat, err := strconv.ParseFloat(price, 64)
	if err != nil {
		http.Error(w, "Error parsing price", http.StatusBadRequest)
		return
	}

	stockInt, err := strconv.Atoi(stock)
	if err != nil {
		http.Error(w, "Error parsing stock", http.StatusBadRequest)
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

	existingProduct := db.Product{ID: idInt}
	if err := existingProduct.GetByID(dbConn); err != nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	if existingProduct.BusinessID != business.ID {
		http.Error(w, "Unauthorized to edit this product", http.StatusForbidden)
		return
	}

	product := db.Product{
		ID:         idInt,
		Title:      title,
		Price:      priceFloat,
		Stock:      stockInt,
		BusinessID: business.ID,
	}

	err = product.Update(dbConn)
	if err != nil {
		http.Error(w, "Error updating product", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	component := views.ProductRow(product)
	component.Render(r.Context(), w)
}

func deleteProduct(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	dbConn := db.Init()
	defer dbConn.Close()

	id := r.FormValue("id")

	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Error parsing product id", http.StatusBadRequest)
		return
	}

	product := db.Product{ID: idInt}
	err = product.Delete(dbConn)
	if err != nil {
		http.Error(w, "Error deleting product", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	component := views.ProductDeleted()
	component.Render(r.Context(), w)
}
