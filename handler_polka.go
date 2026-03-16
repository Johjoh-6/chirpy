package main

import (
	"chirpy/internal/database"
	"chirpy/internal/response"
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerPolkaWH(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}
	var p parameters
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		response.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if p.Event != "user.upgraded" {
		response.RespondWithJSON(w, http.StatusNoContent, nil)
		return
	} else {
		userID, err := uuid.Parse(p.Data.UserID)
		if err != nil {
			response.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		if _, err = cfg.database.UpdateUserIsChirpyRed(r.Context(), database.UpdateUserIsChirpyRedParams{
			ID:          userID,
			IsChirpyRed: true,
		}); err != nil {
			if err == sql.ErrNoRows {
				response.RespondWithError(w, http.StatusNotFound, "user not found")
				return
			}
			response.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		response.RespondWithJSON(w, http.StatusNoContent, nil)
	}

}
