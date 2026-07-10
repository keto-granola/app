package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"time"
)

const (
	firebaseAuthURL = "https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword"
	ctxTimeout      = 10 * time.Second
	httpTimeout     = 10 * time.Second
)

type signInRequest struct {
	Email             string `json:"email"`
	Password          string `json:"password"`
	ReturnSecureToken bool   `json:"returnSecureToken"`
}

type signInResponse struct {
	IDToken      string `json:"idToken"`
	Email        string `json:"email"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    string `json:"expiresIn"`
	LocalID      string `json:"localId"`
}

type errorResponse struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Errors  []struct {
			Message string `json:"message"`
			Domain  string `json:"domain"`
			Reason  string `json:"reason"`
		} `json:"errors"`
	} `json:"error"`
}

func main() {
	var (
		apiKey   string
		email    string
		password string
		verbose  bool
	)

	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()

	flag.StringVar(&apiKey, "api-key", "", "Firebase Web API Key (required)")
	flag.StringVar(&email, "email", "", "User email address (required)")
	flag.StringVar(&password, "password", "", "User password (required)")
	flag.BoolVar(&verbose, "verbose", false, "Show verbose output including token details")
	flag.Parse()

	if err := run(ctx, apiKey, email, password, verbose); err != nil {
		slog.Error("run failed", slog.Any("error", err))
		os.Exit(1)
	}
}

func run(ctx context.Context, apiKey, email, password string, verbose bool) error {
	if apiKey == "" || email == "" || password == "" {
		flag.Usage()
		return fmt.Errorf("api-key, email and password are required")
	}

	token, err := getAccessToken(ctx, apiKey, email, password)
	if err != nil {
		return fmt.Errorf("get access token: %w", err)
	}

	if verbose {
		fmt.Printf("ID Token: %s\n", token.IDToken)
		fmt.Printf("Email: %s\n", token.Email)
		fmt.Printf("UserID: %s\n", token.LocalID)
		fmt.Printf("Expires In: %s seconds\n", token.ExpiresIn)
		fmt.Printf("\nRefresh Token: %s\n", token.RefreshToken)
	} else {
		fmt.Println(token.IDToken)
	}

	return nil
}

func getAccessToken(ctx context.Context, apiKey, email, password string) (*signInResponse, error) {
	reqBody := signInRequest{
		Email:             email,
		Password:          password,
		ReturnSecureToken: true,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	parsedURL, err := url.Parse(firebaseAuthURL)
	if err != nil {
		return nil, fmt.Errorf("parse URL: %w", err)
	}
	q := parsedURL.Query()
	q.Set("key", apiKey)
	parsedURL.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, parsedURL.String(), bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: httpTimeout,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var errResp errorResponse
		if err := json.Unmarshal(body, &errResp); err != nil {
			return nil, fmt.Errorf("authentication failed with status %d: %s", resp.StatusCode, string(body))
		}
		return nil, fmt.Errorf("authentication failed: %s", errResp.Error.Message)
	}

	var signInResp signInResponse
	if err := json.Unmarshal(body, &signInResp); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	return &signInResp, nil
}
