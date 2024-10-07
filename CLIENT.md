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
    - **Result:**
      ```bash
      User created successfully
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
   - **Result:**
     ```bash
     eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySUQiOiIyIn0.58Nr9-P22PnDQ6bwugjEMsd2_86KsZxRA9nUN8AGgzbB7tHh_RS4bJv3-ejiZ61yb5Di5Awm2nCxvS7weqL9QA
     ```

3. **Create Vault**

    - **Description:** Create a new vault entry with a key-value pair.
    - **Usage:**
      ```bash
      ./client savevault <server_url> <token> <key> <value> [<vault_id> only for update]
      ```
    - **Example:**
      ```bash
      ./client savevault http://localhost:8080 eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySUQiOiIyIn0.58Nr9-P22PnDQ6bwugjEMsd2_86KsZxRA9nUN8AGgzbB7tHh_RS4bJv3-ejiZ61yb5Di5Awm2nCxvS7weqL9QA key1 value1
      ```
   - **Result:**
     ```bash
     Vault operation successful: {"id":12,"key":"key1","value":"value1","user_id":2}
     ```

4. **Get Vault**

    - **Description:** Retrieve a vault entry by its ID.
    - **Usage:**
      ```bash
      ./client getvault <server_url> <token> <vault_id>
      ```
    - **Example:**
      ```bash
      ./client getvault http://localhost:8080 eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySUQiOiIyIn0.58Nr9-P22PnDQ6bwugjEMsd2_86KsZxRA9nUN8AGgzbB7tHh_RS4bJv3-ejiZ61yb5Di5Awm2nCxvS7weqL9QA 12
      ```
   - **Result:**
     ```bash
     Vault operation successful: {"id":12,"key":"key1","value":"value1","user_id":2}
     ```

5. **Upload File**

    - **Description:** Upload a binary file to the vault.
    - **Usage:**
      ```bash
      ./client uploadfile <server_url> <token> <filename> <file_path> [<vault_id> only for update]
      ```
    - **Example:**
      ```bash
      ./client uploadfile http://localhost:8080 eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySUQiOiIyIn0.58Nr9-P22PnDQ6bwugjEMsd2_86KsZxRA9nUN8AGgzbB7tHh_RS4bJv3-ejiZ61yb5Di5Awm2nCxvS7weqL9QA file1 ./README.md 100
      ```

6. **Store Login/Password Pair**

    - **Description:** Store a login and password pair in the vault.
    - **Usage:**
      ```bash
      ./client store-login-password <server_url> <token> <login> <password> [<vault_id> only for update]
      ```
    - **Example create:**
      ```bash
      ./client store-login-password http://localhost:8080 eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySUQiOiIyIn0.58Nr9-P22PnDQ6bwugjEMsd2_86KsZxRA9nUN8AGgzbB7tHh_RS4bJv3-ejiZ61yb5Di5Awm2nCxvS7weqL9QA user1 password1
      ```
    - **Result create:**
      ```bash
      Vault operation successful: {"id":13,"key":"login/user1","value":"{\"login\": \"user1\", \"password\": \"password1\"}","user_id":2}
      ```
    - Example update:
      ```bash
      ./client store-login-password http://localhost:8080 eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySUQiOiIyIn0.58Nr9-P22PnDQ6bwugjEMsd2_86KsZxRA9nUN8AGgzbB7tHh_RS4bJv3-ejiZ61yb5Di5Awm2nCxvS7weqL9QA user1 password2 13
      ```
    - **Result update:**
      ```bash
      Vault operation successful: {"id":13,"key":"login/user1","value":"{\"login\": \"user1\", \"password\": \"password2\"}","user_id":2}
      ```

7. **Get Login/Password Pair**

    - **Description:** Retrieve a login and password pair from the vault.
    - **Usage:**
      ```bash
      ./client get-login-password <server_url> <token> <vault_id>
      ```
    - **Example:**
      ```bash
      ./client get-login-password http://localhost:8080 eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySUQiOiIyIn0.58Nr9-P22PnDQ6bwugjEMsd2_86KsZxRA9nUN8AGgzbB7tHh_RS4bJv3-ejiZ61yb5Di5Awm2nCxvS7weqL9QA 13
      ```
    - **Result:**
      ```bash
      Vault operation successful: {"id":13,"key":"login/user1","value":"{\"login\": \"user1\", \"password\": \"password2\"}","user_id":2}
      ```

8. **Store Text**

    - **Description:** Store arbitrary text in the vault.
    - **Usage:**
      ```bash
      ./client store-text <server_url> <token> <key> <text> [<vault_id> only for update]
      ```
    - **Example create:**
      ```bash
      ./client store-text http://localhost:8080 eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySUQiOiIyIn0.58Nr9-P22PnDQ6bwugjEMsd2_86KsZxRA9nUN8AGgzbB7tHh_RS4bJv3-ejiZ61yb5Di5Awm2nCxvS7weqL9QA k1 "Some important text"
      ```
    - **Result create:**
      ```bash
      Vault operation successful: {"id":14,"key":"text/k1","value":"{\"text\": \"Some important text\"}","user_id":2}
      ```
    - Example update:
      ```bash
      ./client store-text http://localhost:8080 eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySUQiOiIyIn0.58Nr9-P22PnDQ6bwugjEMsd2_86KsZxRA9nUN8AGgzbB7tHh_RS4bJv3-ejiZ61yb5Di5Awm2nCxvS7weqL9QA k1 "Some important text" 13
      ```
    - **Result update:**
      ```bash
      Vault operation successful: {"id":14,"key":"text/k1","value":"{\"text\": \"Some important text 2\"}","user_id":2}
      ```

9. **Get Text**

    - **Description:** Retrieve text from the vault.
    - **Usage:**
      ```bash
      ./client get-text <server_url> <token> <vault_id>
      ```
    - **Example:**
      ```bash
      ./client get-text http://localhost:8080 eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySUQiOiIyIn0.58Nr9-P22PnDQ6bwugjEMsd2_86KsZxRA9nUN8AGgzbB7tHh_RS4bJv3-ejiZ61yb5Di5Awm2nCxvS7weqL9QA 14
      ``` 
    - **Result:**
      ```bash
      Vault operation successful: {"id":14,"key":"text/k1","value":"{\"text\": \"Some important text 2\"}","user_id":2}
      ```


10. **Store Binary Data**

    - **Description:** Store arbitrary binary data (e.g., a file) in the vault.
    - **Usage:**
      ```bash
      ./client store-binary <server_url> <token> <key> <file_path> [<vault_id> only for update]
      ```
    - **Example create:**
      ```bash
      ./client store-binary http://localhost:8080 eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySUQiOiIyIn0.58Nr9-P22PnDQ6bwugjEMsd2_86KsZxRA9nUN8AGgzbB7tHh_RS4bJv3-ejiZ61yb5Di5Awm2nCxvS7weqL9QA clientbin ./client
      ```
    - **Result create:**
      ```bash
      Vault saved successfully by id: 18
      ```
    - Example update:
      ```bash
      ./client store-binary http://localhost:8080 eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySUQiOiIyIn0.58Nr9-P22PnDQ6bwugjEMsd2_86KsZxRA9nUN8AGgzbB7tHh_RS4bJv3-ejiZ61yb5Di5Awm2nCxvS7weqL9QA clientbin ./client 18
      ```
    - **Result update:**
      ```bash
      Vault saved successfully by id: 18
      ```

11. **Get Binary Data**

    - **Description:** Retrieve binary data from the vault.
    - **Usage:**
      ```bash
      ./client get-binary <server_url> <token> <vault_id> <file_path>
      ```
    - **Example:**
      ```bash
      ./client get-binary http://localhost:8080 eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySUQiOiIyIn0.58Nr9-P22PnDQ6bwugjEMsd2_86KsZxRA9nUN8AGgzbB7tHh_RS4bJv3-ejiZ61yb5Di5Awm2nCxvS7weqL9QA 18 client2
      ```
    - **Result:**
      ```bash
      Vault 18 saved successfully by path: client2
      ```
    - Note: The binary data is saved in the file client2 in the current directory.

12. **Store Card Information**

    - **Description:** Store bank card information in the vault.
    - **Usage:**
      ```bash
      ./client store-card <server_url> <token> <key> <card_number> <expiry_date> <cvv> [<vault_id> only for update]
      ```
    - **Example create:**
      ```bash
      ./client store-card http://localhost:8080 eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySUQiOiIyIn0.58Nr9-P22PnDQ6bwugjEMsd2_86KsZxRA9nUN8AGgzbB7tHh_RS4bJv3-ejiZ61yb5Di5Awm2nCxvS7weqL9QA tbankvisacard1 4242424242424242 12/25 123
      ```
    - **Result create:**
      ```bash
      Vault operation successful: {"id":15,"key":"card/tbankvisacard1","value":"{\"card_number\": \"4242424242424242\", \"expiry_date\": \"12/25\", \"cvv\": \"123\"}","user_id":2}
      ```
    - Example update:
      ```bash
      ./client store-card http://localhost:8080 eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySUQiOiIyIn0.58Nr9-P22PnDQ6bwugjEMsd2_86KsZxRA9nUN8AGgzbB7tHh_RS4bJv3-ejiZ61yb5Di5Awm2nCxvS7weqL9QA tbankvisacard1 4242424242424242 12/25 333 15
      ```
    - **Result update:**
      ```bash
      Vault operation successful: {"id":15,"key":"card/tbankvisacard1","value":"{\"card_number\": \"4242424242424242\", \"expiry_date\": \"12/25\", \"cvv\": \"333\"}","user_id":2}
      ```

13. **Get Card Information**

    - **Description:** Retrieve bank card information from the vault.
    - **Usage:**
      ```bash
      ./client get-card <server_url> <token> <vault_id>
      ```
    - **Example:**
      ```bash
      ./client get-card http://localhost:8080 eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySUQiOiIyIn0.58Nr9-P22PnDQ6bwugjEMsd2_86KsZxRA9nUN8AGgzbB7tHh_RS4bJv3-ejiZ61yb5Di5Awm2nCxvS7weqL9QA 15
      ```
    - **Result:**
      ```bash
      Vault operation successful: {"id":15,"key":"card/tbankvisacard1","value":"{\"card_number\": \"4242424242424242\", \"expiry_date\": \"12/25\", \"cvv\": \"333\"}","user_id":2}
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
