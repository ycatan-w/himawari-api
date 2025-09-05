package db

import (
	"database/sql"
	"time"
)

var ttl time.Duration = 30 * 24 * time.Hour

func AddUserToken(userID int, token string) error {
	createdAt := time.Now().UTC()
	expiresAt := createdAt.Add(ttl)

	_, err := DB.Exec(`
		INSERT INTO user_tokens (token, user_id, created_at, last_used, expires_at)
		VALUES (?, ?, ?, ?, ?)
	`, token, userID, createdAt.Format(time.RFC3339), createdAt.Format(time.RFC3339), expiresAt.Format(time.RFC3339))

	return err
}

func UpdateLastUsed(token string) error {
	lastUsed := time.Now().UTC()
	expiresAt := lastUsed.Add(ttl)

	_, err := DB.Exec(`
		UPDATE user_tokens
		SET last_used = ?, expires_at = ?
		WHERE token = ?
	`, lastUsed.Format(time.RFC3339), expiresAt.Format(time.RFC3339), token)

	return err
}

func RemoveUserToken(token string) (int64, error) {
	res, err := DB.Exec(`
		DELETE FROM user_tokens
		WHERE token = ?
	`, token)

	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}

func PurgeExpiredTokens() (int64, error) {
	res, err := DB.Exec(`
		DELETE FROM user_tokens
		WHERE expires_at < ?
	`, time.Now().UTC().Format(time.RFC3339))

	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}

func ValidateToken(token string) (int, error) {
	var userID int

	err := DB.QueryRow(`
		SELECT user_id FROM user_tokens
		WHERE token = ?
	`, token).Scan(&userID)

	if err == sql.ErrNoRows {
		return -1, nil
	}
	if err != nil {
		return -1, err
	}

	return userID, nil
}
