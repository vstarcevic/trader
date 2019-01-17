package postgres

import (
	"database/sql"
	"fmt"
)

// Config defines the options that are used when connecting to a PostgreSQL instance
type Config struct {
	Host string
	Port string
	User string
	Pass string
	Name string
}

// Connect creates a connection to the PostgreSQL instance.
// A non-nil error is returned to indicate failure.
func Connect(cfg Config) (*sql.DB, error) {
	url := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s  sslmode=disable", cfg.Host, cfg.Port, cfg.User, cfg.Name, cfg.Pass)

	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	return db, nil
}
