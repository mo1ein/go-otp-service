# OTP Auth Service

[![Go](https://img.shields.io/badge/go-1.21-blue)](https://golang.org)  ![status](https://img.shields.io/badge/status-ready-brightgreen)

> **otp-auth-api** — a compact, production-minded Go service that implements OTP-based login/registration, lightweight user management, JWT authentication, rate-limiting and OpenAPI documentation. Designed to be easy to run for reviewers while following real-world patterns.

---

## Features

- **OTP login/registration**
  - `POST /api/auth/request-otp`: generate a 6-digit OTP (printed to server logs), valid for **2 minutes**.
  - `POST /api/auth/verify`: validate OTP; if user not exists → register, else login. Returns **JWT**.
- **Rate limiting**: max **3 OTP requests per phone** per **10 minutes**.
- **User management** (JWT protected):
  - `GET /api/me`
  - `GET /api/users/{id}`
  - `GET /api/users?search=&page=&page_size=` (pagination + search by phone substring)
- **Storage choice**: In-memory for simplicity and speed in take-home tasks. No external DB required.

## Why In-Memory?

In-memory storage keeps the code and architecture clear while meeting all requirements. It removes infra complexity (migrations, credentials), speeds up evaluation, and still demonstrates proper layering (handlers → services → repositories). If a DB is required later, the repository interfaces make it straightforward to swap implementations.


## Architecture & rationale

The codebase follows a simple, testable layering:

- `internal/handlers` — HTTP routing & request/response handling
- `internal/service` — application rules (OTP lifecycle, rate limiting, JWT issuance)
- `internal/repo` — storage abstractions (in-memory by default)
- `internal/token` — JWT signing & parsing
- `cmd/server` — composition root (wires dependencies)

This separation keeps responsibilities small and makes it trivial to swap the storage implementation for a real DB.

---

## Data model

```json
{
  "id": "uuid",
  "phone": "+989123456789",
  "registered_at": "2025-09-06T12:34:56Z"
}
```

---

```mermaid
sequenceDiagram
    participant Client
    participant API
    participant OTPStore
    participant UserRepo
    participant JWT

    Client->>API: POST /api/auth/request-otp { phone }
    API->>OTPStore: generate & store OTP (2m TTL)
    API-->>Client: 200 { status: ok }

    Note over API,OTPStore: OTP printed to server log for demo

    Client->>API: POST /api/auth/verify { phone, otp }
    API->>OTPStore: verify OTP (one-time)
    alt OTP valid
        API->>UserRepo: get by phone
        alt user not exists
            UserRepo->>UserRepo: create user (registered_at)
        end
        API->>JWT: sign token (sub = userID)
        API-->>Client: 200 { token }
    else invalid
        API-->>Client: 401
    end

    Client->>API: GET /api/me (Authorization: Bearer)
    API->>JWT: parse token
    API->>UserRepo: get by ID
    API-->>Client: 200 { user }
```

---

## Quickstart

### Run locally

```bash
make up
make migration-up
go run cmd/server/main.go
```

The API listens on `http://localhost:8080`.


### Run with Docker


```bash
make up
```

- API: `http://localhost:8080`

---