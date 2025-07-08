package db

import "database/sql"

type Business struct {
	ID      int
	Name    string
	OwnerID string
}

func (b *Business) Set(db *sql.DB) error {
	return nil
}

func (b *Business) Get(db *sql.DB) error {
	return nil
}

func (b *Business) Delete(db *sql.DB) error {
	return nil
}

func (b *Business) Update(db *sql.DB) error {
	return nil
}
