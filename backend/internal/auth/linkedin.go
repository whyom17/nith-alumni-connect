package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	ErrLinkedInAuthFailed = errors.New("linkedin authentication failed")
	ErrInvalidLinkedInCode = errors.New("invalid or expired linkedin authorization code")
)

// LinkedInClient defines the interface for communicating with LinkedIn's OAuth/OIDC endpoints.
type LinkedInClient interface {
	ExchangeCode(ctx context.Context, code, codeVerifier, redirectURI string) (*LinkedInUserInfo, error)
}

// HTTPLinkedInClient is the production implementation contacting LinkedIn servers.
type HTTPLinkedInClient struct {
	clientID     string
	clientSecret string
	httpClient   *http.Client
}

// NewLinkedInClient creates a new HTTPLinkedInClient.
func NewLinkedInClient(clientID, clientSecret string) *HTTPLinkedInClient {
	return &HTTPLinkedInClient{
		clientID:     clientID,
		clientSecret: clientSecret,
		httpClient:   &http.Client{Timeout: 10 * time.Second},
	}
}

type linkedInTokenResp struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	IDToken     string `json:"id_token"`
	Scope       string `json:"scope"`
	Error       string `json:"error"`
	ErrorDesc   string `json:"error_description"`
}

type linkedInUserInfoResp struct {
	Sub     string `json:"sub"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Picture string `json:"picture"`
}

func (c *HTTPLinkedInClient) ExchangeCode(ctx context.Context, code, codeVerifier, redirectURI string) (*LinkedInUserInfo, error) {
	if c.clientID == "" || c.clientSecret == "" {
		return nil, fmt.Errorf("%w: missing linkedin client credentials in environment", ErrLinkedInAuthFailed)
	}

	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", redirectURI)
	data.Set("client_id", c.clientID)
	data.Set("client_secret", c.clientSecret)
	if codeVerifier != "" {
		data.Set("code_verifier", codeVerifier)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://www.linkedin.com/oauth/v2/accessToken", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("creating token request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to connect to linkedin: %v", ErrLinkedInAuthFailed, err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var tokenErr linkedInTokenResp
		_ = json.Unmarshal(bodyBytes, &tokenErr)
		return nil, fmt.Errorf("%w: code exchange failed with status %d: %s (%s)", ErrInvalidLinkedInCode, resp.StatusCode, tokenErr.Error, tokenErr.ErrorDesc)
	}

	var tokenResp linkedInTokenResp
	if err := json.Unmarshal(bodyBytes, &tokenResp); err != nil {
		return nil, fmt.Errorf("decoding token response: %w", err)
	}

	// Fetch userinfo using the obtained access token
	userInfoReq, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.linkedin.com/v2/userinfo", nil)
	if err != nil {
		return nil, fmt.Errorf("creating userinfo request: %w", err)
	}
	userInfoReq.Header.Set("Authorization", "Bearer "+tokenResp.AccessToken)

	userInfoResp, err := c.httpClient.Do(userInfoReq)
	if err != nil {
		return nil, fmt.Errorf("%w: userinfo request failed: %v", ErrLinkedInAuthFailed, err)
	}
	defer userInfoResp.Body.Close()

	if userInfoResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: userinfo returned status %d", ErrLinkedInAuthFailed, userInfoResp.StatusCode)
	}

	var info linkedInUserInfoResp
	if err := json.NewDecoder(userInfoResp.Body).Decode(&info); err != nil {
		return nil, fmt.Errorf("decoding userinfo: %w", err)
	}

	if info.Sub == "" || info.Email == "" {
		return nil, fmt.Errorf("%w: linkedin userinfo missing required sub or email", ErrLinkedInAuthFailed)
	}

	return &LinkedInUserInfo{
		Sub:     info.Sub,
		Email:   info.Email,
		Name:    info.Name,
		Picture: info.Picture,
	}, nil
}
