// Package client provides a client for interacting with the GophKeeper service.
// It allows for user registration, authentication, and vault management.
package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/andreevym/gophkeeper/internal/handlers"
)

// Client represents a client that communicates with the GophKeeper service.
type Client struct {
	serverAddress string
}

// NewClient creates a new instance of Client with the provided server address.
func NewClient(serverAddress string) *Client {
	return &Client{serverAddress: serverAddress}
}

// CreateUser registers a new user with the GophKeeper service.
// It sends a POST request to the /signup endpoint with the provided login and password.
// Returns an error if the request fails or if the server responds with a non-200 status code.
func (c *Client) CreateUser(login, password string) error {
	b, err := json.Marshal(handlers.SignUpRequest{Login: login, Password: password})
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}
	resp, err := http.Post(c.serverAddress+handlers.AuthSignUpURI, "application/json", bytes.NewBuffer(b))
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return handleErrorResponse(resp)
	}
	return nil
}

// SignIn authenticates a user and retrieves an authentication token.
// It sends a POST request to the /signin endpoint with the provided login and password.
// Returns the token if successful, or an error if the request fails or if the server responds with a non-200 status code.
func (c *Client) SignIn(login, password string) (string, error) {
	b, err := json.Marshal(handlers.SignInRequest{Login: login, Password: password})
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}
	resp, err := http.Post(c.serverAddress+handlers.AuthSignInURI, "application/json", bytes.NewBuffer(b))
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return "", handleErrorResponse(resp)
	}

	token := resp.Header.Get("Authorization")
	return strings.TrimPrefix(token, "Bearer "), nil
}

// GetVault retrieves a vault by its ID using the provided authentication token.
// It sends a GET request to the /vault/{vaultID} endpoint with the token in the Authorization header.
// Returns the vault if successful, or an error if the request fails, the server responds with a non-200 status code, or if the vaultID is empty.
func (c *Client) GetVault(token, vaultID string) (handlers.VaultResponse, error) {
	if vaultID == "" {
		return handlers.VaultResponse{}, errors.New("vaultID is empty")
	}
	u, err := url.Parse(fmt.Sprintf("%s%s/%s", c.serverAddress, handlers.VaultURI, vaultID))
	if err != nil {
		return handlers.VaultResponse{}, fmt.Errorf("failed to parse URL: %w", err)
	}
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return handlers.VaultResponse{}, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return handlers.VaultResponse{}, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return handlers.VaultResponse{}, handleErrorResponse(resp)
	}

	var vault handlers.VaultResponse
	if err := json.NewDecoder(resp.Body).Decode(&vault); err != nil {
		return handlers.VaultResponse{}, fmt.Errorf("failed to decode response: %w", err)
	}
	return vault, nil
}

// NewVault creates a new vault with the provided key, value, and optional vaultID using the provided authentication token.
// It sends a POST request to the /vault endpoint with the token in the Authorization header.
// Returns the created vault if successful, or an error if the request fails or if the server responds with a non-200 status code.
func (c *Client) NewVault(token, key, value, vaultID string) (handlers.VaultResponse, error) {
	vaultRequest := handlers.VaultRequest{Key: key, Value: value, ID: vaultID}
	b, err := json.Marshal(vaultRequest)
	if err != nil {
		return handlers.VaultResponse{}, fmt.Errorf("failed to marshal request: %w", err)
	}
	req, err := http.NewRequest(http.MethodPost, c.serverAddress+handlers.VaultURI, bytes.NewBuffer(b))
	if err != nil {
		return handlers.VaultResponse{}, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return handlers.VaultResponse{}, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return handlers.VaultResponse{}, handleErrorResponse(resp)
	}

	var vault handlers.VaultResponse
	if err := json.NewDecoder(resp.Body).Decode(&vault); err != nil {
		return handlers.VaultResponse{}, fmt.Errorf("failed to decode response: %w", err)
	}
	return vault, nil
}

// UploadFile uploads a file to the server and handles both creating a new vault entry and updating an existing one.
// The `vaultID` parameter is optional: if provided, the file will be uploaded to the specified vault.
func (c *Client) UploadFile(token, filename, filePath, vaultID string) (handlers.VaultResponse, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Create a form file part
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return handlers.VaultResponse{}, fmt.Errorf("failed to create form file: %w", err)
	}

	// Open the file to be uploaded
	file, err := os.Open(filePath)
	if err != nil {
		return handlers.VaultResponse{}, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Copy the file contents to the form file part
	_, err = io.Copy(part, file)
	if err != nil {
		return handlers.VaultResponse{}, fmt.Errorf("failed to copy file contents: %w", err)
	}

	// Close the multipart writer to set the terminating boundary
	err = writer.Close()
	if err != nil {
		return handlers.VaultResponse{}, fmt.Errorf("failed to close writer: %w", err)
	}

	// Create a new POST request
	url := c.serverAddress + handlers.FileUploadURI
	if vaultID != "" {
		url += "/" + vaultID
	}
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return handlers.VaultResponse{}, fmt.Errorf("failed to create request: %w", err)
	}

	// Set the necessary headers
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+token)

	// Send the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return handlers.VaultResponse{}, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Check for successful status codes
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return handlers.VaultResponse{}, handleErrorResponse(resp)
	}

	// Decode the response
	var vault handlers.VaultResponse
	if err := json.NewDecoder(resp.Body).Decode(&vault); err != nil {
		return handlers.VaultResponse{}, fmt.Errorf("failed to decode response: %w", err)
	}

	return vault, nil
}

// handleErrorResponse handles HTTP errors by reading the response body and returning a formatted error message.
// It is used internally by the client methods to provide detailed error information when an HTTP request fails.
func handleErrorResponse(resp *http.Response) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}
	return fmt.Errorf("HTTP error: %s, body: %s", resp.Status, string(body))
}
