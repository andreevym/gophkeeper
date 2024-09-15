package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/andreevym/gophkeeper/internal/handlers"
	"github.com/andreevym/gophkeeper/internal/storage"
	"github.com/andreevym/gophkeeper/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("received worn count of arguments: %d, byt expected more than one\n", len(os.Args)-1)
		printHelp()
		os.Exit(1)
	}

	cmd := os.Args[1]
	if cmd == "help" {
		printHelp()
		os.Exit(0)
	}

	if len(os.Args) < 3 {
		fmt.Printf("received worn count of arguments: %d, byt expected more than one\n", len(os.Args)-1)
		printHelp()
		os.Exit(1)
	}

	serverAddress := os.Args[2]
	client := NewClient(serverAddress)

	args := os.Args[2:]

	switch cmd {
	case "signUp":
		if len(args) < 2 {
			logger.Logger().Error("user command requires at least two arguments")
			os.Exit(1)
		}
		login := args[1]
		password := args[2]
		err := client.CreateUser(login, password)
		if err != nil {
			logger.Logger().Error("failed to create user", zap.String("login", login), zap.Error(err))
			os.Exit(1)
		}
		logger.Logger().Info("user command executed successfully")
		os.Exit(0)
	case "signIn":
		if len(args) < 2 {
			logger.Logger().Error("user command requires at least two arguments")
			os.Exit(1)
		}
		login := args[1]
		password := args[2]
		token, err := client.SignIn(login, password)
		if err != nil {
			logger.Logger().Error("no command to execute", zap.Error(err))
			os.Exit(1)
		}
		fmt.Println(token)
		os.Exit(0)
	case "saveVault":
		if len(args) < 2 {
			logger.Logger().Error("vault command requires at least two arguments")
			os.Exit(1)
		}
		token := args[1]
		key := args[2]
		value := args[3]
		vaultID := ""
		if len(args) == 5 {
			vaultID = args[4]
		}
		v, err := client.NewVault(token, key, value, vaultID)
		if err != nil {
			logger.Logger().Error("failed to create vault", zap.Error(err))
			os.Exit(1)
		}
		b, err := json.Marshal(v)
		if err != nil {
			logger.Logger().Error("marshal vault command executed successfully", zap.Error(err))
			os.Exit(1)
		}
		logger.Logger().Info("vault command executed successfully", zap.String("resp", string(b)))
	case "getVault":
		if len(args) < 2 {
			logger.Logger().Error("vault command requires at least two arguments")
			os.Exit(1)
		}
		token := args[1]
		vaultID := args[2]
		v, err := client.GetVault(token, vaultID)
		if err != nil {
			logger.Logger().Error("no vault command executed successfully", zap.Error(err))
			os.Exit(1)
		}
		b, err := json.Marshal(v)
		if err != nil {
			logger.Logger().Error("marshal vault command executed successfully", zap.Error(err))
			os.Exit(1)
		}
		logger.Logger().Info("vault command executed successfully", zap.String("resp", string(b)))
	}

}

type Client struct {
	serverAddress string
}

func NewClient(serverAddress string) *Client {
	return &Client{serverAddress: serverAddress}
}

func (c Client) CreateUser(login string, password string) error {
	b, err := json.Marshal(handlers.SignUpRequest{Login: login, Password: password})
	if err != nil {
		return fmt.Errorf("failed marshal: %w", err)
	}
	resp, err := http.Post(c.serverAddress+handlers.AuthSignUpURI, "application/json", bytes.NewBuffer(b))
	if err != nil {
		return fmt.Errorf("failed post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		readAll, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed read response: %w", err)
		}
		return fmt.Errorf("failed post: %s, body: %s", resp.Status, string(readAll))
	}
	return nil
}

func (c Client) SignIn(login string, password string) (string, error) {
	b, err := json.Marshal(handlers.SignInRequest{Login: login, Password: password})
	if err != nil {
		return "", fmt.Errorf("failed marshal: %w", err)
	}
	resp, err := http.Post(c.serverAddress+handlers.AuthSignInURI, "application/json", bytes.NewBuffer(b))
	if err != nil {
		return "", fmt.Errorf("failed post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		readAll, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", fmt.Errorf("failed read response: %w", err)
		}
		return "", fmt.Errorf("failed post: %s, body %s", resp.Status, string(readAll))
	}

	token := resp.Header.Get("Authorization")
	token = strings.Replace(token, "Bearer ", "", 1)
	return token, nil
}

