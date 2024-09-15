package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/andreevym/gophkeeper/internal/client"
	"github.com/andreevym/gophkeeper/internal/storage"
)

const (
	successColor = "\033[32m" // Green
	errorColor   = "\033[31m" // Red
	resetColor   = "\033[0m"  // Reset
)

// Invoker defines the methods that our Client should implement.
type Invoker interface {
	CreateUser(login, password string) error
	SignIn(login, password string) (string, error)
	GetVault(token, vaultID string) (storage.Vault, error)
	NewVault(token, key, value, vaultID string) (storage.Vault, error)
}

func main() {
	if len(os.Args) < 2 {
		printHelp()
		os.Exit(1)
	}

	cmd := strings.ToLower(os.Args[1])
	if cmd == "help" {
		printHelp()
		os.Exit(0)
	}

	if len(os.Args) < 3 {
		fmt.Printf("%sError: Not enough arguments. Use 'help' command for usage.%s\n", errorColor, resetColor)
		os.Exit(1)
	}

	serverAddress := os.Args[2]
	c := client.NewClient(serverAddress)

	args := os.Args[2:]

	switch cmd {
	case "signup":
		handleSignUp(c, args)
	case "signin":
		handleSignIn(c, args)
	case "savevault":
		handleSaveVault(c, args)
	case "getvault":
		handleGetVault(c, args)
	default:
		fmt.Printf("%sError: Unknown command '%s'. Use 'help' command for usage.%s\n", errorColor, cmd, resetColor)
	}
}

func handleSignUp(invoker Invoker, args []string) {
	if len(args) < 3 {
		fmt.Printf("%sError: Sign-up command requires username and password.%s\n", errorColor, resetColor)
		os.Exit(1)
	}
	login, password := args[1], args[2]
	if err := invoker.CreateUser(login, password); err != nil {
		fmt.Printf("%sError: Failed to create user: %s%s\n", errorColor, err, resetColor)
		os.Exit(1)
	}
	fmt.Printf("%sUser created successfully%s\n", successColor, resetColor)
}

func handleSignIn(invoker Invoker, args []string) {
	if len(args) < 3 {
		fmt.Printf("%sError: Sign-in command requires username and password.%s\n", errorColor, resetColor)
		os.Exit(1)
	}
	login, password := args[1], args[2]
	token, err := invoker.SignIn(login, password)
	if err != nil {
		fmt.Printf("%sError: Sign-in failed: %s%s\n", errorColor, err, resetColor)
		os.Exit(1)
	}
	fmt.Println(token)
}

func handleSaveVault(invoker Invoker, args []string) {
	if len(args) < 4 {
		fmt.Printf("%sError: Save vault command requires token, key, and value.%s\n", errorColor, resetColor)
		os.Exit(1)
	}
	token, key, value := args[1], args[2], args[3]
	vaultID := ""
	if len(args) == 5 {
		vaultID = args[4]
	}
	vault, err := invoker.NewVault(token, key, value, vaultID)
	if err != nil {
		fmt.Printf("%sError: Failed to create vault: %s%s\n", errorColor, err, resetColor)
		os.Exit(1)
	}
	printVault(vault)
}

func handleGetVault(invoker Invoker, args []string) {
	if len(args) < 3 {
		fmt.Printf("%sError: Get vault command requires token and vault ID.%s\n", errorColor, resetColor)
		os.Exit(1)
	}
	token, vaultID := args[1], args[2]
	vault, err := invoker.GetVault(token, vaultID)
	if err != nil {
		fmt.Printf("%sError: Failed to get vault: %s%s\n", errorColor, err, resetColor)
		os.Exit(1)
	}
	printVault(vault)
}

func printVault(v storage.Vault) {
	b, err := json.Marshal(v)
	if err != nil {
		fmt.Printf("%sError: Failed to marshal vault response: %s%s\n", errorColor, err, resetColor)
		os.Exit(1)
	}
	fmt.Printf("%sVault operation successful: %s%s\n", successColor, string(b), resetColor)
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
