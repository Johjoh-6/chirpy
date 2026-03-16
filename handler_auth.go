package main

import (
	"chirpy/internal/auth"
	"chirpy/internal/database"
	"chirpy/internal/response"
	"encoding/json"
	"net/http"
	"time"
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

	token, err := auth.MakeJWT(
		user.ID,
		cfg.JWTSecret,
		time.Hour,
	)
	if err != nil {
		response.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	rftString := auth.MakeRefreshToken()

	refreshToken, err := cfg.database.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		UserID:    user.ID,
		Token:     rftString,
		ExpiresAt: time.Now().AddDate(0, 0, 60),
	})
	if err != nil {
		response.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.RespondWithJSON(w, http.StatusOK, User{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		Token:        token,
		RefreshToken: refreshToken.Token,
	})
}

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	tokenBearer, err := auth.GetBearerToken(r.Header)
	if err != nil {
		response.RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	// get the token from refresh_token
	refreshToken, err := cfg.database.GetUserFromRefreshToken(r.Context(), tokenBearer)
	if err != nil {
		response.RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	token, err := auth.MakeJWT(
		refreshToken.ID,
		cfg.JWTSecret,
		time.Hour,
	)
	if err != nil {
		response.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.RespondWithJSON(w, http.StatusOK, map[string]string{
		"token": token,
	})
}

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	tokenBearer, err := auth.GetBearerToken(r.Header)
	if err != nil {
		response.RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}
	if err = cfg.database.RevokeRefreshToken(r.Context(), tokenBearer); err != nil {
		response.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
