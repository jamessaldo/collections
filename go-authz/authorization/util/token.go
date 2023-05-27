package util

import (
	"authorization/config"
	"authorization/domain/model"
	"fmt"
	"time"

	"encoding/base64"

	"github.com/golang-jwt/jwt/v5"
	uuid "github.com/satori/go.uuid"
)

func CreateToken(userid string, ttl time.Duration, privateKey string) (*model.TokenDetails, error) {
	now := time.Now().UTC()
	td := &model.TokenDetails{
		ExpiresIn: new(int64),
		Token:     new(string),
	}
	*td.ExpiresIn = now.Add(ttl).Unix()
	td.TokenUuid = uuid.NewV4().String()
	td.UserID = userid

	decodedPrivateKey, err := base64.StdEncoding.DecodeString(privateKey)
	if err != nil {
		return nil, fmt.Errorf("could not decode token private key: %w", err)
	}
	key, err := jwt.ParseRSAPrivateKeyFromPEM(decodedPrivateKey)

	if err != nil {
		return nil, fmt.Errorf("create: parse token private key: %w", err)
	}

	atClaims := make(jwt.MapClaims)
	atClaims["sub"] = userid
	atClaims["token_uuid"] = td.TokenUuid
	atClaims["exp"] = td.ExpiresIn
	atClaims["iat"] = now.Unix()
	atClaims["nbf"] = now.Unix()

	*td.Token, err = jwt.NewWithClaims(jwt.SigningMethodRS256, atClaims).SignedString(key)
	if err != nil {
		return nil, fmt.Errorf("create: sign token: %w", err)
	}

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

	return &model.TokenDetails{
		TokenUuid: fmt.Sprint(claims["token_uuid"]),
		UserID:    fmt.Sprint(claims["sub"]),
	}, nil
}

func GenerateTokens(user *model.User) (*model.TokenDetails, *model.TokenDetails, error) {
	token, err := CreateToken(user.ID.String(), config.AppConfig.AccessTokenExpiresIn, config.AppConfig.AccessTokenPrivateKey)
	if err != nil {
		return nil, nil, err
	}

	refreshToken, err := CreateToken(user.ID.String(), config.AppConfig.RefreshTokenExpiresIn, config.AppConfig.RefreshTokenPrivateKey)
	if err != nil {
		return nil, nil, err
	}

	return token, refreshToken, nil
}
