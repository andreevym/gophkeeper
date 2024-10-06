package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/andreevym/gophkeeper/internal/client"
	"github.com/andreevym/gophkeeper/internal/handlers"
)

const (
	successColor = "\033[32m" // Green color for success messages
	errorColor   = "\033[31m" // Red color for error messages
	resetColor   = "\033[0m"  // Reset color to default
)

var gitRef string
var buildTime string
var gitCommit string

func printVersion() {
	if gitRef == "" {
		gitRef = "N/A"
	}
	if buildTime == "" {
		buildTime = "N/A"
	}
	if gitCommit == "" {
		gitCommit = "N/A"
	}

	fmt.Printf("GophKeeper Client\nVersion: %s\nBuild Date: %s\nBuild Commit: %s\n", gitRef, buildTime, gitCommit)
}

// Invoker defines the methods that our Client should implement.
type Invoker interface {
	CreateUser(login, password string) error
	SignIn(login, password string) (string, error)
	GetVault(token, vaultID string) (handlers.VaultResponse, error)
	NewVault(token, key, value, vaultID string) (handlers.VaultResponse, error)
	UploadFile(token, filename, filePath, vaultID string) (handlers.VaultResponse, error)
}

func main() {
	if len(os.Args) < 2 {
		printHelp()
		os.Exit(1)
	}

	cmd := strings.ToLower(os.Args[1])
	switch cmd {
	case "--version":
		printVersion()
		os.Exit(0)
	case "help":
		printHelp()
		os.Exit(0)
	}

	if len(os.Args) < 3 {
		fmt.Printf("%sError: Not enough arguments. Use 'help' command for usage.%s\n", errorColor, resetColor)
		os.Exit(1)
	}

	serverAddress := os.Args[2]
	c := client.NewClient(serverAddress)

	switch cmd {
	case "signup":
		handleSignUp(c, os.Args[3:])
	case "signin":
		handleSignIn(c, os.Args[3:])
	case "savevault":
		handleSaveVault(c, os.Args[3:])
	case "getvault":
		handleGetVault(c, os.Args[3:])
	case "uploadfile":
		handleUploadFile(c, os.Args[3:])
	case "store-login-password":
		handleStoreLoginPasswordVault(c, os.Args[3:])
	case "get-login-password":
		handleGetLoginPasswordVault(c, os.Args[3:])
	case "store-text":
		handleStoreTextVault(c, os.Args[3:])
	case "get-text":
		handleGetTextVault(c, os.Args[3:])
	case "store-binary":
		handleStoreBinaryVault(c, os.Args[3:])
	case "get-binary":
		handleGetBinaryVault(c, os.Args[3:])
	case "store-card":
		handleStoreCardVault(c, os.Args[3:])
	case "get-card":
		handleGetCardVault(c, os.Args[3:])
	default:
		fmt.Printf("%sError: Unknown command '%s'. Use 'help' command for usage.%s\n", errorColor, cmd, resetColor)
	}
}

// handleSignUp processes the sign-up command.
// It requires a username and password to create a new user.
func handleSignUp(invoker Invoker, args []string) {
	if len(args) < 2 {
		fmt.Printf("%sError: Sign-up command requires username and password.%s\n", errorColor, resetColor)
		os.Exit(1)
	}
	login, password := args[0], args[1]
	if err := invoker.CreateUser(login, password); err != nil {
		fmt.Printf("%sError: Failed to create user: %s%s\n", errorColor, err, resetColor)
		os.Exit(1)
	}
	fmt.Printf("%sUser created successfully%s\n", successColor, resetColor)
}

// handleSignIn processes the sign-in command.
// It requires a username and password to authenticate and retrieve a token.
func handleSignIn(invoker Invoker, args []string) {
	if len(args) < 2 {
		fmt.Printf("%sError: Sign-in command requires username and password.%s\n", errorColor, resetColor)
		os.Exit(1)
	}
	login, password := args[0], args[1]
	token, err := invoker.SignIn(login, password)
	if err != nil {
		fmt.Printf("%sError: Sign-in failed: %s%s\n", errorColor, err, resetColor)
		os.Exit(1)
	}
	fmt.Println(token)
}

