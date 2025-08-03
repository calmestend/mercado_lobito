package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/calmestend/mercado_lobito/pkg/env"
	_ "github.com/go-sql-driver/mysql"
)

type database struct {
	User     string
	Password string
	Name     string
	Host     string
}

// Create and connect to mysql db
func Init() *sql.DB {
	var dbVariables *database = getVariables()

	connString := fmt.Sprintf("%s:%s@tcp(%s)/%s", dbVariables.User, dbVariables.Password, dbVariables.Host, dbVariables.Name)

	db, err := sql.Open("mysql", connString)
	if err != nil {
		log.Fatalf("Error %s when opening DB\n", err)
	}

	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()

	_, err = db.ExecContext(ctx, fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", dbVariables.Name))
	if err != nil {
		log.Fatalf("Error %s when creating DB\n", err)
	}

	err = db.PingContext(ctx)
	if err != nil {
		log.Fatalf("Error %s when pinging DB\n", err)
	}

	log.Print("Connected to DB successfully")

	createSchema(db)

	return db
}

// TODO: Add unique and validations
func createSchema(db *sql.DB) error {
	tables := []string{
		`
		CREATE TABLE IF NOT EXISTS users (
			id INT AUTO_INCREMENT PRIMARY KEY,
			middle_names VARCHAR(100) NOT NULL,
			paternal_surname VARCHAR(100) NOT NULL,
			maternal_surname VARCHAR(100) DEFAULT NULL,
			personal_id CHAR(100) DEFAULT NULL,
			email VARCHAR(150) NOT NULL,
			hash CHAR(255) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			update_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
		);
		`,
		`
		CREATE TABLE IF NOT EXISTS admins (
			id INT AUTO_INCREMENT PRIMARY KEY,
			user_id INT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			update_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		);
		`,
		`
		CREATE TABLE IF NOT EXISTS students (
			id CHAR(10) PRIMARY KEY,
			grade VARCHAR(20) NOT NULL,
			class_group VARCHAR(10) NOT NULL,
			user_id INT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			update_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		);
		`,
		`
		CREATE TABLE IF NOT EXISTS businesses (
			id INT AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			type VARCHAR(100) DEFAULT NULL,
			description TEXT DEFAULT NULL,
			owner_id CHAR(10),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			update_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			FOREIGN KEY (owner_id) REFERENCES students(id) ON DELETE CASCADE
		);
		`,
		`
		CREATE TABLE IF NOT EXISTS businesses_collaborators (
			id INT AUTO_INCREMENT PRIMARY KEY,
			business_id INT,
			collaborator_id INT,
			FOREIGN KEY (business_id) REFERENCES businesses(id) ON DELETE CASCADE,
			FOREIGN KEY (collaborator_id) REFERENCES users(id) ON DELETE CASCADE
		);
		`,
		`
		CREATE TABLE IF NOT EXISTS products (
			id INT AUTO_INCREMENT PRIMARY KEY,
			title VARCHAR(100) NOT NULL,
			price DECIMAL(10,2) NOT NULL,
			stock INT NOT NULL DEFAULT 0,
			business_id INT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			update_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			FOREIGN KEY (business_id) REFERENCES businesses(id) ON DELETE CASCADE
		);
		`,
		`
		CREATE TABLE IF NOT EXISTS sessions (
			id INT AUTO_INCREMENT PRIMARY KEY,
			uuid CHAR(255) NOT NULL,
			user_id INT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			update_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		);
		`,
	}

	for _, table := range tables {
		stmt, err := db.Prepare(table)
		if err != nil {
			return err
		}

		stmt.Exec()
	}

	return nil
}

func getVariables() *database {
	user, err := env.GetEnv("MYSQL_USER")
	if err != nil {
		log.Fatal(err)
	}

	pass, err := env.GetEnv("MYSQL_PASSWORD")
	if err != nil {
		log.Fatal(err)
	}

	dbName, err := env.GetEnv("MYSQL_DATABASE_NAME")
	if err != nil {
		log.Fatal(err)
	}

	host, err := env.GetEnv("MYSQL_HOST")
	if err != nil {
		log.Fatal(err)
	}

	return &database{
		User:     user,
		Password: pass,
		Name:     dbName,
		Host:     host,
	}
}
