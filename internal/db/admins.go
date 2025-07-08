package db

import "database/sql"

type Admin struct {
	Id     int
	UserID int
}

func (a *Admin) Set(db *sql.DB) error {
	return nil
}

func (a *Admin) Get(db *sql.DB) error {
	return nil
}

func (a *Admin) Delete(db *sql.DB) error {
	return nil
}

func (a *Admin) Update(db *sql.DB) error {
	return nil
}
