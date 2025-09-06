# OTP Auth Service (Go)

A clean, minimal backend implementing OTP-based login/registration, basic user management, JWT auth, rate limiting, and OpenAPI docs.

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

## Why In-Memory? (DB Justification)

In-memory storage keeps the code and architecture clear while meeting all requirements. It removes infra complexity (migrations, credentials), speeds up evaluation, and still demonstrates proper layering (handlers → services → repositories). If a DB is required later, the repository interfaces make it straightforward to swap implementations.

## Getting Started
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