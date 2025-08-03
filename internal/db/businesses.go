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

func (b *Business) Update(db *sql.DB) error {
	stmt := `
		UPDATE businesses
		SET name = ?, type = ?, description = ?
		WHERE id = ?
	`
	_, err := db.Exec(stmt, b.Name, b.Type, b.Description, b.ID)
	return err
}

func (b *Business) Delete(db *sql.DB) error {
	stmt := `DELETE FROM businesses WHERE id = ?`
	_, err := db.Exec(stmt, b.ID)
	return err
}
