# go-cli-db

This project is a Command Line Interface (CLI) tool that interfaces with a PostgreSQL database. It allows users to perform various database operations through a simple command-line interface.

## Features

- Connect to a PostgreSQL database
- Execute SQL queries
- Manage database entities using defined models
- Utility functions for common tasks

## Installation

To install the CLI tool, clone the repository and navigate to the project directory:

```bash
git clone <repository-url>
cd go-cli-db
```

Then, run the following command to install the necessary dependencies:

```bash
go mod tidy
```

## Usage

To run the CLI tool, use the following command:

```bash
go run cmd/main.go
```

You can pass various command-line arguments to interact with the database.
