package main

import (
	"chirpy/internal/database"
	"chirpy/internal/response"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

type chirpInsert struct {
	Body   string    `json:"body"`
	UserID uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerChirp(w http.ResponseWriter, r *http.Request) {
	// read the body request
	var c chirpInsert
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		response.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	if len(c.Body) == 0 {
		response.RespondWithError(w, http.StatusBadRequest, "Chirp is required")
		return
	}
	if len(c.Body) > 140 {
		response.RespondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}
	// replace profanity
	c.Body = cfg.profanity.RemoveProfanity(c.Body)

	// store chirp in database
	chirpCreated, err := cfg.database.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   c.Body,
		UserID: c.UserID,
	})
	if err != nil {
		response.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.RespondWithJSON(w, http.StatusCreated, Chirp{
		ID:        chirpCreated.ID,
		CreatedAt: chirpCreated.CreatedAt,
		UpdatedAt: chirpCreated.UpdatedAt,
		Body:      chirpCreated.Body,
		UserID:    chirpCreated.UserID,
	})
}