func (c Client) GetVault(token string, vaultID string) (storage.Vault, error) {
	if vaultID == "" {
		return storage.Vault{}, errors.New("vaultID is empty")
	}
	u, err := url.Parse(c.serverAddress + handlers.VaultURI + "/" + vaultID)
	if err != nil {
		return storage.Vault{}, fmt.Errorf("failed parse url: %w", err)
	}
	request := &http.Request{
		Method: http.MethodGet,
		URL:    u,
		Header: http.Header{
			"Authorization": []string{fmt.Sprintf("Bearer %s", token)},
		},
	}
	httpClient := &http.Client{}
	resp, err := httpClient.Do(request)
	if err != nil {
		return storage.Vault{}, fmt.Errorf("failed post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		readAll, err := io.ReadAll(resp.Body)
		if err != nil {
			return storage.Vault{}, fmt.Errorf("failed read response: %w", err)
		}
		return storage.Vault{}, fmt.Errorf("failed post: %s, body %s", resp.Status, string(readAll))
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return storage.Vault{}, fmt.Errorf("failed read all: %w", err)
	}

	vaultResponse := storage.Vault{}
	err = json.Unmarshal(b, &vaultResponse)
	if err != nil {
		return storage.Vault{}, fmt.Errorf("failed unmarshal: %w", err)
	}
	return vaultResponse, nil
}

func (c Client) NewVault(token string, key string, value string, vaultID string) (storage.Vault, error) {
	vaultRequest := handlers.VaultRequest{
		Key:   key,
		Value: value,
	}
	if vaultID != "" {
		vaultRequest.ID = vaultID
	}
	b, err := json.Marshal(vaultRequest)
	if err != nil {
		return storage.Vault{}, fmt.Errorf("failed marshal: %w", err)
	}
	httpClient := &http.Client{}
	request, err := http.NewRequest(http.MethodPost, c.serverAddress+handlers.VaultURI, bytes.NewBuffer(b))
	if err != nil {
		return storage.Vault{}, fmt.Errorf("failed post: %w", err)
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+token)
	resp, err := httpClient.Do(request)
	if err != nil {
		return storage.Vault{}, fmt.Errorf("failed post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		readAll, err := io.ReadAll(resp.Body)
		if err != nil {
			return storage.Vault{}, fmt.Errorf("failed read response: %w", err)
		}
		return storage.Vault{}, fmt.Errorf("failed post: %s, body %s", resp.Status, string(readAll))
	}

	b, err = io.ReadAll(resp.Body)
	if err != nil {
		return storage.Vault{}, fmt.Errorf("failed read all: %w", err)
	}

	vaultResponse := storage.Vault{}
	err = json.Unmarshal(b, &vaultResponse)
	if err != nil {
		return storage.Vault{}, fmt.Errorf("failed unmarshal: %w", err)
	}
	return vaultResponse, nil
}

func printHelp() {
	fmt.Println("GophKeeper CLI Help")
	fmt.Println("---------------------")
	fmt.Println("This is a CLI tool to interact with the GophKeeper service. Below are the available commands:")
	fmt.Println()
	fmt.Println("1. Sign Up")
	fmt.Println("Description: Register a new user.")
	fmt.Println("Usage: ./client signUp <server_url> <username> <password>")
	fmt.Println("Example:")
	fmt.Println("  ./client signUp http://localhost:8080 testName testPassword")
	fmt.Println()
	fmt.Println("2. Sign In")
	fmt.Println("Description: Log in a user and retrieve an authentication token.")
	fmt.Println("Usage: ./client signIn <server_url> <username> <password>")
	fmt.Println("Example:")
	fmt.Println("  ./client signIn http://localhost:8080 testName testPassword")
	fmt.Println("Response:")
	fmt.Println("  eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySUQiOiIxIn0.bMFEsrvtCxd5i3SMn3E_8HcRx6RzNfTX2PI1eWXJsbNUbeG_VaEpf9trTcm4KsYqYp_wpLzMYEYKQCtQykb4lQ")
	fmt.Println()
	fmt.Println("3. Create Vault")
	fmt.Println("Description: Create a new vault using an authentication token.")
	fmt.Println("Usage: ./client saveVault <server_url> <token> <key> <value>")
	fmt.Println("Example:")
	fmt.Println("  ./client saveVault http://localhost:8080 <token> k1 v1")
	fmt.Println()
	fmt.Println("4. Get Vault")
	fmt.Println("Description: Retrieve vault details using an authentication token and vault ID.")
	fmt.Println("Usage: ./client getVault <server_url> <token> <vault_id>")
	fmt.Println("Example:")
	fmt.Println("  ./client getVault http://localhost:8080 <token> 1")
	fmt.Println()
}
