package database

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	CreateUserParams
}

type CreateUserParams struct {
	Email    string `json:"email"`
	Password string `json:"-"`
}

type UpdateUserParams struct {
	ID             uuid.UUID
	Email          string
	HashedPassword string
}

func (c *Client) GetUsers() ([]User, error) {
	query := `
		SELECT
			id,
			email, created_at, updated_at
		FROM users
	`

	rows, err := c.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []User{}
	for rows.Next() {
		var user User
		var id string
		if err := rows.Scan(&id, &user.Email); err != nil {
			return nil, err
		}
		user.ID, err = uuid.Parse(id)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (c *Client) GetUserByEmail(email string) (User, error) {
	query := `
		SELECT id, created_at, updated_at, email, password
		FROM users
		WHERE email = $1
	`
	var user User
	var id string
	err := c.db.QueryRow(query, email).Scan(&id, &user.CreatedAt, &user.UpdatedAt, &user.Email, &user.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, nil
		}
		return User{}, err
	}
	user.ID, err = uuid.Parse(id)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (c *Client) GetUserByRefreshToken(token string) (User, error) {
	query := `
		SELECT u.id, u.email, u.created_at, u.updated_at
		FROM users u
		JOIN refresh_tokens rt ON u.id = rt.user_id
		WHERE rt.token = $1
	`

	var user User
	var id string
	err := c.db.QueryRow(query, token).Scan(&id, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, nil
		}
		return User{}, err
	}
	user.ID, err = uuid.Parse(id)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (c *Client) CreateUser(params CreateUserParams) (User, error) {
	id := uuid.New()

	query := `
		INSERT INTO users
		    (id, created_at, updated_at, email, password)
		VALUES
		    ($1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, $2, $3)
	`
	_, err := c.db.Exec(query, id.String(), params.Email, params.Password)
	if err != nil {
		return User{}, err
	}

	return c.GetUserById(id)
}

func (c *Client) GetUserById(id uuid.UUID) (User, error) {
	query := `
		SELECT id, created_at, updated_at, email, password
		FROM users
		WHERE id = $1
	`
	var user User
	var idStr string
	err := c.db.QueryRow(query, id.String()).Scan(&idStr, &user.CreatedAt, &user.UpdatedAt, &user.Email, &user.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, nil
		}
		return User{}, err
	}
	user.ID, err = uuid.Parse(idStr)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (c *Client) UpdateUser(params UpdateUserParams) (User, error) {
	query := `UPDATE users SET email = $2, hashed_password = $3 WHERE id = $1`
	_, err := c.db.Exec(query, params.ID.String(), params.Email, params.HashedPassword)
	if err != nil {
		return User{}, err
	}

	return c.GetUserById(params.ID)
}

func (c *Client) DeleteUser(id uuid.UUID) error {
	query := `
		DELETE FROM users
		WHERE id = $2
	`
	_, err := c.db.Exec(query, id.String())
	return err
}
