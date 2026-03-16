package main

import (
	"chirpy/internal/auth"
	"chirpy/internal/database"
	"chirpy/internal/response"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"sort"
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

func (cfg *apiConfig) handlerChirp(w http.ResponseWriter, r *http.Request) {

	userToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		response.RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}
	userID, err := auth.ValidateJWT(userToken, cfg.JWTSecret)
	if err != nil {
		response.RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	// read the body request
	type parameters struct {
		Body string `json:"body"`
	}
	var c parameters
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
		UserID: userID,
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

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	authorIDString := r.URL.Query().Get("author_id")
	sortString := r.URL.Query().Get("sort")
	if sortString != "" && sortString != "asc" && sortString != "desc" {
		response.RespondWithError(w, http.StatusBadRequest, "invalid sort")
		return
	}

	authorID, err := uuid.Parse(authorIDString)
	if err != nil && authorIDString != "" {
		response.RespondWithError(w, http.StatusBadRequest, "invalid author_id")
		return
	}
	var chirps []database.Chirp
	var error error
	if authorIDString != "" {
		chirps, error = cfg.database.GetChirpsByAuthor(context.Background(), authorID)
	} else {
		chirps, error = cfg.database.GetChirps(context.Background())
	}
	if error != nil {
		response.RespondWithError(w, http.StatusInternalServerError, error.Error())
		return
	}
	// change chirps from database to Chirps struct
	var chirpResponses []Chirp
	for _, chirp := range chirps {
		chirpResponses = append(chirpResponses, Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		})
	}

	// sort in memory
	if sortString == "desc" {
		sort.Slice(chirpResponses, func(i, j int) bool {
			return chirpResponses[i].CreatedAt.After(chirpResponses[j].CreatedAt)
		})
	}
	response.RespondWithJSON(w, http.StatusOK, chirpResponses)
}

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	chirpIDString := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDString)
	if err != nil {
		response.RespondWithError(w, http.StatusBadRequest, "chirpID is required")
		return
	}
	chirp, err := cfg.database.GetChirp(context.Background(), chirpID)
	if err != nil {
		if err == sql.ErrNoRows {
			response.RespondWithError(w, http.StatusNotFound, "chirp not found")
			return
		}
		response.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.RespondWithJSON(w, http.StatusOK, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
}

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	chirpIDString := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDString)
	if err != nil {
		response.RespondWithError(w, http.StatusBadRequest, "chirpID is required")
		return
	}
	userToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		response.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	userID, err := auth.ValidateJWT(userToken, cfg.JWTSecret)
	if err != nil {
		response.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	chirp, err := cfg.database.GetChirp(context.Background(), chirpID)
	if err != nil {
		if err == sql.ErrNoRows {
			response.RespondWithError(w, http.StatusNotFound, "chirp not found")
			return
		}
		response.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if chirp.UserID != userID {
		response.RespondWithError(w, http.StatusForbidden, "Forbidden")
		return
	}
	err = cfg.database.DeleteChirp(context.Background(), chirpID)
	if err != nil {
		response.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.RespondWithJSON(w, http.StatusNoContent, nil)
}
