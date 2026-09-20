package internal

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	auth "github.com/montruh-afk/chirpy/internal/authentication"
	"github.com/montruh-afk/chirpy/internal/database"
)

type ApiConfig struct {
	FileServerHits atomic.Int32
	Db             *database.Queries
	Platform       string
	TknScrt        string
	Exp            string
}

type Duration struct {
	duration time.Duration
}

const (
	metrics = `<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`
)

func (cfg *ApiConfig) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.FileServerHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *ApiConfig) RequestHits(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	if _, err := fmt.Fprintf(w, metrics, cfg.FileServerHits.Load()); err != nil {
		log.Printf("Something broke on out end: %s", err)
	}
}

func (cfg *ApiConfig) Reset(w http.ResponseWriter, r *http.Request) {
	if cfg.Platform != "dev" {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("You do not have the permission to interact via this endpoint"))
		return
	}
	ctx := r.Context()
	if err := cfg.Db.ResetUsers(ctx); err != nil {
		log.Panicf("Something went wrong: %s", err)
	}
	cfg.FileServerHits.Store(0)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
}

func (cfg *ApiConfig) CreateUser(w http.ResponseWriter, r *http.Request) {
	createUser, err := verifyDetails(r)
	if err != nil {
		respondWithError(w, 500, "Something broke while attempting to verify input", err)
		return
	}
	
	ctx := r.Context()
	//Check if the user exists to avoid sql unique key constraint violation
	if _, err := cfg.Db.GetuserByEmail(ctx, createUser.Email); err == nil {
		respondWithError(w, http.StatusConflict, "An account with that email already exists", nil)
		return
	}

	hashedPass, err := auth.HashPassword(createUser.Password)
	if err != nil {
		respondWithError(w, http.StatusForbidden, "Something went wrong", err)
		return
	} 
	user, err := cfg.Db.CreateUser(ctx, database.CreateUserParams{
		Email:          createUser.Email,
		HashedPassword: hashedPass,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user", err)
		return
	}

	respondWithJson(w, http.StatusCreated, &user)

}

func (cfg *ApiConfig) CreateChirp(w http.ResponseWriter, r *http.Request) {
	cleanToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusForbidden, "Could not fetch user", err)
		return
	}

	userid, err := auth.ValidateJWT(cleanToken, cfg.TknScrt)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Failed to authenticate user", err)
		return
	}

	log.Printf("Current user id: %v\n", userid)

	params, err := validateChirp(r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	ctx := r.Context()
	//Verify user
	valid, err := cfg.Db.GetUserByID(ctx, userid)
	if err != nil {
		respondWithError(w, 401, "User does not exist", err)
		return
	}

	
	data, err := cfg.Db.CreateChirp(ctx, database.CreateChirpParams{
		UserID: valid.ID,
		Body:   params.Body,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not create chirp", err)
		return
	}

	respondWithJson(w, http.StatusCreated, &data)
}

func (cfg *ApiConfig) GetChirps(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	chirps, err := cfg.Db.GetChirps(ctx)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong while attempting to retrieve from server", err)
		return
	}

	respondWithJson(w, 200, &chirps)
}

func (cfg *ApiConfig) GetChirp(w http.ResponseWriter, r *http.Request) {
	chirp, err := fetchChirp(cfg, r)
	if err != nil {
		respondWithError(w, 404, "Something went wrong", err)
		return
	}
	
	respondWithJson(w, 200, &chirp)
}

func (cfg *ApiConfig) HandlerLogin(w http.ResponseWriter, r *http.Request) {
	//user instance
	deets, err := verifyDetails(r)
	if err != nil {
		respondWithError(w, 500, "Something went wrong while attempting to handle incoming request", err)
		return
	}

	ctx := r.Context()
	userDB, err := cfg.Db.GetuserByEmail(ctx, deets.Email)
	if err != nil {
		respondWithError(w, 401, "Incorrect email or password", nil)
		return
	}

	match, err := auth.CheckPasswordHash(deets.Password, userDB.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
	}

	if !match {
		respondWithError(w, 401, "Incorrect email or password", nil)
		return
	}

	params := &Duration{}
	parseDuration(params, w, cfg.Exp)

	userDB.Token, err = auth.MakeJWT(userDB.ID, cfg.TknScrt, params.duration)
	if err != nil {
		respondWithError(w, 500, "Failed to administer user token", err)
		return
	}

	refresher := auth.MakeRefreshToken()
	if err := cfg.Db.StoreRefreshToken(ctx, database.StoreRefreshTokenParams{
		Token:     refresher,
		UserID:    userDB.ID,
		ExpiresAt: time.Now().AddDate(0, 0, 60),
	}); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not validate token", err)
		return
	}
	userDB.RefreshToken = refresher

	respondWithJson(w, 200, &userDB)

}