// handleSaveVault processes the save vault command.
// It requires a token, key, and value to create a new vault entry. Optionally, a vault ID can be provided.
func handleSaveVault(invoker Invoker, args []string) {
	if len(args) < 3 {
		fmt.Printf("%sError: Save vault command requires token, key, and value.%s\n", errorColor, resetColor)
		os.Exit(1)
	}
	token, key, value := args[0], args[1], args[2]
	vaultID := ""
	if len(args) == 4 {
		vaultID = args[3]
	}
	vault, err := invoker.NewVault(token, key, value, vaultID)
	if err != nil {
		fmt.Printf("%sError: Failed to create vault: %s%s\n", errorColor, err, resetColor)
		os.Exit(1)
	}
	printVault(vault)
}

// handleGetVault processes the get vault command.
// It requires a token and vault ID to retrieve vault details.
func handleGetVault(invoker Invoker, args []string) {
	if len(args) < 2 {
		fmt.Printf("%sError: Get vault command requires token and vault ID.%s\n", errorColor, resetColor)
		os.Exit(1)
	}
	token, vaultID := args[0], args[1]
	vault, err := invoker.GetVault(token, vaultID)
	if err != nil {
		fmt.Printf("%sError: Failed to get vault: %s%s\n", errorColor, err, resetColor)
		os.Exit(1)
	}
	printVault(vault)
}

// handleUploadFile processes the upload file command.
// It requires a token, filename, file path, and optionally a vault ID for updating an existing entry.
func handleUploadFile(invoker Invoker, args []string) {
	if len(args) < 3 {
		fmt.Printf("%sError: Upload file command requires token, filename, and file path.%s\n", errorColor, resetColor)
		os.Exit(1)
	}
	token, filename, filePath := args[0], args[1], args[2]
	vaultID := ""
	if len(args) == 4 {
		vaultID = args[3]
	}

	vault, err := invoker.UploadFile(token, filename, filePath, vaultID)
	if err != nil {
		fmt.Printf("%sError: Failed to upload file: %s%s\n", errorColor, err, resetColor)
		os.Exit(1)
	}
	printVault(vault)
}

// printVault outputs the vault response in a JSON format.
func printVault(v handlers.VaultResponse) {
	b, err := json.Marshal(v)
	if err != nil {
		fmt.Printf("%sError: Failed to marshal vault response: %s%s\n", errorColor, err, resetColor)
		os.Exit(1)
	}
	fmt.Printf("%sVault operation successful: %s%s\n", successColor, string(b), resetColor)
}

// handleStoreLoginPasswordVault handles storing a login/password pair in the vault.
func handleStoreLoginPasswordVault(invoker Invoker, args []string) {
	if len(args) < 4 {
		fmt.Printf("%sError: Store login/password command requires token, vault ID, login, and password.%s\n", errorColor, resetColor)
		os.Exit(1)
	}
	token, vaultID, login, password := args[0], args[1], args[2], args[3]
	key := "login/" + login
	value := fmt.Sprintf(`{"login": "%s", "password": "%s"}`, login, password)

	vault, err := invoker.NewVault(token, key, value, vaultID)
	if err != nil {
		fmt.Printf("%sError: Failed to store login/password in vault: %s%s\n", errorColor, err, resetColor)
		os.Exit(1)
	}
	printVault(vault)
}

// handleStoreTextVault handles storing arbitrary text data in the vault.
func handleStoreTextVault(invoker Invoker, args []string) {
	if len(args) < 3 {
		fmt.Printf("%sError: Store text vault command requires token, vault ID, and text.%s\n", errorColor, resetColor)
		os.Exit(1)
	}
	token, vaultID, text := args[0], args[1], args[2]
	key := "text/" + vaultID
	value := fmt.Sprintf(`{"text": "%s"}`, text)

	vault, err := invoker.NewVault(token, key, value, vaultID)
	if err != nil {
		fmt.Printf("%sError: Failed to store text in vault: %s%s\n", errorColor, err, resetColor)
		os.Exit(1)
	}
	printVault(vault)
}

// handleGetLoginPasswordVault handles retrieving a login/password pair from the vault.
func handleGetLoginPasswordVault(invoker Invoker, args []string) {
	if len(args) < 2 {
		fmt.Printf("%sError: Get login/password command requires token and vault ID.%s\n", errorColor, resetColor)
		os.Exit(1)
	}
	token, vaultID := args[0], args[1]
	vault, err := invoker.GetVault(token, vaultID)
	if err != nil {
		fmt.Printf("%sError: Failed to retrieve login/password from vault: %s%s\n", errorColor, err, resetColor)
		os.Exit(1)
	}
	printVault(vault)
}

