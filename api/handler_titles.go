package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/chaeanthony/go-netflix/utils"
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

// Gets all titles
func (cfg *APIConfig) HandlerTitlesGet(w http.ResponseWriter, r *http.Request) {
	titles, err := cfg.DB.GetTitles()
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Couldn't get title", err)
		return
	}

	utils.RespondJSON(w, http.StatusOK, titles)
}

func (cfg *APIConfig) HandlerTitleGetById(w http.ResponseWriter, r *http.Request) {
	titleIdString := r.PathValue("titleId")
	titleId, err := strconv.Atoi(titleIdString)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "failed to parse title id", err)
		return
	}

	title, err := cfg.DB.GetTitleById(titleId)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Couldn't get title", err)
		return
	}

	utils.RespondJSON(w, http.StatusOK, title)
}

// Gets all titles of type show
func (cfg *APIConfig) HandlerShowsGet(w http.ResponseWriter, r *http.Request) {
	shows, err := cfg.DB.GetShows()
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Couldn't get shows", err)
		return
	}

	utils.RespondJSON(w, http.StatusOK, shows)
}

// Gets all titles of type movie
func (cfg *APIConfig) HandlerMoviesGet(w http.ResponseWriter, r *http.Request) {
	shows, err := cfg.DB.GetMovies()
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Couldn't get shows", err)
		return
	}

	utils.RespondJSON(w, http.StatusOK, shows)
}
