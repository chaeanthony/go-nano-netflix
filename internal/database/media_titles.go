package database

import "time"

const (
	TypeShow  = "show"
	TypeMovie = "movie"
)

type MediaTitle struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	OriginDate  time.Time `json:"origin_date"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (c *Client) GetTitles() ([]MediaTitle, error) {
	query := `SELECT id, title, type, description, origin_date, created_at, updated_at FROM media_titles`

	rows, err := c.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	titles := []MediaTitle{}
	for rows.Next() {
		var title MediaTitle
		if err := rows.Scan(&title.ID, &title.Title, &title.Type, &title.Description, &title.OriginDate, &title.CreatedAt, &title.UpdatedAt); err != nil {
			return nil, err
		}
		titles = append(titles, title)
	}

	return titles, nil
}

func (c *Client) getTitlesByType(titleType string) ([]MediaTitle, error) {
	query := `SELECT id, title, type, description, origin_date, created_at, updated_at FROM media_titles WHERE type = $1`

	rows, err := c.db.Query(query, titleType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	titles := []MediaTitle{}
	for rows.Next() {
		var title MediaTitle
		if err := rows.Scan(&title.ID, &title.Title, &title.Type, &title.Description, &title.OriginDate, &title.CreatedAt, &title.UpdatedAt); err != nil {
			return nil, err
		}
		titles = append(titles, title)
	}

	return titles, nil
}

func (c *Client) GetTitleById(id int) (MediaTitle, error) {
	query := `SELECT id, title, type, description, origin_date, created_at, updated_at FROM media_titles WHERE id = $1`

	var title MediaTitle
	err := c.db.QueryRow(query, id).Scan(&title.ID, &title.Title, &title.Type, &title.Description, &title.OriginDate, &title.CreatedAt, &title.UpdatedAt)
	if err != nil {
		return MediaTitle{}, err
	}

	return title, nil
}

func (c *Client) GetShows() ([]MediaTitle, error) {
	return c.getTitlesByType(TypeShow)
}

func (c *Client) GetMovies() ([]MediaTitle, error) {
	return c.getTitlesByType(TypeMovie)
}

func (c *Client) CreateTitle(params MediaTitle) error {
	query := `
		INSERT INTO media_titles 
			(title, type, description, origin_date, created_at, updated_at) 
		VALUES 
			($1, $2, $3, $4, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
	`

	_, err := c.db.Exec(query, params.Title, params.Type, params.Description, params.OriginDate)
	return err
}
