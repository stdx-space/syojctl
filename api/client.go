package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"time"

	"github.com/charmbracelet/log"
)

// Client represents the SYOJ API client
type Client struct {
	httpClient *http.Client
	baseURL    string
	logger     *log.Logger
	token      string
	tokenId    string
}

// Credentials represents the authentication credentials
type Credentials struct {
	Token   string `json:"token"`
	TokenId string `json:"token_id"`
}

// LoginRequest represents the login request payload
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse represents the login response
type LoginResponse struct {
	Message string `json:"message"`
}

// Problem represents a problem from SYOJ
type Problem struct {
	ID              string        `json:"id"`
	Title           string        `json:"title"`
	Description     string        `json:"description"`
	Difficulty      string        `json:"difficulty"`
	MemoryLimit     int           `json:"memoryLimit"`
	TimeLimit       int           `json:"timeLimit"`
	Notes           string        `json:"notes"`
	Tags            []Tag         `json:"tags"`
	AllowedLanguages []Language   `json:"allowedLanguages"`
	Author          Author        `json:"author"`
	TestCases       []TestCaseGroup `json:"testCases"`
	ProblemSection  []Section     `json:"ProblemSection"`
	AllowSubmit     bool          `json:"allowSubmit"`
}

// Tag represents a problem tag
type Tag struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Language represents an allowed programming language
type Language struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Value    string `json:"value"`
	Highlight string `json:"highlight"`
}

// Author represents the problem author
type Author struct {
	DisplayName string `json:"displayName"`
}

// TestCaseGroup represents a group of test cases
type TestCaseGroup struct {
	TestCases []TestCase `json:"testCases"`
}

// TestCase represents a single test case
type TestCase struct {
	Input  string `json:"input"`
	Output string `json:"output"`
}

// Section represents a problem section (Input, Output, Constraints, etc.)
type Section struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Order   int    `json:"order"`
}

// NewClient creates a new SYOJ API client
func NewClient() (*Client, error) {
	// Create a cookie jar to store cookies
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	// Create an HTTP client with the cookie jar
	httpClient := &http.Client{
		Jar:     jar,
		Timeout: 30 * time.Second,
	}

	// Create a logger with options
	logger := log.NewWithOptions(os.Stderr, log.Options{
		ReportTimestamp: true,
		ReportCaller:    true,
	})

	return &Client{
		httpClient: httpClient,
		baseURL:    "https://syoj.org",
		logger:     logger,
	}, nil
}

// NewClientWithCredentials creates a new SYOJ API client with authentication credentials
func NewClientWithCredentials(token, tokenId string) (*Client, error) {
	client, err := NewClient()
	if err != nil {
		return nil, err
	}

	client.token = token
	client.tokenId = tokenId

	// Set the cookies for authentication
	u, _ := url.Parse(client.baseURL)
	cookies := []*http.Cookie{
		{
			Name:  "Token",
			Value: token,
		},
		{
			Name:  "TokenId",
			Value: tokenId,
		},
	}
	client.httpClient.Jar.SetCookies(u, cookies)

	return client, nil
}

// Login performs login to SYOJ and returns credentials
func (c *Client) Login(email, password string) (*Credentials, error) {
	// Prepare the login request
	loginData := LoginRequest{
		Email:    email,
		Password: password,
	}

	jsonData, err := json.Marshal(loginData)
	if err != nil {
		c.logger.Error("Failed to marshal login data", "error", err)
		return nil, err
	}

	// Create the request
	req, err := http.NewRequest("POST", c.baseURL+"/api/login", bytes.NewBuffer(jsonData))
	if err != nil {
		c.logger.Error("Failed to create request", "error", err)
		return nil, err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "syojctl/1.0")

	// Perform the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.logger.Error("Failed to make request", "error", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.logger.Error("Failed to read response body", "error", err)
		return nil, err
	}

	// Log response details
	c.logger.Info("Received response", "status", resp.StatusCode, "body", string(body))

	// Parse the response
	var loginResp LoginResponse
	if err := json.Unmarshal(body, &loginResp); err != nil {
		c.logger.Error("Failed to parse response", "error", err)
		return nil, err
	}

	// Check if login was successful
	if resp.StatusCode == http.StatusOK && loginResp.Message == "Successfully logged in" {
		c.logger.Info("Login successful!")

		// Extract cookies
		u, _ := url.Parse(c.baseURL)
		cookies := c.httpClient.Jar.Cookies(u)

		// Create credentials object
		creds := Credentials{}

		c.logger.Debug("Extracting cookies")
		for _, cookie := range cookies {
			c.logger.Debug("Cookie", "name", cookie.Name, "value", cookie.Value)
			// Save the important cookies
			if cookie.Name == "Token" {
				creds.Token = cookie.Value
			} else if cookie.Name == "TokenId" {
				creds.TokenId = cookie.Value
			}
		}

		return &creds, nil
	} else {
		c.logger.Error("Login failed", "status", resp.StatusCode, "message", loginResp.Message)
		return nil, fmt.Errorf("login failed with status %d: %s", resp.StatusCode, loginResp.Message)
	}
}

// GetProblem fetches a problem by its ID
func (c *Client) GetProblem(problemID string) (*Problem, error) {
	// Create the request
	req, err := http.NewRequest("GET", c.baseURL+"/api/problems/"+problemID, nil)
	if err != nil {
		c.logger.Error("Failed to create request", "error", err)
		return nil, err
	}

	// Set headers
	req.Header.Set("User-Agent", "syojctl/1.0")

	// Perform the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.logger.Error("Failed to make request", "error", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Check if we got a successful response
	if resp.StatusCode != http.StatusOK {
		c.logger.Error("Failed to fetch problem", "status", resp.StatusCode)
		return nil, fmt.Errorf("failed to fetch problem with status %d", resp.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.logger.Error("Failed to read response body", "error", err)
		return nil, err
	}

	// Log response details
	c.logger.Info("Received problem response", "status", resp.StatusCode, "body_length", len(body))

	// Parse the response
	var problem Problem
	if err := json.Unmarshal(body, &problem); err != nil {
		c.logger.Error("Failed to parse problem response", "error", err)
		return nil, err
	}

	return &problem, nil
}