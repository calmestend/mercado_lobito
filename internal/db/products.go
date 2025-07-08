package db

import "database/sql"

type Product struct {
	ID         int
	Title      string
	Price      float64
	Stock      int
	BusinessID int
}

func (p *Product) Set(db *sql.DB) error {
	return nil
}

func (p *Product) Get(db *sql.DB) error {
	return nil
}

func (p *Product) Delete(db *sql.DB) error {
	return nil
}

func (p *Product) Update(db *sql.DB) error {
	return nil
}
