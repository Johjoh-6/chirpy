package main

import (
	"chirpy/internal/response"
	"net/http"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if cfg.Platform != "dev" {
		response.RespondWithError(w, http.StatusForbidden, "Platform is not dev")
		return
	}

	// Delete all users
	if err := cfg.database.DeleteUsers(r.Context()); err != nil {
		response.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Reset"))
}
