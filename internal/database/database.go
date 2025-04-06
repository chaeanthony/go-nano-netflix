package database

import (
	"database/sql"
	"errors"
	"fmt"
)

type Client struct {
	db *sql.DB
}

func NewClient(pathToDB string) (*Client, error) {
	db, err := sql.Open("postgres", pathToDB) // pathToDB = conn string in postgres
	if err != nil {
		return nil, err
	}
	c := &Client{db}
	err = c.autoMigrate()
	if err != nil {
		return nil, errors.New("failed to auto migrate. " + err.Error())
	}
	return c, nil
}

func (c *Client) autoMigrate() error {
	userTable := `
	CREATE TABLE IF NOT EXISTS users (
		id UUID PRIMARY KEY,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		password TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE
	);
	`
	_, err := c.db.Exec(userTable)
	if err != nil {
		return err
	}

	refreshTokenTable := `
	CREATE TABLE IF NOT EXISTS refresh_tokens (
		token TEXT PRIMARY KEY,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		revoked_at TIMESTAMPTZ,
		user_id UUID NOT NULL,
		expires_at TIMESTAMPTZ NOT NULL,
		FOREIGN KEY(user_id) REFERENCES users(id)
	);
	`
	_, err = c.db.Exec(refreshTokenTable)
	if err != nil {
		return err
	}

	mediaTitleTable := `
	CREATE TABLE IF NOT EXISTS media_titles (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
		type VARCHAR(15) NOT NULL, 
		description TEXT,
		origin_date TIMESTAMP,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err = c.db.Exec(mediaTitleTable)
	if err != nil {
		return err
	}

	watchlistTable := `
	CREATE TABLE watchlists (
    watchlist_id SERIAL PRIMARY KEY, 
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,            
    title_id INT NOT NULL REFERENCES media_titles(id),           
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err = c.db.Exec(watchlistTable)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) Reset() error {
	if _, err := c.db.Exec("DELETE FROM refresh_tokens"); err != nil {
		return fmt.Errorf("failed to reset table refresh_tokens: %w", err)
	}
	if _, err := c.db.Exec("DELETE FROM users"); err != nil {
		return fmt.Errorf("failed to reset table users: %w", err)
	}
	if _, err := c.db.Exec("DELETE FROM media_titles"); err != nil {
		return fmt.Errorf("failed to reset table media_titles: %w", err)
	}
	return nil
}
