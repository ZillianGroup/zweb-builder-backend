package oauthgoogle

import (
	"encoding/json"

	"github.com/go-resty/resty/v2"
	"github.com/zilliangroup/zweb-builder-backend/src/utils/config"
)

const (
	GOOGLE_OAUTH2_API = "https://oauth2.googleapis.com/token"
)

type RefreshTokenResponse struct {
	AccessToken string `json:"access_token"`
	Expiry      int    `json:"expires_in"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}

type ExchangeTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Expiry       int    `json:"expires_in"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
}

func NewRefreshTokenResponse() *RefreshTokenResponse {
	return &RefreshTokenResponse{}
}

func NewExchangeTokenResponse() *ExchangeTokenResponse {
	return &ExchangeTokenResponse{}
}

func (resp *RefreshTokenResponse) ExportAccessToken() string {
	return resp.AccessToken
}

func (resp *ExchangeTokenResponse) ExportAccessToken() string {
	return resp.AccessToken
}

func ExchangeOAuthToken(code string) (*ExchangeTokenResponse, error) {
	conf := config.GetInstance()
	googleOAuthClientID := conf.GetZWebGoogleSheetsClientID()
	googleOAuthClientSecret := conf.GetZWebGoogleSheetsClientSecret()
	googleOAuthRedirectURI := conf.GetZWebGoogleSheetsRedirectURI()
	client := resty.New()
	// request
	resp, errInPost := client.R().
		SetFormData(map[string]string{
			"client_id":     googleOAuthClientID,
			"client_secret": googleOAuthClientSecret,
			"code":          code,
			"grant_type":    "authorization_code",
			"redirect_uri":  googleOAuthRedirectURI,
		}).
		Post(GOOGLE_OAUTH2_API)
	if resp.IsError() {
		return nil, errInPost
	}
	// unmarshal
	exchangeTokenResponse := NewExchangeTokenResponse()
	errInUnmarshal := json.Unmarshal(resp.Body(), &exchangeTokenResponse)
	if errInUnmarshal != nil {
		return nil, errInUnmarshal
	}
	return exchangeTokenResponse, nil
}

func RefreshOAuthToken(refreshToken string) (*RefreshTokenResponse, error) {
	conf := config.GetInstance()
	googleOAuthClientID := conf.GetZWebGoogleSheetsClientID()
	googleOAuthClientSecret := conf.GetZWebGoogleSheetsClientSecret()
	client := resty.New()
	// request
	resp, errInPost := client.R().
		SetFormData(map[string]string{
			"client_id":     googleOAuthClientID,
			"client_secret": googleOAuthClientSecret,
			"refresh_token": refreshToken,
			"grant_type":    "refresh_token",
		}).
		Post(GOOGLE_OAUTH2_API)

	if resp.IsError() {
		return nil, errInPost
	}
	// unmarshal
	refreshTokenResponse := NewRefreshTokenResponse()
	errInUnmarshal := json.Unmarshal(resp.Body(), &refreshTokenResponse)
	if errInUnmarshal != nil {
		return nil, errInUnmarshal
	}
	return refreshTokenResponse, nil
}
