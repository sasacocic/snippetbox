package models

import (
	"database/sql"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmnt := `INSERT INTO users (name, email, hashed_password, created)
    VALUES ($1, $2, $3, NOW())`

	_, err = m.DB.Exec(stmnt, name, email, string(hashedPassword))
	if err != nil {

		// need to do some error handling here regarding the type of error
		// but I'm not particularaly interested in that right now - since I decided
		// to use postgresql instead of MySQL it's going to be a bit different
		panic("not handling this error  " + err.Error())
	}

	return nil
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	var id int
	var hashedPassword []byte

	stmt := "SELECT id, hashed_password FROM users WHERE email = $1"

	err := m.DB.QueryRow(stmt, email).Scan(&id, &hashedPassword)
	if err != nil {
		return 0, err // TODO: i'm not doing the right thing here
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		return 0, nil // TODO: not doing proper error handling the way the book does here
	}

	return id, nil
}

func (m *UserModel) Exists(id int) (bool, error) {

	var exists bool

	stmt := "SELECT EXISTS(SELECT true FROM users WHERE id = $1)"

	err := m.DB.QueryRow(stmt, id).Scan(&exists)

	return exists, err
}
