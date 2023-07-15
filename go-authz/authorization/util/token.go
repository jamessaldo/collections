package util

import (
	"authorization/config"
	"authorization/domain/model"
	"fmt"
	"time"

	"encoding/base64"

	"github.com/golang-jwt/jwt/v5"
	"github.com/oklog/ulid/v2"
	uuid "github.com/satori/go.uuid"
)

func CreateToken(userID uuid.UUID, ttl time.Duration, privateKey string) (*model.TokenDetails, error) {
	now := time.Now().UTC()

	decodedPrivateKey, err := base64.StdEncoding.DecodeString(privateKey)
	if err != nil {
		return nil, fmt.Errorf("could not decode token private key: %w", err)
	}
	key, err := jwt.ParseRSAPrivateKeyFromPEM(decodedPrivateKey)

	if err != nil {
		return nil, fmt.Errorf("create: parse token private key: %w", err)
	}

	atClaims := make(jwt.MapClaims)
	atClaims["sub"] = userID
	atClaims["token_ulid"] = ulid.Make()
	atClaims["exp"] = now.Add(ttl).Unix()
	atClaims["iat"] = now.Unix()
	atClaims["nbf"] = now.Unix()

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodRS256, atClaims)
	jwtToken.Header["kid"] = config.AppConfig.AccessTokenKID

	token, err := jwtToken.SignedString(key)
	if err != nil {
		return nil, fmt.Errorf("create: sign token: %w", err)
	}

	td := model.NewTokenDetails(token, atClaims["token_ulid"].(ulid.ULID), userID, atClaims["exp"].(int64))
	return td, nil
}

func ValidateToken(token string, publicKey string) (*model.TokenDetails, error) {
	decodedPublicKey, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		return nil, fmt.Errorf("could not decode: %w", err)
	}

	key, err := jwt.ParseRSAPublicKeyFromPEM(decodedPublicKey)

	if err != nil {
		return nil, fmt.Errorf("validate: parse key: %w", err)
	}

	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected method: %s", t.Header["alg"])
		}
		return key, nil
	})

	if err != nil {
		return nil, fmt.Errorf("validate: %w", err)
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		return nil, fmt.Errorf("validate: invalid token")
	}

	tokenUlid, err := ulid.Parse(fmt.Sprint(claims["token_ulid"]))
	if err != nil {
		return nil, fmt.Errorf("validate: parse token ulid: %w", err)
	}
	expirationIn, err := claims.GetExpirationTime()
	if err != nil {
		return nil, fmt.Errorf("validate: get expiration time: %w", err)
	}

	td := model.NewTokenDetails(token, tokenUlid, uuid.FromStringOrNil(fmt.Sprint(claims["sub"])), expirationIn.Unix())

	return td, nil
}

func GenerateTokens(user *model.User) (*model.TokenDetails, *model.TokenDetails, error) {
	token, err := CreateToken(user.ID, config.AppConfig.AccessTokenExpiresIn, config.AppConfig.AccessTokenPrivateKey)
	if err != nil {
		return nil, nil, err
	}

	refreshToken, err := CreateToken(user.ID, config.AppConfig.RefreshTokenExpiresIn, config.AppConfig.RefreshTokenPrivateKey)
	if err != nil {
		return nil, nil, err
	}

	return token, refreshToken, nil
}
