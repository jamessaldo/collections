package util

import (
	"authorization/config"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

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

var googleOauthConfig = &oauth2.Config{
	RedirectURL:  config.AppConfig.GoogleOAuthRedirectUrl,
	ClientID:     config.AppConfig.GoogleClientID,
	ClientSecret: config.AppConfig.GoogleClientSecret,
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
	Endpoint:     google.Endpoint,
}

const oauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v2/userinfo?alt=json&access_token="

func GetGoogleUser(code string) (*GoogleUserResult, error) {
	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("code exchange wrong: %s", err.Error())
	}
	response, err := http.Get(oauthGoogleUrlAPI + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed read response: %s", err.Error())
	}

	var GoogleUserRes map[string]interface{}

	if err := json.Unmarshal(contents, &GoogleUserRes); err != nil {
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
