package authentication

import (
	"github.com/alexedwards/argon2id"
	"runtime"
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