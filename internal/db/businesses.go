package db

import (
	"database/sql"
	"errors"
)

type Business struct {
	ID          int
	Name        string
	Type        string
	Description string
	OwnerID     string
}

func (b *Business) Set(db *sql.DB) error {
	stmt, err := db.Prepare(`
		insert INTO businesses (name, type, description, owner_id)
		VALUES (?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(b.Name, b.Type, b.Description, b.OwnerID)
	if err != nil {
		return err
	}

	lastID, err := res.LastInsertId()
	if err == nil {
		b.ID = int(lastID)
	}

	return nil
}

func (b *Business) Get(db *sql.DB) error {
	stmt := `
		SELECT id, name, type, description, owner_id
		FROM businesses
		WHERE id = ?
	`
	row := db.QueryRow(stmt, b.ID)
	err := row.Scan(&b.ID, &b.Name, &b.Type, &b.Description, &b.OwnerID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("business not found")
		}
		return err
	}
	return nil
}

func (b *Business) GetByOwnerID(db *sql.DB) error {
	stmt := `
		SELECT id, name, type, description, owner_id
		FROM businesses
		WHERE owner_id = ?
	`
	row := db.QueryRow(stmt, b.OwnerID)
	err := row.Scan(&b.ID, &b.Name, &b.Type, &b.Description, &b.OwnerID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("business not found")
		}
		return err
	}
	return nil
}

func (b *Business) GetProductsByOwnerID(db *sql.DB) ([]Product, error) {
	stmt := `
		SELECT p.id, p.title, p.price, p.stock, p.business_id
		FROM products p
		INNER JOIN businesses b ON p.business_id = b.id
		WHERE b.owner_id = ?
	`

	rows, err := db.Query(stmt, b.OwnerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var product Product
		err := rows.Scan(&product.ID, &product.Title, &product.Price, &product.Stock, &product.BusinessID)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func (b *Business) Delete(db *sql.DB) error {
	stmt := `DELETE FROM businesses WHERE id = ?`
	_, err := db.Exec(stmt, b.ID)
	return err
}

func (b *Business) Update(db *sql.DB) error {
	stmt := `
		UPDATE businesses
		SET name = ?, type = ?, description = ?
		WHERE id = ?
	`
	_, err := db.Exec(stmt, b.Name, b.Type, b.Description, b.ID)
	return err
}
