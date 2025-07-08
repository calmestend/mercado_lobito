package db

import "database/sql"

type BusinessCollaborator struct {
	ID             int
	BusinessID     int
	CollaboratorID int
}

func (bc *BusinessCollaborator) Set(db *sql.DB) error {
	return nil
}

func (bc *BusinessCollaborator) Get(db *sql.DB) error {
	return nil
}

func (bc *BusinessCollaborator) Delete(db *sql.DB) error {
	return nil
}

func (bc *BusinessCollaborator) Update(db *sql.DB) error {
	return nil
}
