# GophKeeper Client

GophKeeper Client is a command-line tool for interacting with the GophKeeper service, which allows users to securely store various types of sensitive information, such as login credentials, binary files, and text notes.

This document describes how to build and use the GophKeeper Client to manage your vaults.

## Table of Contents

1. [Build the Client Application](#build-the-client-application)
2. [Usage of Client Application](#usage-of-client-application)
    - [Commands Overview](#commands-overview)
    - [Command Details](#command-details)
---

## Build the Client Application

To build the GophKeeper Client application, run:

```bash
go build -o client cmd/client/main.go
```

This creates a binary named `client` that can be used for user registration, sign-in, vault management, and file uploads.

## Usage of Client Application

The GophKeeper Client provides a variety of commands for different tasks. Below is an overview of the available commands, followed by detailed usage information.

### Commands Overview

| Command                  | Description                                               |
|--------------------------|-----------------------------------------------------------|
| `signup`                 | Register a new user                                       |
| `signin`                 | Authenticate and retrieve an access token                 |
| `savevault`              | Create a new vault entry with a key-value pair            |
| `getvault`               | Retrieve a vault entry by its ID                          |
| `uploadfile`             | Upload a binary file to the vault                         |
| `store-login-password`   | Store a login and password pair in the vault              |
| `get-login-password`     | Retrieve a login and password pair from the vault         |
| `store-text`             | Store arbitrary text in the vault                         |
| `get-text`               | Retrieve text from the vault                              |
| `store-binary`           | Store arbitrary binary data in the vault                  |
| `get-binary`             | Retrieve binary data from the vault                       |
| `store-card`             | Store bank card information in the vault                  |
| `get-card`               | Retrieve bank card information from the vault             |
| `--version`              | Display version information                               |
| `help`                   | Display help information for all commands                 |

### Command Details

1. **Sign Up**

    - **Description:** Register a new user.
    - **Usage:**
      ```bash
      ./client signup <server_url> <username> <password>
      ```
    - **Example:**
      ```bash
      ./client signup http://localhost:8080 user1 password1
      ```

2. **Sign In**

    - **Description:** Authenticate and retrieve an access token.
    - **Usage:**
      ```bash
      ./client signin <server_url> <username> <password>
      ```
    - **Example:**
      ```bash
      ./client signin http://localhost:8080 user1 password1
      ```

3. **Create Vault**

    - **Description:** Create a new vault entry with a key-value pair.
    - **Usage:**
      ```bash
      ./client savevault <server_url> <token> <key> <value> [<vault_id>]
      ```
    - **Example:**
      ```bash
      ./client savevault http://localhost:8080 <token> key1 value1
      ./client savevault http://localhost:8080 <token> key1 value1 vaultID123
      ```

4. **Get Vault**

    - **Description:** Retrieve a vault entry by its ID.
    - **Usage:**
      ```bash
      ./client getvault <server_url> <token> <vault_id>
      ```
    - **Example:**
      ```bash
      ./client getvault http://localhost:8080 <token> vaultID123
      ```

5. **Upload File**

    - **Description:** Upload a binary file to the vault.
    - **Usage:**
      ```bash
      ./client uploadfile <server_url> <token> <filename> <file_path> [<vault_id>]
      ```
    - **Example:**
      ```bash
      ./client uploadfile http://localhost:8080 <token> file1 /path/to/file.txt
      ./client uploadfile http://localhost:8080 <token> file1 /path/to/file.txt vaultID123
      ```

6. **Store Login/Password Pair**

    - **Description:** Store a login and password pair in the vault.
    - **Usage:**
      ```bash
      ./client store-login-password <server_url> <token> <vault_id> <login> <password>
      ```
    - **Example:**
      ```bash
      ./client store-login-password http://localhost:8080 <token> vaultID123 user1 password1
      ```

7. **Get Login/Password Pair**

    - **Description:** Retrieve a login and password pair from the vault.
    - **Usage:**
      ```bash
      ./client get-login-password <server_url> <token> <vault_id>
      ```
    - **Example:**
      ```bash
      ./client get-login-password http://localhost:8080 <token> vaultID123
      ```

8. **Store Text**

    - **Description:** Store arbitrary text in the vault.
    - **Usage:**
      ```bash
      ./client store-text <server_url> <token> <vault_id> <text>
      ```
    - **Example:**
      ```bash
      ./client store-text http://localhost:8080 <token> vaultID123 "Some important text"
      ```

9. **Get Text**

    - **Description:** Retrieve text from the vault.
    - **Usage:**
      ```bash
      ./client get-text <server_url> <token> <vault_id>
      ```
    - **Example:**
      ```bash
      ./client get-text http://localhost:8080 <token> vaultID123
      ```

10. **Store Binary Data**

    - **Description:** Store arbitrary binary data (e.g., a file) in the vault.
    - **Usage:**
      ```bash
      ./client store-binary <server_url> <token> <vault_id> <file_path>
      ```
    - **Example:**
      ```bash
      ./client store-binary http://localhost:8080 <token> vaultID123 /path/to/file.bin
      ```

11. **Get Binary Data**

    - **Description:** Retrieve binary data from the vault.
    - **Usage:**
      ```bash
      ./client get-binary <server_url> <token> <vault_id>
      ```
    - **Example:**
      ```bash
      ./client get-binary http://localhost:8080 <token> vaultID123
      ```

12. **Store Card Information**

    - **Description:** Store bank card information in the vault.
    - **Usage:**
      ```bash
      ./client store-card <server_url> <token> <vault_id> <card_number> <expiry_date> <cvv>
      ```
    - **Example:**
      ```bash
      ./client store-card http://localhost:8080 <token> vaultID123 1234567812345678 12/25 123
      ```

13. **Get Card Information**

    - **Description:** Retrieve bank card information from the vault.
    - **Usage:**
      ```bash
      ./client get-card <server_url> <token> <vault_id>
      ```
    - **Example:**
      ```bash
      ./client get-card http://localhost:8080 <token> vaultID123
      ```

14. **Version Information**

    - **Description:** Display version and build information of the client.
    - **Usage:**
      ```bash
      ./client --version
      ```

15. **Help**

    - **Description:** Display help information for all commands.
    - **Usage:**
      ```bash
      ./client help
      ```
