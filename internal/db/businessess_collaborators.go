package db

import (
	"database/sql"
	"errors"
	"log"
)

type BusinessCollaborator struct {
	ID             int
	BusinessID     int
	CollaboratorID int
}

func (bc *BusinessCollaborator) Set(db *sql.DB) error {
	stmt, err := db.Prepare(`
		insert INTO businesses_collaborators (business_id, collaborator_id)
		VALUES (?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(bc.BusinessID, bc.CollaboratorID)
	if err != nil {
		return err
	}

	lastID, err := res.LastInsertId()
	if err == nil {
		bc.ID = int(lastID)
	}

	return nil
}

func (bc *BusinessCollaborator) Get(db *sql.DB) error {
	stmt := `
		SELECT id, business_id, collaborator_id
		FROM businesses_collaborators
		WHERE id = ?
	`
	row := db.QueryRow(stmt, bc.ID)
	err := row.Scan(&bc.ID, &bc.BusinessID, &bc.CollaboratorID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("business not found")
		}
		return err
	}
	return nil
}

func (bc *BusinessCollaborator) GetByBusinessAndCollaborator(db *sql.DB) error {
	stmt := `SELECT business_id, collaborator_id FROM business_collaborators 
			 WHERE business_id = ? AND collaborator_id = ?`

	row := db.QueryRow(stmt, bc.BusinessID, bc.CollaboratorID)
	err := row.Scan(&bc.BusinessID, &bc.CollaboratorID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("business collaborator relationship not found")
		}
		return err
	}
	return nil
}

func (b *Business) GetCollaboratorsByBusinessID(db *sql.DB) ([]User, error) {
	stmt := `
		SELECT u.id, u.middle_names, u.paternal_surname, u.maternal_surname, u.email, u.personal_id
		FROM users u
		INNER JOIN businesses_collaborators bc ON u.id = bc.collaborator_id
		WHERE bc.business_id = ?
	`

	rows, err := db.Query(stmt, b.ID)
	if err != nil {
		return nil, err
	}

	log.Print(rows)
	defer rows.Close()

	var collaborators []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.MiddleNames, &user.PaternalSurname,
			&user.MaternalSurname, &user.Email, &user.PersonalID)
		if err != nil {
			return nil, err
		}
		collaborators = append(collaborators, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return collaborators, nil
}

func (bc *BusinessCollaborator) Delete(db *sql.DB) error {
	stmt := `DELETE FROM businesses_collaborators WHERE id = ?`
	_, err := db.Exec(stmt, bc.ID)
	return err
}
