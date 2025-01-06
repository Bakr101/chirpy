package auth

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error){
	hashedPassword, err :=  bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func CheckPasswordHash(password, hash string) error{
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil{
		return err
	}
	return nil
}

type CustomClaims struct{
	jwt.RegisteredClaims
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error){
	claims := CustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: "chirpy",
			IssuedAt: jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn * time.Second)),
			Subject: userID.String(),
		},
	}
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	stringifiedToken, err := newToken.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}
	return stringifiedToken, nil
}



func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error){
	parsedToken, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok{
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(tokenSecret), nil
	})

	if err != nil {
		return uuid.Nil, fmt.Errorf("error parsing token: %w", err)
	}
	claims, ok := parsedToken.Claims.(*CustomClaims)
	if !ok || !parsedToken.Valid{
		return uuid.Nil, fmt.Errorf("err validating token claims: %v", parsedToken.Claims)
	}
	parsedUUID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid subject UUID: %w", err)
	} 

	return parsedUUID, nil
}

func GetBearerToken(headers http.Header) (string, error){
	authorizationValues := headers.Values("Authorization")
	if len(authorizationValues) == 0 {
		return "", fmt.Errorf("no values in authorization header")
	}
	authValuesSplit := strings.Fields(authorizationValues[0])

	trimmedToken := strings.TrimSpace(authValuesSplit[1]) 

	return trimmedToken, nil
}