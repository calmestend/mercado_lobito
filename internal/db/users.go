package db

type User struct {
	id               int
	middle_names     string
	paternal_surname string
	maternal_surname string
	personal_id      int
	email            string
	hash             string
	created_at       int
	update_at        int
}
