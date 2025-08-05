package db

import (
	"database/sql"
	"errors"
)

type Product struct {
	ID         int
	Title      string
	Price      float64
	Stock      int
	BusinessID int
}

func (p *Product) Set(db *sql.DB) error {
	stmt, err := db.Prepare(`
		INSERT INTO products(title, price, stock, business_id)
		values (?, ?, ?, ?)
		`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(p.Title, p.Price, p.Stock, p.BusinessID)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err == nil {
		p.ID = int(id)
	}

	return nil
}

func (p *Product) GetByID(db *sql.DB) error {
	stmt := `
		SELECT id, title, price, stock, business_id
		FROM products
		WHERE id = ?
	`

	row := db.QueryRow(stmt, p.ID)
	err := row.Scan(&p.ID, &p.Title, &p.Price, &p.Stock, &p.BusinessID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("product not found")
		}
		return err
	}

	return nil
}

func (p *Product) Update(db *sql.DB) error {
	stmt := `
		UPDATE products
		SET title = ?, price = ?, stock = ?
		WHERE id = ?
	`

	_, err := db.Exec(stmt, p.Title, p.Price, p.Stock, p.ID)

	return err
}

func (p *Product) Delete(db *sql.DB) error {
	stmt := `DELETE FROM products WHERE id = ?`
	_, err := db.Exec(stmt, p.ID)
	return err
}
