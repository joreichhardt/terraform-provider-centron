package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Client struct {
	BaseURL      string
	ClientID     string
	ClientSecret string
	HTTPClient   *http.Client

	token   *TokenResponse
	tokenMu sync.RWMutex
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	tokenGetTime time.Time
}

func NewClient(baseURL, clientID, clientSecret string) *Client {
	if baseURL == "" {
		baseURL = "https://ccenter.centron.de/api/v1"
	}
	return &Client{
		BaseURL:      strings.TrimSuffix(baseURL, "/"),
		ClientID:     clientID,
		ClientSecret: clientSecret,
		HTTPClient:   &http.Client{Timeout: 60 * time.Second},
	}
}

// GetToken handles the OAuth2 client credentials flow
func (c *Client) GetToken(ctx context.Context) (string, error) {
	c.tokenMu.RLock()
	// Check if token exists and is valid (with a 60-second buffer)
	if c.token != nil && time.Since(c.token.tokenGetTime) < time.Duration(c.token.ExpiresIn-60)*time.Second {
		tok := c.token.AccessToken
		c.tokenMu.RUnlock()
		return tok, nil
	}
	c.tokenMu.RUnlock()

	c.tokenMu.Lock()
	defer c.tokenMu.Unlock()

	// Check again in case another goroutine just updated it
	if c.token != nil && time.Since(c.token.tokenGetTime) < time.Duration(c.token.ExpiresIn-60)*time.Second {
		return c.token.AccessToken, nil
	}

	authURL := c.BaseURL + "/oauth/token"
	payload := map[string]string{
		"grant_type":    "client_credentials",
		"client_id":     c.ClientID,
		"client_secret": c.ClientSecret,
		"scope":         "*",
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", authURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to get token, status: %d, body: %s", resp.StatusCode, string(body))
	}

	var tr TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tr); err != nil {
		return "", err
	}
	tr.tokenGetTime = time.Now()
	c.token = &tr

	return c.token.AccessToken, nil
}

func (c *Client) DoRequest(ctx context.Context, method, path string, body interface{}, result interface{}) error {
	token, err := c.GetToken(ctx)
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	url := c.BaseURL + path
	var reqBody io.Reader

	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return err
		}
		reqBody = bytes.NewBuffer(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")
	if reqBody != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	if result != nil && resp.StatusCode != http.StatusNoContent {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return err
		}
	}

	return nil
}
