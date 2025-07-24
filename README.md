# Chat Application with Gorilla WebSocket and Gin Authentication

This project is a chat application built using **Gorilla WebSocket** for real-time communication and **Gin** as the HTTP framework for API routing and authentication. The application supports user authentication and secure WebSocket connections for one-on-one messaging.

---

## Project Pattern
This project is using the repository pattern, which is beneficial because it separates data access logic from business logic, resulting in a cleaner and more maintainable codebase. This separation enhances testability by allowing the use of mock repositories without relying on the actual database. It also increases flexibility, as changes to the underlying data source require minimal modifications to the business layer. By promoting adherence to the Single Responsibility Principle, the repository pattern helps ensure a well-structured and scalable architecture.

## Features

- **User Authentication:**
  - Access and refresh token-based authentication.
  - Refresh tokens are stored securely in Redis.
- **WebSocket Communication:**
  - Establish WebSocket connections with authentication.
  - Supports direct messaging between users (one-on-one chat).
  - Messages are persisted when the recipient is offline and delivered upon reconnection.
- **Token Validation:**
  - Access tokens are validated for each request.
  - refresh token is compatible and stored in redis.
- **Secure Routes:**
  - Protects WebSocket endpoints with middleware to validate user credentials.

---

## Requirements

- Go 1.23+
- Redis
- Gorilla WebSocket library
- Gin HTTP framework

---

## Project Structure

```plaintext
├───cmd
│   └───server
├───configs
├───docs
├───internal
│   ├───app
│   ├───controller
│   │   ├───auth
│   │   ├───chat
│   │   └───user
│   ├───http
│   │   ├───request
│   │   ├───response
│   │   └───server
│   ├───middleware
│   ├───model
│   ├───notification
│   ├───repository
│   │   ├───auth
│   │   ├───chat
│   │   └───user
│   ├───service
│   │   ├───auth
│   │   ├───chat
│   │   ├───email
│   │   └───user
│   ├───storage
│   └───util
│       └───error
├───log
├───pkg
│   ├───db
│   └───util
├───scripts
├───swagger
└───test
```

---

## Installation and Setup

1. **Clone the repository:**

   ```bash
   git clone https://github.com/your-repository/chat-app.git
   cd chat-app
   ```

2. **Install dependencies:**

   ```bash
   go mod tidy
   ```

3. **Configure environment variables:**
   Create a `.env` file in the root directory and add the following variables:

   ```env
    APP_NAME="your-app"
    APP_HOST="localhost:3000"
    APP_PORT=3000
    JWT_SECRET="your-secret"
    DB_HOST="postgres://postgres:password@localhost:5432/database"
    GIN_MODE="debug"
    REDIS_HOST="example:6379"
    REDIS_PASS="your pass"
    SMTP_HOST=smtp.example.com
    SMTP_PORT="number"
    SMTP_USER="string@mail"
    SMTP_PASS="string"
   ```

4. **Run the application:**
   ```bash
   go run main.go
   ```

---

## API Endpoints

### Authentication Endpoints

1. **Sign Up -> request otp:**

   ```http
   POST /api/v1/auth/signup/request-otp
   ```

   Request body:

   ```json
   {
     "email": "string"
   }
   ```

2. **Sign Up -> verify otp:**

   ```http
   POST /api/v1/auth/signup/verify-otp
   ```

   Request body:

   ```json
   {
     "email": "string",
     "otp": "string",
     "password": "string",
     "username": "string"
   }
   ```

3. **Sign In:**

   ```http
   POST /api/v1/auth/signin
   ```

   Request body:

   ```json
   {
     "email": "string",
     "password": "string"
   }
   ```

4. **Refresh Token:**
   ```http
   POST /api/v1/auth/refresh
   ```
   Request body:
   ```json
   {
     "refresh_token": "string"
   }
   ```

### WebSocket Endpoint

**Connect to WebSocket:**

```http
GET /api/v1/chat/ws
```

Headers:

```http
Authorization: Bearer <access_token>
```

---

## Middleware

### Authentication Middleware

- **Description:** Ensures that all requests have a valid access token.
- **Functionality:**
  - Verifies JWT tokens.
  - Extracts user credentials and injects them into the Gin context.

---

## WebSocket Message Structure

### Client to Server

```json
{
  "receiver_id": "string",
  "content": "string"
}
```

### Server to Client

```json
{
  "sender_id": "string",
  "content": "string",
  "timestamp": "string"
}
```

---

## Contributing

Feel free to fork the repository and submit pull requests. Please ensure your code follows the project’s style guidelines and is thoroughly tested.
