package db

import (
	"database/sql"
	"errors"
	"github.com/google/uuid"
)

type Session struct {
	ID     int
	UUID   string
	UserID int
}

func (s *Session) Set(db *sql.DB) error {
	if s.UUID == "" {
		s.UUID = uuid.New().String()
	}

	stmt, err := db.Prepare(`
		INSERT INTO sessions(uuid, user_id)
		VALUES (?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(s.UUID, s.UserID)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err == nil {
		s.ID = int(id)
	}

	return nil
}

func (s *Session) Get(db *sql.DB) error {
	stmt := `
		SELECT id, uuid, user_id 
		FROM sessions
		WHERE uuid = ?
	`
	row := db.QueryRow(stmt, s.UUID)
	err := row.Scan(&s.ID, &s.UUID, &s.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("session not found")
		}
		return err
	}
	return nil
}

func (s *Session) Delete(db *sql.DB) error {
	stmt := `DELETE FROM sessions WHERE uuid = ?`
	_, err := db.Exec(stmt, s.UUID)
	return err
}