func (cfg *ApiConfig) Refresh(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 500, "Could not validate token", err)
		return
	}
	tok, err := cfg.Db.CheckRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, 401, "Failed to authenticate", err)
		return
	}
	if tok.RevokedAt.Valid {
		respondWithError(w, 401, "Token has been revoked", nil)
		return
	}

	if time.Now().Compare(tok.ExpiresAt) >= 0 {
		respondWithError(w, 401, "Expired token, please log in again", nil)
		return
	}

	params := &Duration{}
	//Helper function to avoid repeating time.ParseDuration()
	parseDuration(params, w, cfg.Exp)
	tokn, err := auth.MakeJWT(tok.UserID, cfg.TknScrt, params.duration)
	if err != nil {
		respondWithError(w, 500, "Something went wrong", err)
		return
	}
	type response struct {
		Token string `json:"token"`
	}

	resp := response{
		Token: tokn,
	}

	respondWithJson(w, 200, &resp)
}

func (cfg *ApiConfig) Revoke(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 500, "Could not validate token", err)
		return
	}

	if err := cfg.Db.RevokeToken(r.Context(), token); err != nil {
		respondWithError(w, 500, "Something broke on our end", err)
		return
	}

	respondWithJson(w, 204, nil)

}

func (cfg *ApiConfig) UpdateUserLogin(w http.ResponseWriter, r *http.Request) {
	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "Invalid access token", err)
		return
	}

	userID, err := auth.ValidateJWT(accessToken, cfg.TknScrt)
	if err != nil {
		respondWithError(w, 401, "Invalid access token", err)
		return
	}

	deets, err := verifyDetails(r)
	if err != nil {
		respondWithError(w, 500, "Something broke while attempting to verify input", err)
		return
	}

	hashed_pass, err := auth.HashPassword(deets.Password)
	if err != nil {
		respondWithError(w, 500, "Could not store updated password, your records remain unchanged", err)
		return
	}
	updatedUser, err := cfg.Db.UpdateUser(r.Context(), database.UpdateUserParams{
		Email: deets.Email,
		ID: userID,
		HashedPassword: hashed_pass,
	})
	if err != nil {
		respondWithError(w, 401, "Something went wrong", err)
		return
	}

	respondWithJson(w, 200, &updatedUser)
}

func (cfg *ApiConfig) DeleteChirp(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "Invalid token", err)
		return
	}

	valid, err := auth.ValidateJWT(token, cfg.TknScrt)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorised token", err)
		return
	}

	chirp, err := fetchChirp(cfg, r)
	if err != nil {
		respondWithError(w, 404, "Something went wrong", err)
		return
	}
	if chirp.UserID != valid {
		respondWithError(w, 403, "You do not have the required permissions to carry out that action", nil)
		return
	}
	if err := cfg.Db.DeleteChirp(r.Context(), chirp.ID); err != nil {
		respondWithError(w, 500, "Something went wrong", err)
		return
	}

	respondWithJson(w, 204, nil)
}

func (cfg *ApiConfig) HandlerPolka(w http.ResponseWriter, r *http.Request) {
	type eventHandler struct {
		Event string`json:"event"`
		Data struct {
			UserId string `json:"user_id"`
		}`json:"data"`
	}

	event := eventHandler{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&event); err != nil {
		respondWithError(w, 500, "Something broke on our end", err)
		return
	}
	if event.Event != "user.upgraded" {
		respondWithJson(w, 204, nil)
		return
	}
	
	if err := uuid.Validate(event.Data.UserId); err != nil {
		respondWithError(w, http.StatusForbidden, "Invalid user id format", err)
		return
	}

	id, err := uuid.Parse(event.Data.UserId)
	if err != nil {
		respondWithError(w, http.StatusForbidden, "Invalid user id", err)
		return
	}
	if err := cfg.Db.SetChirpyRed(r.Context(), id); err != nil {
		respondWithError(w, 404, "Something went wrong", err)
		return
	}

	respondWithJson(w, 204, nil)
}