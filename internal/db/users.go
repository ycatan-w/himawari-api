package db

import (
	"database/sql"
	"time"
)

type User struct {
	Username string
	Password string
}

type PartialUserDB struct {
	ID       int
	Username string
	Password string
}

func UserExists(username string) (bool, error) {
	var exists bool
	err := DB.QueryRow(`
	SELECT 1
	FROM users
	WHERE username = ?`, username).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func FindUserByUsername(username string) *PartialUserDB {
	var partialUser PartialUserDB

	err := DB.QueryRow(`
	SELECT id, username, password
	FROM users
	WHERE username = ?
	`, username).Scan(&partialUser.ID, &partialUser.Username, &partialUser.Password)
	if err != nil {
		return nil
	}

	return &partialUser
}

func CreateUser(user User) (int, error) {
	result, err := DB.Exec(`
	INSERT INTO users (username, password, created_at)
	VALUES (?, ?, ?)`, user.Username, user.Password, time.Now().UTC().Format(time.RFC3339))
	if err != nil {
		return -1, err
	}
	id, _ := result.LastInsertId()

	return int(id), nil
}
