# Auth Service with RBAC

A production-grade authentication and authorization service built from scratch in Go. No third-party auth libraries — just clean architecture, JWT-based auth, role-based access control, Redis caching, background workers, and a full audit trail.

---

## Features

- **JWT Authentication** — secure token generation, validation, and expiry handling
- **Role-Based Access Control (RBAC)** — fine-grained permission checks via a dedicated middleware layer
- **Audit Logging** — every mutating action (who updated what, who deleted whom) is recorded in a dedicated PostgreSQL audit table with actor, target, action type, and timestamp
- **Redis Caching** — frequent auth lookups cached in Redis to reduce DB hits on every authenticated request
- **Background Workers** — async jobs like token cleanup connected to Redis queues; completed and pending jobs tracked separately
- **Pagination and Result Capping** — all list endpoints paginated with a hard cap to prevent unbounded queries
- **Separate DB and Response Structs** — sensitive fields like passwords never leak into API responses

---

## Architecture

Strict three-layer architecture — nothing leaks between layers:

```
Handler → Service → Repository
```

- **Handler** — parses and validates the incoming request
- **Service** — contains all business logic
- **Repository** — handles all database interaction

Auth and authorization are intentionally split into two separate middleware layers:

```
Request → Auth Middleware (JWT) → Permission Middleware (RBAC) → Handler
```

---

## Tech Stack

| Layer | Technology |
|---|---|
| Language | Go (net/http) |
| Database | PostgreSQL |
| Caching & Workers | Redis |
| Auth | JWT |

---

## Project Structure

```
.
├── cmd/
│   └── main.go               # Entry point
├── internal/
│   ├── handler/              # HTTP handlers — request parsing and response writing
│   ├── service/              # Business logic
│   ├── repository/           # Database queries
│   ├── middleware/
│   │   ├── auth.go           # JWT validation middleware
│   │   └── permission.go     # RBAC permission middleware
│   ├── model/
│   │   ├── db/               # DB structs (never exposed to client)
│   │   └── response/         # User-facing response structs
│   └── worker/               # Background workers connected to Redis
├── pkg/
│   └── redis/                # Redis client setup
├── go.mod
└── go.sum
```

---

## Getting Started

### Prerequisites

- Go 1.21+
- PostgreSQL
- Redis

### Setup

1. **Clone the repository**

```bash
git clone https://github.com/Deepesh-Sabran/your-repo-name.git
cd your-repo-name
```

2. **Create a `.env` file** in the root directory

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_NAME=auth_service

REDIS_ADDR=localhost:6379

JWT_SECRET=your_jwt_secret
JWT_EXPIRY=24h
```

3. **Install dependencies**

```bash
go mod tidy
```

4. **Run the server**

```bash
go run cmd/main.go
```

Server starts on `http://localhost:8080`

---

## API Overview

| Method | Endpoint | Description | Auth Required |
|---|---|---|---|
| POST | `/auth/register` | Register a new user | No |
| POST | `/auth/login` | Login and receive JWT | No |
| GET | `/users` | List all users (paginated) | Yes + Admin |
| PUT | `/users/:id` | Update a user | Yes + Admin |
| DELETE | `/users/:id` | Delete a user | Yes + Admin |
| GET | `/audit` | View audit logs | Yes + Admin |

---

## Key Design Decisions

**Why split auth and permission into two middleware layers?**
Keeping them separate means each layer has one responsibility. Adding a new role or permission never requires touching the JWT validation logic.

**Why separate DB and response structs?**
The struct saved to PostgreSQL should never be the struct returned to the client. Passwords, internal IDs, and sensitive fields stay inside the DB model and never accidentally leak into API responses.

**Why audit logging?**
Most side projects skip it. Every real production system needs it. Knowing who did what and when is non-negotiable for any multi-role system.

---

## Author

**Deepesh Sabran**
[LinkedIn](https://www.linkedin.com/in/deepesh-sabran-13b733201) · [GitHub](https://github.com/Deepesh-Sabran)