// handleGetTextVault handles retrieving text data from the vault.
func handleGetTextVault(invoker Invoker, args []string) {
	if len(args) < 2 {
		fmt.Printf("%sError: Get text vault command requires token and vault ID.%s\n", errorColor, resetColor)
		os.Exit(1)
	}
	token, vaultID := args[0], args[1]
	vault, err := invoker.GetVault(token, vaultID)
	if err != nil {
		fmt.Printf("%sError: Failed to retrieve text from vault: %s%s\n", errorColor, err, resetColor)
		os.Exit(1)
	}
	printVault(vault)
}

// handleGetBinaryVault handles retrieving binary data from the vault.
func handleGetBinaryVault(invoker Invoker, args []string) {
	if len(args) < 2 {
		fmt.Printf("%sError: Get binary vault command requires token and vault ID.%s\n", errorColor, resetColor)
		os.Exit(1)
	}
	token, vaultID := args[0], args[1]
	vault, err := invoker.GetVault(token, vaultID)
	if err != nil {
		fmt.Printf("%sError: Failed to retrieve binary data from vault: %s%s\n", errorColor, err, resetColor)
		os.Exit(1)
	}
	printVault(vault)
}

// handleGetCardVault handles retrieving card data from the vault.
func handleGetCardVault(invoker Invoker, args []string) {
	if len(args) < 2 {
		fmt.Printf("%sError: Get card vault command requires token and vault ID.%s\n", errorColor, resetColor)
		os.Exit(1)
	}
	token, vaultID := args[0], args[1]
	vault, err := invoker.GetVault(token, vaultID)
	if err != nil {
		fmt.Printf("%sError: Failed to retrieve card data from vault: %s%s\n", errorColor, err, resetColor)
		os.Exit(1)
	}
	printVault(vault)
}

// handleStoreBinaryVault handles storing arbitrary binary data in the vault.
func handleStoreBinaryVault(invoker Invoker, args []string) {
	if len(args) < 3 {
		fmt.Printf("%sError: Store binary vault command requires token, vault ID, and file path.%s\n", errorColor, resetColor)
		os.Exit(1)
	}
	token, vaultID, filePath := args[0], args[1], args[2]

	fileData, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("%sError: Failed to read file %s: %s%s\n", errorColor, filePath, err, resetColor)
		os.Exit(1)
	}

	key := "binary/" + vaultID
	value := fmt.Sprintf(`{"file_name": "%s", "data": "%x"}`, filePath, fileData)

	vault, err := invoker.NewVault(token, key, value, vaultID)
	if err != nil {
		fmt.Printf("%sError: Failed to store binary data in vault: %s%s\n", errorColor, err, resetColor)
		os.Exit(1)
	}
	printVault(vault)
}

// handleStoreCardVault handles storing bank card data in the vault.
func handleStoreCardVault(invoker Invoker, args []string) {
	if len(args) < 5 {
		fmt.Printf("%sError: Store card vault command requires token, vault ID, card number, expiry date, and CVV.%s\n", errorColor, resetColor)
		os.Exit(1)
	}
	token, vaultID, cardNumber, expiryDate, cvv := args[0], args[1], args[2], args[3], args[4]
	key := "card/" + vaultID
	value := fmt.Sprintf(`{"card_number": "%s", "expiry_date": "%s", "cvv": "%s"}`, cardNumber, expiryDate, cvv)

	vault, err := invoker.NewVault(token, key, value, vaultID)
	if err != nil {
		fmt.Printf("%sError: Failed to store card data in vault: %s%s\n", errorColor, err, resetColor)
		os.Exit(1)
	}
	printVault(vault)
}

