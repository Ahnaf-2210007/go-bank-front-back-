# GoBank — Online Banking System

A lightweight RESTful banking API built in **Go** with **PostgreSQL**, providing account management, JWT-based authentication, email verification, money transfers with email notifications, coupon-based promotions, and transaction history.

---

## Table of Contents

- [Overview](#overview)
- [Tech Stack](#tech-stack)
- [Features](#features)
- [Project Structure](#project-structure)
- [Database Schema](#database-schema)
- [API Endpoints](#api-endpoints)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Setup](#setup)
  - [Build & Run](#build--run)
  - [Seed Test Data](#seed-test-data)
  - [Run Tests](#run-tests)
- [Authentication](#authentication)
- [Configuration](#configuration)

---

## Overview

GoBank is a backend banking system that exposes a RESTful API for managing bank accounts and performing money transfers. It handles user authentication with JSON Web Tokens (JWT), requires email verification when registering new accounts, sends email notifications after transfers, and stores all data in a PostgreSQL database.

---

## Tech Stack

| Category | Technology |
|---|---|
| Language | Go |
| HTTP Router | [Gorilla Mux](https://github.com/gorilla/mux) v1.8.1 |
| Database | PostgreSQL (`lib/pq` driver) |
| Authentication | JWT ([golang-jwt/jwt](https://github.com/golang-jwt/jwt) v5) |
| Password Hashing | bcrypt (`golang.org/x/crypto`) |
| Email | Gmail SMTP (`net/smtp`) |
| Environment | [godotenv](https://github.com/joho/godotenv) |
| Testing | [Testify](https://github.com/stretchr/testify) |
| Build Tool | Make |

---

## Features

- **Account registration with email verification** — New accounts go through a two-step flow: a 6-digit code is sent to the provided email address and must be verified before the account is activated (code expires in 5 minutes)
- **Account retrieval** — List all accounts or fetch a specific account by ID
- **Account deletion** — Remove accounts from the system
- **JWT authentication** — Secure login that returns a signed JWT token (valid for 24 hours)
- **Protected endpoints** — Middleware validates JWT tokens and enforces account ownership
- **Account profile updates** — Update name, email (two-step OTP flow), or password via a single protected endpoint
- **Money transfers** — Transfer funds between accounts with atomic database transactions and insufficient-funds protection
- **Transfer email notifications** — Both sender and recipient receive an email receipt after each transfer
- **Transaction history** — Paginated transaction history with optional filtering by type and month
- **Coupon/offer redemption** — Apply a one-time promotional coupon (`OFFER1000`) to add 1,000 to an account balance
- **Database migrations** — Automatic schema initialization and backward-compatible column migrations on startup
- **Environment configuration** — Supports `.env` file via godotenv with fallback to system environment variables

---

## Project Structure

```
Online-Banking-System/
├── main.go          # Entry point: loads .env, connects to DB, starts HTTP server on :3000
├── api.go           # HTTP handlers, routing, JWT middleware, email helpers
├── storage.go       # PostgreSQL data access layer
├── types.go         # Data models and request/response structs
├── types_test.go    # Unit tests for account creation logic
├── makefile         # Build, run, and test commands
├── go.mod           # Go module definition
└── go.sum           # Dependency checksums
```

### File Responsibilities

| File | Purpose |
|---|---|
| `main.go` | Loads `.env`, initializes all DB tables, optionally seeds data (`-seed` flag), and starts the API server |
| `api.go` | Defines all HTTP routes and handlers, JWT creation/validation, email sending (SMTP), and a request-wrapping helper for error handling |
| `storage.go` | Implements the `Storage` interface with PostgreSQL — creates, reads, updates, and deletes accounts; manages pending accounts, transfers, coupon redemptions, and transaction history |
| `types.go` | Defines `Account`, `PendingAccount`, `LoginRequest`, `LoginResponse`, `CreateAccountRequest`, `VerificationRequest`, `TransferRequest`, `TransferResult`, `TransactionRecord`, and `OfferRequest` structs |
| `types_test.go` | Unit tests verifying that `NewAccount` correctly creates an account, hashes the password, and generates a valid account number |

---

## Database Schema

The application uses five tables, all created automatically on startup.

### Entity Relationship Diagram

```mermaid
erDiagram
    account {
        SERIAL id PK
        VARCHAR first_name
        VARCHAR last_name
        BIGINT number
        VARCHAR email
        VARCHAR encrypted_password
        BIGINT balance
        TIMESTAMP created_at
    }

    pending_account {
        SERIAL id PK
        VARCHAR first_name
        VARCHAR last_name
        VARCHAR email
        BIGINT number
        VARCHAR encrypted_password
        VARCHAR verification_code
        TIMESTAMP expires_at
    }

    pending_profile_update {
        SERIAL id PK
        INT account_id FK
        VARCHAR new_email
        VARCHAR verification_code
        TIMESTAMP expires_at
        TIMESTAMP created_at
    }

    transfer {
        SERIAL id PK
        VARCHAR transaction_id
        INT from_account_id FK
        INT to_account_id FK
        BIGINT amount
        VARCHAR transaction_type
        VARCHAR status
        TIMESTAMP created_at
    }

    coupon_redemption {
        SERIAL id PK
        VARCHAR coupon_code
        INT account_id FK
        TIMESTAMP redeemed_at
    }

    account ||--o{ pending_profile_update : "has pending update"
    account ||--o{ transfer : "sends (from_account_id)"
    account ||--o{ transfer : "receives (to_account_id)"
    account ||--o| coupon_redemption : "redeems coupon"
    account ||--o| pending_account : "pending accounts"
```

### `account`

```sql
CREATE TABLE IF NOT EXISTS account (
    id                 SERIAL PRIMARY KEY,
    first_name         VARCHAR(255),
    last_name          VARCHAR(255),
    number             BIGINT,
    email              VARCHAR(255),
    encrypted_password VARCHAR(255) NOT NULL DEFAULT '',
    balance            BIGINT,
    created_at         TIMESTAMP
);
```

### `pending_account`

Holds accounts that have been registered but not yet email-verified. Rows are deleted after successful verification or when the code expires (5-minute TTL).

```sql
CREATE TABLE IF NOT EXISTS pending_account (
    id                 SERIAL PRIMARY KEY,
    first_name         VARCHAR(255),
    last_name          VARCHAR(255),
    email              VARCHAR(255) NOT NULL,
    number             BIGINT NOT NULL,
    encrypted_password VARCHAR(255) NOT NULL,
    verification_code  VARCHAR(6) NOT NULL,
    expires_at         TIMESTAMP NOT NULL
);
```

### `pending_profile_update`

Holds pending email-change requests awaiting OTP confirmation (5-minute TTL). Only one active pending update per account is allowed; a new request replaces the previous one.

```sql
CREATE TABLE IF NOT EXISTS pending_profile_update (
    id                SERIAL PRIMARY KEY,
    account_id        INT NOT NULL REFERENCES account(id),
    new_email         VARCHAR(255) NOT NULL,
    verification_code VARCHAR(6) NOT NULL,
    expires_at        TIMESTAMP NOT NULL,
    created_at        TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE (account_id, new_email)
);
```

### `transfer`

Records every completed money transfer.

```sql
CREATE TABLE IF NOT EXISTS transfer (
    id               SERIAL PRIMARY KEY,
    transaction_id   VARCHAR(32) NOT NULL UNIQUE,
    from_account_id  INT NOT NULL REFERENCES account(id),
    to_account_id    INT NOT NULL REFERENCES account(id),
    amount           BIGINT NOT NULL,
    transaction_type VARCHAR(50) NOT NULL DEFAULT 'transfer',
    status           VARCHAR(20) NOT NULL DEFAULT 'completed',
    created_at       TIMESTAMP NOT NULL
);
```

### `coupon_redemption`

Tracks which accounts have already redeemed a coupon (one redemption per account enforced by a unique index).

```sql
CREATE TABLE IF NOT EXISTS coupon_redemption (
    id          SERIAL PRIMARY KEY,
    coupon_code VARCHAR(255) NOT NULL,
    account_id  INT NOT NULL,
    redeemed_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

---

## API Endpoints

| Method | Path | Auth Required | Description |
|---|---|---|---|
| `POST` | `/login` | No | Authenticate with account number + password; returns a JWT token |
| `GET` | `/account` | No | List all accounts |
| `POST` | `/account` | No | Start account registration — validates input and sends a 6-digit verification code to the provided email |
| `POST` | `/account/verification` | No | Verify the 6-digit code; creates and returns the fully activated account |
| `POST` | `/account/update` | **JWT** | Update profile name, email (two-step OTP), or password — see action-based request shapes below |
| `GET` | `/account/{id}` | **JWT** | Get a specific account by ID |
| `DELETE` | `/account/{id}` | **JWT** | Delete an account by ID |
| `POST` | `/account/{id}/offer` | **JWT** | Redeem the `OFFER1000` coupon to add 1,000 to balance (one-time per account) |
| `GET` | `/account/transactions` | **JWT** | Get paginated transaction history for the authenticated user |
| `POST` | `/transfer` | **JWT** | Transfer funds between accounts; sends email receipts to both parties |

### Request / Response Examples

**POST `/account/update`** *(Authorization header required — action-based)*

All four actions share the same endpoint. Include a valid JWT in the `Authorization` header.

*Update name (`action = profile`)*
```json
// Request
{
  "action": "profile",
  "firstName": "Alice",
  "lastName": "Smith"
}

// Response
{
  "id": 1,
  "firstName": "Alice",
  "lastName": "Smith",
  "number": 482910,
  "email": "alice@example.com",
  "balance": 1000,
  "createdAt": "2024-01-15T10:30:00Z"
}
```

*Request email change (`action = email_request`) — step 1*
```json
// Request
{ "action": "email_request", "newEmail": "newalice@example.com", "password": "secret123" }

// Response
{ "message": "verification code sent to new email address" }
```

*Confirm email change (`action = email_verify`) — step 2*
```json
// Request
{ "action": "email_verify", "newEmail": "newalice@example.com", "otp": "482910" }

// Response
{ "message": "email updated successfully" }
```

*Update password (`action = password`)*
```json
// Request
{
  "action": "password",
  "currentPassword": "secret123",
  "newPassword": "newSecret456",
  "confirmPassword": "newSecret456"
}

// Response
{ "message": "password updated successfully" }
```

**POST `/login`**
```json
// Request
{ "number": 123456, "password": "password123" }

// Response
{ "token": "<jwt>", "number": 123456 }
```

**POST `/account`** *(Step 1 of registration)*
```json
// Request
{ "firstName": "Jane", "lastName": "Doe", "email": "jane@example.com", "password": "secret123" }

// Response
{ "message": "verification code sent to email" }
```

**POST `/account/verification`** *(Step 2 of registration)*
```json
// Request
{ "code": "482910" }

// Response — full Account object
{
  "id": 1,
  "firstName": "Jane",
  "lastName": "Doe",
  "number": 482910,
  "email": "jane@example.com",
  "balance": 0,
  "createdAt": "2024-01-15T10:30:00Z"
}
```

**POST `/transfer`** *(Authorization header required)*
```json
// Request
{ "toAccount": 567890, "amount": 500 }

// Response
{
  "status": "transfer successful",
  "transactionId": "TXN-000001234567",
  "transferredAt": "2024-01-15T10:35:00Z"
}
```

**GET `/account/transactions`** *(query parameters optional)*
```
GET /account/transactions?limit=10&offset=0&type=transfer&month=january
Authorization: <jwt>
```

**POST `/account/{id}/offer`** *(Authorization header required)*
```json
// Request
{ "couponCode": "OFFER1000" }

// Response
{ "status": "offer applied, 1000 added to your balance" }
```

---

## Getting Started

### Prerequisites

- **Go** 1.21 or later
- **PostgreSQL** running and accessible (default: `localhost:5432`)
- **`JWT_SECRET`** environment variable set for signing tokens
- *(Optional)* Gmail SMTP credentials for email verification and transfer notifications

### Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/Ahnaf-2210007/Online-Banking-System.git
   cd Online-Banking-System
   ```

2. Install Go dependencies:
   ```bash
   go mod download
   ```

3. Create a `.env` file in the project root (or export the variables into your shell):
   ```env
   JWT_SECRET=your-secret-key

   # PostgreSQL — individual fields (used when DATABASE_URL is not set)
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=postgres
   DB_PASSWORD=gobank
   DB_NAME=postgres

   # Or use a full DSN instead of the individual fields above:
   # DATABASE_URL=postgres://postgres:gobank@localhost:5432/postgres?sslmode=disable

   # Optional — Gmail SMTP for email verification and transfer notifications
   SMTP_EMAIL=your-gmail@gmail.com
   SMTP_PASSWORD=your-app-password
   ```

   > **Note:** When `SMTP_EMAIL` or `SMTP_PASSWORD` are not set, verification codes and transfer notifications are only logged to stdout. This is convenient for local development.

### Build & Run

```bash
# Build the binary
make build

# Build and run in one step
make run
```

The server starts on **`http://localhost:3000`**.

### Seed Test Data

Pass the `-seed` flag to create a sample account on startup:

```bash
./bin/gobank -seed
```

This creates an account for **John Doe** (`john.doe@example.com`) with password `password123`, which you can use immediately to test the `/login` endpoint.

### Run Tests

```bash
make test
```

---

## Authentication

Most write operations and all per-account reads require a valid JWT token in the `Authorization` HTTP header:

```
Authorization: <token>
```

Tokens are obtained from the `POST /login` endpoint and are valid for **24 hours**.

Routes protected by the `withJWTAuth` middleware (`GET /account/{id}`, `DELETE /account/{id}`, `POST /account/{id}/offer`) also verify that the JWT's embedded account number matches the requested account, preventing cross-account access.

The `POST /transfer`, `GET /account/transactions`, and `POST /account/update` endpoints also require a valid JWT but perform their own inline token validation to identify the caller's account.

---

## Configuration

| Setting | Environment Variable | Default / Notes |
|---|---|---|
| Server port | *(hardcoded)* | `:3000` |
| Full DB connection string | `DATABASE_URL` | Overrides individual `DB_*` variables when set |
| Database host | `DB_HOST` | `localhost` |
| Database port | `DB_PORT` | `5432` |
| Database user | `DB_USER` | `postgres` |
| Database password | `DB_PASSWORD` | `gobank` |
| Database name | `DB_NAME` | `postgres` |
| JWT signing secret | `JWT_SECRET` | Required |
| SMTP sender email | `SMTP_EMAIL` | Optional; if absent, emails are only logged to stdout |
| SMTP sender password | `SMTP_PASSWORD` | Optional; if absent, emails are only logged to stdout |
