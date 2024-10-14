# Referral System

## Description

A simple RESTful API service for managing a referral system.

## Features

- User registration and authentication
- Creating and deleting referral codes
- Retrieving referral code by email
- Registering via referral code
- Retrieving information about referrals
- API Documentation (Swagger)

## Technology Stack

- Go
- Gin
- GORM
- PostgreSQL
- JWT
- Swagger

## Installation

1. **Clone the repository:**

    ```bash
    git clone https://github.com/serlenario/referral-system.git
    cd referral-system
    ```

2. **Create a `.env` file based on `.env.example` and fill in the environment variables:**

    ```env
    DB_HOST=localhost
    DB_PORT=5432
    DB_USER=postgres
    DB_PASSWORD=password
    DB_NAME=referral_db
    JWT_SECRET=your_jwt_secret
    ```

3. **Install dependencies:**

    ```bash
    go mod download
    ```

4. **Run migrations:**

    With GORM, migrations are automatically handled when the application starts.

5. **Generate Swagger documentation:**

    ```bash
    swag init
    ```

6. **Start the application:**

    ```bash
    go run cmd/main.go
    ```

## Testing

Use [Postman](https://www.postman.com/) to send requests to the API.

## Documentation

Available at: `http://localhost:8080/swagger/index.html`
