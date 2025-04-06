package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type WatchlistItem struct {
	ID         int        `json:"watchlist_id"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	MediaTitle MediaTitle `JSON:"media_title"`
	CreateWatchlistItem
}

type CreateWatchlistItem struct {
	UserID  uuid.UUID `json:"user_id"`
	TitleID int       `json:"title_id"`
}

func (c *Client) GetWatchlist(userID uuid.UUID) ([]WatchlistItem, error) {
	query := `
	SELECT 
		w.ID, 
		w.created_at AS watchlist_created_at, 
		w.updated_at AS watchlist_updated_at,
		t.id AS media_title_id,
		t.title,
		t.type,
		t.description,
		t.origin_date,
		t.created_at AS media_title_created_at,
		t.updated_at AS media_title_updated_at
	FROM watchlists w
	JOIN users u ON w.user_id = u.id
	JOIN media_titles t ON w.title_id = t.id
	WHERE u.user_id = $1;
	`

	rows, err := c.db.Query(query, userID.String())
	if err != nil {
		return []WatchlistItem{}, err
	}
	defer rows.Close()

	var watchlistItems []WatchlistItem
	for rows.Next() {
		var item WatchlistItem
		var mediaTitle MediaTitle

		err := rows.Scan(
			&item.ID,
			&item.TitleID,
			&item.CreatedAt, // Watchlist created at
			&item.UpdatedAt, // Watchlist updated at
			&mediaTitle.ID,
			&mediaTitle.Title,
			&mediaTitle.Type,
			&mediaTitle.Description,
			&mediaTitle.OriginDate,
			&mediaTitle.CreatedAt, // Media title created at
			&mediaTitle.UpdatedAt, // Media title updated at
		)
		if err != nil {
			return []WatchlistItem{}, err
		}

		item.UserID = userID
		item.MediaTitle = mediaTitle
		watchlistItems = append(watchlistItems, item)
	}

	if err := rows.Err(); err != nil {
		return []WatchlistItem{}, err
	}

	return watchlistItems, nil
}

func (c *Client) getWatchlistItemById(id int) (WatchlistItem, error) {
	query := `
	SELECT 
    w.id AS watchlist_id, 
    w.created_at AS watchlist_created_at, 
    w.updated_at AS watchlist_updated_at,
    t.id AS media_title_id,
    t.title,
    t.type,
    t.description,
    t.origin_date,
    t.created_at AS media_title_created_at,
    t.updated_at AS media_title_updated_at
	FROM watchlists w
	JOIN media_titles t ON w.title_id = t.id
	WHERE w.id = $1;
	`

	var watchlistItem WatchlistItem
	err := c.db.QueryRow(query, id).Scan(
		&watchlistItem.ID,
		&watchlistItem.CreatedAt,
		&watchlistItem.UpdatedAt,
		&watchlistItem.MediaTitle.ID,
		&watchlistItem.MediaTitle.Title,
		&watchlistItem.MediaTitle.Type,
		&watchlistItem.MediaTitle.Description,
		&watchlistItem.MediaTitle.OriginDate,
		&watchlistItem.MediaTitle.CreatedAt,
		&watchlistItem.MediaTitle.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return WatchlistItem{}, fmt.Errorf("no watchlist found with id %d", id)
		}
		return WatchlistItem{}, err
	}

	return watchlistItem, nil
}

func (c *Client) existsWatchlistItem(userID uuid.UUID, titleID int) (bool, error) {
	query := `
	SELECT EXISTS (
    SELECT 1 
    FROM watchlists 
    WHERE user_id = $1 AND title_id = $2
	);
	`
	var exists bool
	err := c.db.QueryRow(query, userID.String(), titleID).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (c *Client) AddWatchlistItem(userID uuid.UUID, titleID int) (WatchlistItem, error) {
	exists, err := c.existsWatchlistItem(userID, titleID)
	if err != nil {
		return WatchlistItem{}, err
	} else if exists {
		return WatchlistItem{}, fmt.Errorf("title %d already exists", titleID)
	}

	query := `
	INSERT INTO watchlists 
		(user_id, title_id, created_at, updated_at)
	VALUES 
		($1, $2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	RETURNING id;
	`
	var watchlistID int
	err = c.db.QueryRow(query, userID.String(), titleID).Scan(&watchlistID)
	if err != nil {
		return WatchlistItem{}, err
	}

	return c.getWatchlistItemById(watchlistID)
}
