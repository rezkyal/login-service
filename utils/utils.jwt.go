package utils

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

func GenerateToken(id int64) (string, error) {
	keyData, err := os.ReadFile("./rsakey/jwtrsa256.key")
	if err != nil {
		log.Println("[ERROR][GenerateToken] failed to read private key", err)
		return "", errors.WithStack(errors.New("Error when generate token"))
	}
	key, err := jwt.ParseRSAPrivateKeyFromPEM(keyData)
	if err != nil {
		log.Println("[ERROR][GenerateToken] failed to parse rsa private key from PEM", err)
		return "", errors.WithStack(errors.New("Error when generate token"))
	}

	tokenLifespanStr := os.Getenv("JWT_LIVESPAN")
	tokenLifespan, err := strconv.Atoi(tokenLifespanStr)
	if err != nil {
		log.Println("[WARN][GenerateToken] error when converting tokenLifespan", errors.WithStack(err))
	}
	if tokenLifespan == 0 {
		tokenLifespan = 60
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"id":  id,
		"exp": time.Now().Add(time.Minute * time.Duration(tokenLifespan)).Unix(),
	})

	tokenString, err := token.SignedString(key)

	if err != nil {
		return "", fmt.Errorf("[GenerateToken] error when SignedString, err: %+v", errors.WithStack(err))
	}

	return tokenString, nil
}

func TokenValid(ctx echo.Context) (int64, error) {
	keyData, _ := os.ReadFile("rsakey/jwtrsa256.key.pub")
	key, _ := jwt.ParseRSAPublicKeyFromPEM(keyData)
	tokenString := ExtractToken(ctx)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return key, nil
	})
	if err != nil {
		return 0, errors.WithStack(err)
	}

	claims := token.Claims.(jwt.MapClaims)
	idRaw := claims["id"].(float64)

	return int64(idRaw), nil
}

func ExtractToken(ctx echo.Context) string {
	if len(ctx.Request().Header["Authorization"]) == 0 {
		return ""
	}
	bearerToken := ctx.Request().Header["Authorization"][0]
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}
	return ""
}
