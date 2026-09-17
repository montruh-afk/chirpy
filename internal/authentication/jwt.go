package authentication

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenType string

const (
	TokenAccess TokenType = "chirpy-access"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer: string(TokenAccess),
		IssuedAt: &jwt.NumericDate{
			Time: time.Now().UTC(),
		},
		ExpiresAt: &jwt.NumericDate{
			Time: time.Now().Add(expiresIn),
		},
		Subject: userID.String(),
	})
	signed, err := token.SignedString([]byte(tokenSecret))
	return signed, err
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claims := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	})

	if err != nil {
		return uuid.Nil, err
	}

	idString, err:= token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, err
	}

	if issuer, err := token.Claims.GetIssuer(); err != nil {
		return uuid.Nil, err
	} else if issuer != string(TokenAccess) {
		return uuid.Nil, fmt.Errorf("Tampered token")
	}

	userID, err := uuid.Parse(idString)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid user ID: %w", err)
	}
	return userID, nil
}