// printHelp displays usage information for the CLI tool.
func printHelp() {
	fmt.Println("GophKeeper CLI Help")
	fmt.Println("---------------------")
	fmt.Println("This is a CLI tool to interact with the GophKeeper service. Below are the available commands:")
	fmt.Println()
	fmt.Println("1. Sign Up")
	fmt.Println("Description: Register a new user.")
	fmt.Println("Usage: ./client signup <server_url> <username> <password>")
	fmt.Println("Example:")
	fmt.Println("  ./client signup http://localhost:8080 testName testPassword")
	fmt.Println()
	fmt.Println("2. Sign In")
	fmt.Println("Description: Log in a user and retrieve an authentication token.")
	fmt.Println("Usage: ./client signin <server_url> <username> <password>")
	fmt.Println("Example:")
	fmt.Println("  ./client signin http://localhost:8080 testName testPassword")
	fmt.Println()
	fmt.Println("3. Create Vault")
	fmt.Println("Description: Create a new vault using an authentication token.")
	fmt.Println("Usage: ./client savevault <server_url> <token> <key> <value> [<vault_id>]")
	fmt.Println("Example:")
	fmt.Println("  ./client savevault http://localhost:8080 <token> k1 v1")
	fmt.Println("  ./client savevault http://localhost:8080 <token> k1 v1 <vault_id>")
	fmt.Println()
	fmt.Println("4. Get Vault")
	fmt.Println("Description: Retrieve vault details using an authentication token and vault ID.")
	fmt.Println("Usage: ./client getvault <server_url> <token> <vault_id>")
	fmt.Println("Example:")
	fmt.Println("  ./client getvault http://localhost:8080 <token> 1")
	fmt.Println()
	fmt.Println("5. Upload File")
	fmt.Println("Description: Upload a binary file to the server using an authentication token.")
	fmt.Println("Usage: ./client uploadfile <server_url> <token> <filename> <file_path> [<vault_id>]")
	fmt.Println("Example:")
	fmt.Println("  ./client uploadfile http://localhost:8080 <token> filename3 /home/user/file.md")
	fmt.Println("  ./client uploadfile http://localhost:8080 <token> filename3 /home/user/file.md <vault_id>")
	fmt.Println()
	fmt.Println("6. Store Login/Password Vault")
	fmt.Println("Description: Store a login and password pair in the vault.")
	fmt.Println("Usage: ./client store-login-password <server_url> <token> <vault_id> <login> <password>")
	fmt.Println("Example:")
	fmt.Println("  ./client store-login-password http://localhost:8080 <token> 1 user1 password1")
	fmt.Println()
	fmt.Println("7. Get Login/Password Vault")
	fmt.Println("Description: Retrieve a login and password pair from the vault.")
	fmt.Println("Usage: ./client get-login-password <server_url> <token> <vault_id>")
	fmt.Println("Example:")
	fmt.Println("  ./client get-login-password http://localhost:8080 <token> 1")
	fmt.Println()
	fmt.Println("8. Store Text Vault")
	fmt.Println("Description: Store arbitrary text data in the vault.")
	fmt.Println("Usage: ./client store-text <server_url> <token> <vault_id> <text>")
	fmt.Println("Example:")
	fmt.Println("  ./client store-text http://localhost:8080 <token> 1 \"This is some text to store\"")
	fmt.Println()
	fmt.Println("9. Get Text Vault")
	fmt.Println("Description: Retrieve text data from the vault.")
	fmt.Println("Usage: ./client get-text <server_url> <token> <vault_id>")
	fmt.Println("Example:")
	fmt.Println("  ./client get-text http://localhost:8080 <token> 1")
	fmt.Println()
	fmt.Println("10. Store Binary Vault")
	fmt.Println("Description: Store arbitrary binary data (e.g., a file) in the vault.")
	fmt.Println("Usage: ./client store-binary <server_url> <token> <vault_id> <file_path>")
	fmt.Println("Example:")
	fmt.Println("  ./client store-binary http://localhost:8080 <token> 1 /path/to/file.bin")
	fmt.Println()
	fmt.Println("11. Get Binary Vault")
	fmt.Println("Description: Retrieve binary data from the vault.")
	fmt.Println("Usage: ./client get-binary <server_url> <token> <vault_id>")
	fmt.Println("Example:")
	fmt.Println("  ./client get-binary http://localhost:8080 <token> 1")
	fmt.Println()
	fmt.Println("12. Store Card Vault")
	fmt.Println("Description: Store bank card data in the vault.")
	fmt.Println("Usage: ./client store-card <server_url> <token> <vault_id> <card_number> <expiry_date> <cvv>")
	fmt.Println("Example:")
	fmt.Println("  ./client store-card http://localhost:8080 <token> 1 1234567890123456 12/24 123")
	fmt.Println()
	fmt.Println("13. Get Card Vault")
	fmt.Println("Description: Retrieve bank card data from the vault.")
	fmt.Println("Usage: ./client get-card <server_url> <token> <vault_id>")
	fmt.Println("Example:")
	fmt.Println("  ./client get-card http://localhost:8080 <token> 1")
	fmt.Println()
	fmt.Println("14. Version Information")
	fmt.Println("Description: Get the version and build date of the client.")
	fmt.Println("Usage: ./client --version")
	fmt.Println("Example:")
	fmt.Println("  ./client --version")
}
