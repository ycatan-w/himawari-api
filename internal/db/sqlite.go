package db

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/ycatan-w/himawari-api/internal/output"
	_ "modernc.org/sqlite"
)

var DB *sql.DB

func ConnectDB() error {
	err := openSQL()
	if err != nil {
		output.PrintFail("Failed to open database")
		return fmt.Errorf("failed to open database: %w", err)
	}

	if !tableExists("users") {
		output.PrintFail("Database must be initialize first")
		return errors.New("database not initialized, please run with --init-db first")
	}

	return nil
}
func InitDB() error {
	err := openSQL()
	if err != nil {
		output.PrintFail("Failed to open database")
		return fmt.Errorf("failed to open database: %w", err)
	}

	if err := createTables(); err != nil {
		output.PrintFail("Table creation fail")
		return fmt.Errorf("failed to create tables: %w", err)
	}

	return nil
}

var DBPath string

func openSQL() error {
	var err error
	dbPath := os.Getenv("HIMAWARI_DB_PATH")
	if len(dbPath) > 0 {
		DBPath = dbPath
	} else {
		if runtime.GOOS == "windows" {
			base := os.Getenv("ProgramData")
			if base == "" {
				base = "C:\\ProgramData"
			}
			DBPath = filepath.Join(base, "Himawari", "himawari.db")
		} else {
			DBPath = "/var/lib/himawari/himawari.db"
		}
	}
	dir := filepath.Dir(DBPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	DB, err = sql.Open("sqlite", DBPath)
	if err != nil {
		return err
	}

	return nil
}

func tableExists(name string) bool {
	query := `SELECT name FROM sqlite_master WHERE type='table' AND name=?;`
	row := DB.QueryRow(query, name)
	var table string
	if err := row.Scan(&table); err != nil {
		return false
	}
	return table == name
}

func createTables() error {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username VARCHAR(100) NOT NULL UNIQUE,
			password TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);`,

		`CREATE TABLE IF NOT EXISTS user_tokens (
			token TEXT PRIMARY KEY,
			user_id INTEGER NOT NULL,
			created_at DATETIME NOT NULL,
			last_used DATETIME,
			expires_at DATETIME,
			FOREIGN KEY(user_id) REFERENCES users(id)
		);`,
		`CREATE INDEX IF NOT EXISTS idx_user_tokens_user_id ON user_tokens(user_id);`,

		`CREATE TABLE IF NOT EXISTS events (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title VARCHAR(100) NOT NULL,
			description TEXT,
			date VARCHAR(10) NOT NULL,
			start INTEGER NOT NULL,
			end INTEGER NOT NULL,
			user_id INTEGER NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY(user_id) REFERENCES users(id)
		);`,
		`CREATE INDEX IF NOT EXISTS idx_events_user_date ON events(user_id, date);`,

		`CREATE TABLE IF NOT EXISTS logs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			text TEXT NOT NULL,
			date VARCHAR(10) NOT NULL,
			user_id INTEGER NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY(user_id) REFERENCES users(id)
		);`,
		`CREATE INDEX IF NOT EXISTS idx_logs_user_date ON logs(user_id, date);`,
	}

	for _, stmt := range stmts {
		if _, err := DB.Exec(stmt); err != nil {
			return err
		}
	}

	return nil
}
