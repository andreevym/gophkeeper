# GophKeeper Server

The GophKeeper Server is a crucial component of the GophKeeper system, designed to handle user authentication, data storage, and secure communication with the client application. It operates by exposing endpoints that the client can interact with to manage and access secure vaults.

---

## Table of Contents

1. [Overview](Overview)
2. [Getting Started](#getting-started)
   - [Prerequisites](#prerequisites)
   - [Build the Application](#build-the-application)
   - [Running the Application](#running-the-application)
3. [Configuration Options](#configuration-options)
   - [Environment Variables](#environment-variables)
   - [Command-line Flags](#command-line-flags)
   - [JSON Configuration File](#json-configuration-file)
4. [Usage](#usage)
5. [Examples](#examples)
   - [Example 1: Running with Environment Variables](#example-1-running-with-environment-variables)
   - [Example 2: Running with Command-line Flags](#example-2-running-with-command-line-flags)
   - [Example 3: Running with a JSON Configuration File](#example-3-running-with-a-json-configuration-file)

---

Here's an **Overview** section for the **GophKeeper Server** documentation:

---

## Overview

The **GophKeeper Server** is a crucial component of the GophKeeper system, designed to handle user authentication, data storage, and secure communication with the client application. It operates by exposing endpoints that the client can interact with to manage and access secure vaults.

### Key Features:
- **Authentication**: The server manages user registration and login processes, issuing JSON Web Tokens (JWTs) for secure authentication.
- **Data Management**: It handles the creation, storage, and retrieval of user vaults. Vaults are used to store sensitive information such as passwords, binary data, and other private details.
- **Configuration**: The server offers flexible configuration options, allowing you to specify settings through environment variables, command-line flags, or a JSON configuration file.
- **Logging**: It supports configurable logging levels to help with debugging and monitoring server operations.

### Components:
1. **Server Application**: The main executable that runs the server, listens for incoming requests, processes them, and interacts with the database.
2. **Database**: Stores user credentials, vault data, and other application-related information. The default setup uses PostgreSQL, but you can configure it to use other databases if needed.

### Configuration:
The server can be configured in several ways:
- **Environment Variables**: Set environment variables to configure the server's address, database connection, log level, and JWT secret key.
- **Command-line Flags**: Pass configuration options as flags when starting the server.
- **JSON Configuration File**: Use a JSON file to specify all configuration settings in a structured format.

This flexibility in configuration allows you to tailor the server to different deployment environments and operational requirements, ensuring that it meets the specific needs of your application.

---

## Getting Started

To run the application, you can configure the server using one of the following methods:
1. **Environment Variables**
2. **Command-line Flags**
3. **JSON Configuration File**

### Prerequisites

- Go 1.18+ installed
- PostgreSQL database (if you're using the default connection string)

### Build the Application

To build the application, run the following command:

```bash
go build -o client
```

### Running the Application

Once built, run the application using:

```bash
./client
```

You can provide configuration values using environment variables, flags, or a JSON file, as described below.

## Configuration Options

The application can be configured using the following parameters, either via environment variables or command-line flags.

| Option           | Environment Variable | Command-line Flag | Default Value                                                  | Description                         |
|------------------|----------------------|-------------------|-----------------------------------------------------------------|-------------------------------------|
| Server Address   | `ADDRESS`             | `-a`              | `:8080`                                                         | The address the server listens on   |
| Database URI     | `DATABASE_URI`        | `-d`              | `postgres://gophkeeper:gophkeeper@localhost:5432/gophkeeper?sslmode=disable` | The URI for the database connection |
| Log Level        | `LOG_LEVEL`           | `-l`              | `info`                                                          | The log level (e.g., `info`, `debug`, `error`) |
| JWT Secret Key   | `JWT_SECRET_KEY`      | `-j`              | None                                                            | The secret key for JWT token signing |

### Environment Variables

You can set the following environment variables to configure the application:

- `ADDRESS`: The server address (e.g., `:8080`, `127.0.0.1:9000`).
- `DATABASE_URI`: The database connection URI.
- `LOG_LEVEL`: The log level (e.g., `info`, `debug`, `error`).
- `JWT_SECRET_KEY`: The secret key used to sign JWT tokens.

Example:

```bash
export ADDRESS=":9090"
export DATABASE_URI="postgres://user:password@localhost:5432/mydb?sslmode=disable"
export LOG_LEVEL="debug"
export JWT_SECRET_KEY="my-secret-key"
```

### Command-line Flags

You can also use flags to configure the application when starting it:

- `-a`: Server address (e.g., `-a :9090`).
- `-d`: Database URI (e.g., `-d postgres://user:password@localhost:5432/mydb`).
- `-l`: Log level (e.g., `-l debug`).
- `-j`: JWT Secret Key (e.g., `-j my-secret-key`).

Example:

```bash
./client -a :9090 -d postgres://user:password@localhost:5432/mydb -l debug -j my-secret-key
```

### JSON Configuration File

You can also use JSON file to configure the application when starting it:

Example `config.json`:

```json
{
    "Address": ":9090",
    "DatabaseURI": "postgres://user:password@localhost:5432/mydb?sslmode=disable",
    "LogLevel": "debug",
    "JWTSecretKey": "my-secret-key"
}
```

## Usage

After setting up the configuration, you can run the server:

1. Using environment variables:
   ```bash
   ADDRESS=":9090" DATABASE_URI="postgres://user:password@localhost:5432/mydb?sslmode=disable" LOG_LEVEL="debug" JWT_SECRET_KEY="my-secret-key" ./clien
   ```

2. Using command-line flags:
   ```bash
   ./clien -a :9090 -d postgres://user:password@localhost:5432/mydb -l debug -j my-secret-key
   ```

3. Using a JSON config file:
   ```bash
   ./clien
   ```

## Examples

### Example 1: Running with Environment Variables

```bash
export ADDRESS=":9090"
export DATABASE_URI="postgres://user:password@localhost:5432/mydb?sslmode=disable"
export LOG_LEVEL="debug"
export JWT_SECRET_KEY="my-secret-key"
./clien
```

### Example 2: Running with Command-line Flags

```bash
./clien -a :9090 -d postgres://user:password@localhost:5432/mydb -l debug -j my-secret-key
```

### Example 3: Running with a JSON Configuration File

First, create a JSON file (`config.json`):

```json
{
    "Address": ":9090",
    "DatabaseURI": "postgres://user:password@localhost:5432/mydb",
    "LogLevel": "debug",
    "JWTSecretKey": "my-secret-key"
}
```

```bash
./clien
```