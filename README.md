# Banking Ledger System

A modular, scalable banking ledger system built with Go, designed to handle core banking operations such as account management, transaction processing, and ledger maintenance. This project emphasizes clean architecture, separation of concerns, and ease of deployment using Docker.

## ✨ Features

- **User Registration & Login**: with unique emails
- **Account Management**: Create customer accounts
- **Transaction Processing**: Deposits and withdrawals with per-user isolation
- **Ledger Maintenance**: Consistent and auditable transaction ledger
- **Scoped Access**: Users can only perform transactions on their own accounts
- **Modular Architecture** with clear layering of handlers, services, and models
- **Structured Logging** for observability
- **Middleware Support**: Auth, logging, and request validation
- **Dockerized Setup** for easy local development

## 🛠️ Tech Stack

### Backend:
- **Go (Golang)** – Service logic
- **Gin** – Web framework for APIs
- **JWT (github.com/golang-jwt/jwt)** – Authentication and user scoping

### Database:
- **PostgreSQL** – Stores user accounts, balances, and transactions
- **MongoDB**  – For logs or audit trails

### Messaging:
- **Kafka**  – For async transaction handling and event-based updates

### DevOps:
- **Docker & Docker Compose** – Containerized development environment

### Testing & Tooling:
- **Zap** – Structured logger
- **godotenv** – Manage environment variables
- **Swagger** – API documentation 
- **Postman** - API testing

## 🔐 Authentication

The system uses **JWT (JSON Web Tokens)** for user authentication. Here's how it works:

1. **User Registration/Login**: Upon successful login, the user receives a JWT token.
2. **Protected Endpoints**: All account and transaction APIs require the token in the `Authorization` header (`Bearer <token>`).
3. **Scoped Access**:
   - Every account is associated with a specific user.
   - Middleware extracts the user ID from the JWT and injects it into the request context.
   - Transaction endpoints verify ownership before processing.
   - Users **cannot access or transact on accounts that don’t belong to them**, ensuring security and data isolation.

Example header:
```
Authorization: Bearer <your-jwt-token>
```

## 🚀 Getting Started

### Prerequisites

- Go 1.21+
- Docker
- Docker Compose

### Installation

```bash
git clone https://github.com/abhinaya-k/banking_ledger.git
cd banking_ledger
```

Set up environment variables in a `.env` file:

```env
# Server Config
SERVICE_NAME ="banking_ledger"
SERVICE_BASE_PATH="/bankingLedger"
ENV ="local"
SERVER_PORT =8080

# PostgreSQL Config
DB_HOST="localhost"
DB_PORT=5432
DB_USER="yourusername"
DB_PASSWORD="yourpassword"
DB_NAME="banking_ledger"

# Kafka Config
KAFKA_BROKER="ledger-kafka:29092"
KAFKA_USERNAME="admin"
KAFKA_PASSWORD="admin"
TRANSACTION_PROCESSING_KAFKA_TOPIC="your-kafka-topic"
TRANSACTION_PROCESSING_KAFKA_CG = "your-kafka-consumer-group"

# MongoDB Config
MONGO_HOST="ledger-mongo"
MONGO_PORT="27017"
MONGO_DB_NAME="bankingLedger"

# JWT Secret
JWT_SECRET=your_jwt_secret_key

API_KEY="SoTq5ZoWt8jjl8z7MoAGiHN1BATI5j6k"
```

Start the services:

```bash
docker-compose up --build
```

## 📘 API Endpoints (Sample)

- `POST /bankingLedger/user/v1/register`: Register a new user as an admin or a normal user
- `POST /bankingLedger/user/v1/login`: Login and receive JWT

*NOTE: The above two apis are authenticated using an api key from the .env file and the apis below are authenticated using a JWT token*
- `POST /bankingLedger/v1/account`: Create a user-scoped account
- `PATCH /bankingLedger/v1/account/transaction`: Deposit or withdraw from own account
- `POST /bankingLedger/v1/account/ledger`: View transaction history (admins can view history of all users, user roles can only view their own transactions)
