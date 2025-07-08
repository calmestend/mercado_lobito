package db

import "database/sql"

type Student struct {
	ID         string
	Grade      string
	ClassGroup string
	UserID     int
}

func (s *Student) Set(db *sql.DB) error {
	return nil
}
func (s *Student) Get(db *sql.DB) error {
	return nil
}
func (s *Student) Delete(db *sql.DB) error {
	return nil
}
func (s *Student) Update(db *sql.DB) error {
	return nil
}
