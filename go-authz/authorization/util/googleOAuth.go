package util

import (
	"auth/config"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"time"
)

type GoogleOauthToken struct {
	Access_token string
	Id_token     string
}

type GoogleUserResult struct {
	Id            string
	Email         string
	VerifiedEmail bool
	Name          string
	GivenName     string
	FamilyName    string
	Picture       string
	Locale        string
}

func GetGoogleOauthToken(ctx context.Context, code string) (*GoogleOauthToken, error) {
	const rootURl = "https://oauth2.googleapis.com/token"

	values := url.Values{}
	values.Add("grant_type", "authorization_code")
	values.Add("code", code)
	values.Add("client_id", config.AppConfig.GoogleClientID)
	values.Add("client_secret", config.AppConfig.GoogleClientSecret)
	values.Add("redirect_uri", config.AppConfig.GoogleOAuthRedirectUrl)

	query := values.Encode()

	req, err := http.NewRequestWithContext(ctx, "POST", rootURl, bytes.NewBufferString(query))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := http.Client{
		Timeout: time.Second * 30,
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		fmt.Print(config.AppConfig.GoogleClientID)
		return nil, errors.New("could not retrieve token")
	}

	var resBody bytes.Buffer
	_, err = io.Copy(&resBody, res.Body)
	if err != nil {
		return nil, err
	}

	var GoogleOauthTokenRes map[string]interface{}

	if err := json.Unmarshal(resBody.Bytes(), &GoogleOauthTokenRes); err != nil {
		return nil, err
	}

	tokenBody := &GoogleOauthToken{
		Access_token: GoogleOauthTokenRes["access_token"].(string),
		Id_token:     GoogleOauthTokenRes["id_token"].(string),
	}

	return tokenBody, nil
}

func GetGoogleUser(ctx context.Context, access_token string, id_token string) (*GoogleUserResult, error) {
	rootUrl := fmt.Sprintf("https://www.googleapis.com/oauth2/v2/userinfo?alt=json&access_token=%s", access_token)

	req, err := http.NewRequestWithContext(ctx, "GET", rootUrl, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", id_token))

	client := http.Client{
		Timeout: time.Second * 30,
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("could not retrieve user")
	}

	var resBody bytes.Buffer
	_, err = io.Copy(&resBody, res.Body)
	if err != nil {
		return nil, err
	}

	var GoogleUserRes map[string]interface{}

	if err := json.Unmarshal(resBody.Bytes(), &GoogleUserRes); err != nil {
		return nil, err
	}

	userBody := &GoogleUserResult{
		Id:            GoogleUserRes["id"].(string),
		Email:         GoogleUserRes["email"].(string),
		VerifiedEmail: GoogleUserRes["verified_email"].(bool),
		Name:          GoogleUserRes["name"].(string),
		GivenName:     GoogleUserRes["given_name"].(string),
		Picture:       GoogleUserRes["picture"].(string),
	}

	familyName, ok := GoogleUserRes["family_name"].(string)
	if ok {
		userBody.FamilyName = familyName
	}

	return userBody, nil
}

func RandomUsername(username string) string {
	number := rand.Intn(99999-10000) + 10000
	return fmt.Sprintf("%s%d", username, number)
}
