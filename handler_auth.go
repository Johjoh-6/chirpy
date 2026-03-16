package main

import (
	"chirpy/internal/auth"
	"chirpy/internal/response"
	"encoding/json"
	"net/http"
)

type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	var userLogin UserLogin
	if err := json.NewDecoder(r.Body).Decode(&userLogin); err != nil {
		response.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	if len(userLogin.Email) == 0 {
		response.RespondWithError(w, http.StatusBadRequest, "Email is required")
		return
	}
	if len(userLogin.Password) == 0 {
		response.RespondWithError(w, http.StatusBadRequest, "Password is required")
		return
	}

	user, err := cfg.database.GetUserByEmail(r.Context(), userLogin.Email)
	if err != nil {
		response.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	ok, err := auth.CheckPasswordHash(userLogin.Password, user.HashedPassword)
	if err != nil {
		response.RespondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}
	if !ok {
		response.RespondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	response.RespondWithJSON(w, http.StatusOK, User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	})
}
