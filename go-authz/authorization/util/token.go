package util

import (
	"authorization/config"
	"authorization/domain/model"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(ttl time.Duration, payload interface{}, secretJWTKey, jwtKid string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	now := time.Now().UTC()
	claims := token.Claims.(jwt.MapClaims)

	claims["sub"] = payload.(*model.User).ID
	claims["exp"] = now.Add(ttl).Unix()
	claims["iat"] = now.Unix()
	claims["nbf"] = now.Unix()

	token.Header["kid"] = jwtKid

	tokenString, err := token.SignedString([]byte(secretJWTKey))

	if err != nil {
		return "", fmt.Errorf("generating JWT Token failed: %w", err)
	}

	return tokenString, nil
}

func GenerateRefreshToken(ttl time.Duration, payload interface{}, secretRefreshJWTKey string) (string, error) {
	refreshToken := jwt.New(jwt.SigningMethodHS256)

	now := time.Now().UTC()
	claims := refreshToken.Claims.(jwt.MapClaims)

	claims["sub"] = payload
	claims["exp"] = now.Add(ttl).Unix()

	refreshTokenString, err := refreshToken.SignedString([]byte(secretRefreshJWTKey))

	if err != nil {
		return "", fmt.Errorf("generating JWT Refresh Token failed: %w", err)
	}

	return refreshTokenString, nil
}

func ValidateToken(token string, signedJWTKey string) (interface{}, error) {
	tok, err := jwt.Parse(token, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected method: %s", jwtToken.Header["alg"])
		}

		return []byte(signedJWTKey), nil
	})
	if err != nil {
		return nil, fmt.Errorf("invalidate token: %w", err)
	}

	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok || !tok.Valid {
		return nil, fmt.Errorf("invalid token claim")
	}

	return claims["sub"], nil
}

func GenerateTokens(user *model.User) (string, string, error) {
	token, err := GenerateToken(config.AppConfig.TokenExpiresIn, user, config.AppConfig.JWTTokenSecret, config.AppConfig.JWTKid)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := GenerateRefreshToken(config.AppConfig.RefreshTokenExpiresIn, user, config.AppConfig.RefreshJWTTokenSecret)
	if err != nil {
		return "", "", err
	}

	return token, refreshToken, nil
}
