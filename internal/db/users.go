package db

import (
	"database/sql"
	"errors"
)

type User struct {
	ID              int
	MiddleNames     string
	PaternalSurname string
	MaternalSurname string
	PersonalID      int
	Email           string
	Hash            string
}

func (u *User) Set(db *sql.DB) error {
	stmt, err := db.Prepare(`
		INSERT INTO users(middle_names, paternal_surname, maternal_surname, personal_id, email, hash)
		VALUES (?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(u.MiddleNames, u.PaternalSurname, u.MaternalSurname, u.PersonalID, u.Email, u.Hash)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err == nil {
		u.ID = int(id)
	}

	return nil
}

func (u *User) GetByID(db *sql.DB) error {
	stmt := `
		SELECT id, middle_names, paternal_surname, maternal_surname, personal_id, email, hash
		FROM users
		WHERE id = ?
	`
	row := db.QueryRow(stmt, u.ID)
	err := row.Scan(&u.ID, &u.MiddleNames, &u.PaternalSurname, &u.MaternalSurname, &u.PersonalID, &u.Email, &u.Hash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("user not found")
		}
		return err
	}

	return nil
}

func (u *User) GetByEmail(db *sql.DB) error {
	stmt := `
		SELECT id, middle_names, paternal_surname, maternal_surname, personal_id, email, hash
		FROM users
		WHERE email = ?
	`

	row := db.QueryRow(stmt, u.Email)
	err := row.Scan(&u.ID, &u.MiddleNames, &u.PaternalSurname, &u.MaternalSurname, &u.PersonalID, &u.Email, &u.Hash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("user not found")
		}
		return err
	}
	return nil
}

func (u *User) Update(db *sql.DB) error {
	stmt := `
		UPDATE users
		SET middle_names = ?, paternal_surname = ?, maternal_surname = ?, personal_id = ?, email = ?, hash = ?
		WHERE id = ?
	`
	_, err := db.Exec(stmt, u.MiddleNames, u.PaternalSurname, u.MaternalSurname, u.PersonalID, u.Email, u.Hash, u.ID)
	return err
}

func (u *User) Delete(db *sql.DB) error {
	stmt := `DELETE FROM users WHERE id = ?`
	_, err := db.Exec(stmt, u.ID)
	return err
}
