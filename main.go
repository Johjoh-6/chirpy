package main

import (
	"chirpy/internal/database"
	"chirpy/internal/response"
	"chirpy/internal/utils"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	profanity      *utils.Profanity
	database       *database.Queries
}

func main() {
	// load the environment variables from .env
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	dbQueries := database.New(db)

	mux := http.NewServeMux()

	svr := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	apiCfg := &apiConfig{
		profanity: &utils.Profanity{
			Word:     []string{"kerfuffle", "sharbert", "fornax"},
			Replacer: "****",
		},
		database: dbQueries,
	}

	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir("./")))))
	mux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetric)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)

	mux.HandleFunc("POST /api/validate_chirp", apiCfg.handlerValidateChirp)

	svr.ListenAndServe()
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handlerMetric(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	tmpl := fmt.Sprintf(`
	<html>
	  <body>
	    <h1>Welcome, Chirpy Admin</h1>
	    <p>Chirpy has been visited %d times!</p>
	  </body>
	</html>
	`, cfg.fileserverHits.Load())
	w.Write([]byte(tmpl))
}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Reset"))
}

func (cfg *apiConfig) handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	// read the body request
	type chirp struct {
		Body string `json:"body"`
	}
	var c chirp
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

	response.RespondWithJSON(w, http.StatusOK, map[string]string{"cleaned_body": c.Body})
}
