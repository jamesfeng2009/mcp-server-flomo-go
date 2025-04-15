package flomo

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

// Client represents a client for interacting with the Flomo API
type Client struct {
	apiURL string
	logger *log.Logger
}

// Response represents the response from Flomo API
type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Memo    struct {
		Slug      string   `json:"slug"`
		CreatorID int      `json:"creator_id"`
		Source    string   `json:"source"`
		Content   string   `json:"content"`
		Tags      []string `json:"tags"`
		UpdatedAt string   `json:"updated_at"`
		CreatedAt string   `json:"created_at"`
	} `json:"memo"`
}

// NewClient creates a new Flomo client
func NewClient(apiURL string) *Client {
	return &Client{
		apiURL: apiURL,
		logger: log.New(log.Writer(), "[Flomo] ", log.LstdFlags|log.Lmsgprefix),
	}
}

// WriteNote writes a note to Flomo
func (c *Client) WriteNote(content string) (*Response, error) {
	startTime := time.Now()
	c.logger.Printf("Starting to send note (length: %d characters)", len(content))

	if content == "" {
		c.logger.Println("Error: Empty content provided")
		return nil, errors.New("content cannot be empty")
	}

	// Prepare request body
	reqBody := map[string]string{
		"content": content,
	}
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		c.logger.Printf("Error: Failed to marshal request body: %v", err)
		return nil, fmt.Errorf("failed to marshal request body: %v", err)
	}
	c.logger.Printf("Request body prepared (size: %d bytes)", len(jsonBody))

	// Create request
	req, err := http.NewRequest("POST", c.apiURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		c.logger.Printf("Error: Failed to create request: %v", err)
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "MCP-Server-Flomo-Go/1.0")
	c.logger.Println("Request headers set")

	// Send request
	c.logger.Printf("Sending request to %s", c.apiURL)
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		c.logger.Printf("Error: Failed to send request: %v", err)
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()
	c.logger.Printf("Response received (status: %s)", resp.Status)

	// Read response
	var flomoResp Response
	if err := json.NewDecoder(resp.Body).Decode(&flomoResp); err != nil {
		c.logger.Printf("Error: Failed to decode response: %v", err)
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		c.logger.Printf("Error: Request failed with status %d: %s", resp.StatusCode, flomoResp.Message)
		return &flomoResp, fmt.Errorf("request failed with status %d: %s", resp.StatusCode, flomoResp.Message)
	}

	duration := time.Since(startTime)
	c.logger.Printf("Note sent successfully (took %v)", duration)
	c.logger.Printf("Memo details - CreatedAt: %s, Tags: %v", flomoResp.Memo.CreatedAt, flomoResp.Memo.Tags)

	return &flomoResp, nil
} 