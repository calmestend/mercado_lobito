package db

import (
	"database/sql"
	"errors"
)

type Student struct {
	ID         string
	Grade      string
	ClassGroup string
	UserID     int
}

func (s *Student) Set(db *sql.DB) error {
	stmt, err := db.Prepare(`
		INSERT INTO students(id, grade, class_group, user_id) values (?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(s.ID, s.Grade, s.ClassGroup, s.UserID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Student) GetByID(db *sql.DB) error {
	stmt := `
		SELECT id, grade, class_group, user_id FROM students WHERE id = ?
	`
	row := db.QueryRow(stmt, s.ID)
	err := row.Scan(&s.ID, &s.Grade, &s.ClassGroup, &s.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("user not found")
		}
		return err
	}

	return nil
}

func (s *Student) Update(db *sql.DB) error {
	return nil
}

func (s *Student) GetByUserID(db *sql.DB) error {
	stmt := `
		SELECT id, grade, class_group, user_id FROM students WHERE user_id = ?
	`
	row := db.QueryRow(stmt, s.UserID)
	err := row.Scan(&s.ID, &s.Grade, &s.ClassGroup, &s.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("student not found")
		}
		return err
	}
	return nil
}

func (s *Student) Delete(db *sql.DB) error {
	return nil
}
