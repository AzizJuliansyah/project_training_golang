package models

import (
	"database/sql"
	"financial_record/entities"
)

type AuthModel struct {
	db *sql.DB
}

func NewAuthModel(db *sql.DB) *AuthModel {
	return &AuthModel{
		db: db,
	}
}

func (model AuthModel) FindUserByEmail(email string) (entities.Auth, error) {
	var user entities.Auth
	query := `
		SELECT id, email, name, password FROM users WHERE email = ?
	`

	err := model.db.QueryRow(query, email).Scan(
		&user.Id,
		&user.Email,
		&user.Name,
		&user.Password,
	)

	if err != nil {
		return user, err
	}

	return user, nil
}

func (model AuthModel) Register(user entities.Register) error {
	_, err := model.db.Exec(
		"INSERT INTO users (name, email, password) VALUES (?,?,?)",
		user.Name, user.Email, user.Password,
	)

	return err
}