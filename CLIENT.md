# GophKeeper Client

GophKeeper Client is the client-side application of the GophKeeper system, which allows users to securely store and manage their vaults containing logins, passwords, binary data, and other private information.

This document provides instructions for building and using the GophKeeper Client to interact with the server.

---

## Table of Contents

1. [Build the Client Application](#build-the-client-application)
2. [Usage of Client Application](#usage-of-client-application)
   - [1. Sign Up](#1-sign-up)
   - [2. Sign In to Receive Token](#2-sign-in-to-receive-token)
   - [3. Create a New Vault with Token](#3-create-a-new-vault-with-token)
   - [4. Get Vault with Token](#4-get-vault-with-token)
3. [Summary of Commands](#summary-of-commands)

---

## Build the Client Application

To build the GophKeeper Client application, run the following command:

```bash
go build -o client
```

This will create a binary named `client` that you can use to perform operations such as signing up, signing in, creating vaults, and retrieving vault information.

## Usage of Client Application

The client application allows you to create your own secure vault. Below are the instructions for various client operations.

### 1. Sign Up

To register a new user, use the `signUp` command along with the server URL, username, and password.

Example:

```bash
./client signUp http://localhost:8080 testName testPassword
```

This will register a new user with the username `testName` and the password `testPassword`.

### 2. Sign In to Receive Token

To log in and receive an authentication token, use the `signIn` command with the server URL, username, and password.

Example:

```bash
./client signIn http://localhost:8080 testName testPassword
```

#### Response

Upon successful login, you will receive a JWT token in the response, which you will need for further actions such as creating or retrieving vaults.

Example response:

```bash
eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySUQiOiIxIn0.bMFEsrvtCxd5i3SMn3E_8HcRx6RzNfTX2PI1eWXJsbNUbeG_VaEpf9trTcm4KsYqYp_wpLzMYEYKQCtQykb4lQ
```

### 3. Create a New Vault with Token

To create a vault, use the `saveVault` command along with the server URL, the JWT token, a key, and a value for the vault.

Example:

```bash
./client saveVault http://localhost:8080 "eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySUQiOiIxIn0.iok4gCKCJP3d7vXMUyDFEvgZQ2-hyyk85gvHvmoGkx5-aMByqGyq8GjfNcpgY1Mc31xRn-d0BHnmy3H1kwNWXg" k1 v1
```

In this example:
- `k1` is the key of the vault entry.
- `v1` is the value associated with the key `k1`.

### 4. Get Vault with Token

To retrieve a vault by its ID, use the `getVault` command along with the server URL, the JWT token, and the vault ID.

Example:

```bash
./client getVault http://localhost:8080 "eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySUQiOiIxIn0.iok4gCKCJP3d7vXMUyDFEvgZQ2-hyyk85gvHvmoGkx5-aMByqGyq8GjfNcpgY1Mc31xRn-d0BHnmy3H1kwNWXg" 1
```

In this example:
- `1` is the ID of the vault entry you want to retrieve.

### Summary of Commands

| Command        | Description                                      | Example Usage                                                                                      |
|----------------|--------------------------------------------------|-----------------------------------------------------------------------------------------------------|
| `signUp`       | Register a new user                              | `./client signUp http://localhost:8080 testName testPassword`                                       |
| `signIn`       | Log in and get a JWT token                       | `./client signIn http://localhost:8080 testName testPassword`                                       |
| `saveVault`    | Create a new vault entry                         | `./client saveVault http://localhost:8080 <JWT_TOKEN> k1 v1`                                        |
| `getVault`     | Retrieve a vault entry by ID                     | `./client getVault http://localhost:8080 <JWT_TOKEN> 1`                                             |

This documentation provides basic usage instructions for the GophKeeper client application, enabling you to manage secure vault entries effectively.