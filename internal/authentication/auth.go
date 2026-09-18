package authentication

import (
	"github.com/alexedwards/argon2id"
	"runtime"
	"net/http"
	"errors"
	"strings"
)

func HashPassword(pass string) (string, error) {
	params := &argon2id.Params{
		Parallelism: uint8(runtime.NumCPU()),
		Iterations: argon2id.DefaultParams.Iterations,
		Memory: argon2id.DefaultParams.Memory,
		SaltLength: argon2id.DefaultParams.SaltLength,
		KeyLength: argon2id.DefaultParams.KeyLength,
	}

	hash, err := argon2id.CreateHash(pass, params)
	if err != nil {
		return "", err
	}
	return hash, nil
}

func CheckPasswordHash(password, hash string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(password, hash)
	return match, err
}

func GetBearerToken(headers http.Header) (string, error) {
	auth := headers.Get("Authorization")

	if auth == "" {
		return auth, errors.New("Unauthorised: Please login in to access features")
	}
	if !strings.HasPrefix(strings.ToLower(auth), "bearer ") {
		return "", errors.New("Malformed authorization header received")
	}
	token := auth[len("bearer "):]
	return strings.TrimSpace(token), nil
}

