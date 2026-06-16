# Secret Management Guide

## Overview
This project uses environment variables to manage secrets in a production-grade manner. The `secrets.go` package provides a centralized configuration system that loads secrets from `.env` files (development) or environment variables (production).

## Key Features
✅ **Centralized Secret Management** - All secrets managed in one place
✅ **Environment-Aware** - Different handling for dev vs production
✅ **Automatic Validation** - Fails fast if required secrets are missing
✅ **No Hardcoding** - All sensitive data in environment variables
✅ **Connection Pooling** - Optimized for production with connection limits

---

## Setup Instructions

### 1. Copy Example Config
```bash
cp .env.example .env
```

### 2. Update `.env` with Your Secrets
```env
# Generate a strong JWT secret (minimum 32 characters)
JWT_SECRET=your_super_secret_jwt_key_here_min_32_chars

# Database URL (PostgreSQL format)
DATABASE_URL=user=apple password=yourpassword dbname=goapp host=localhost port=5432 sslmode=disable

# Redis address
REDIS_ADDR=localhost:6379

# Environment mode
ENVIRONMENT=development
```

### 3. Start Application
```bash
go run cmd/main.go
```

The app will automatically load secrets from `.env` file.

---

## Production Deployment

### Important: Production doesn't use .env files

In production, **DO NOT use .env files**. Instead, set environment variables directly:

**Using Docker:**
```dockerfile
ENV JWT_SECRET=your_production_secret
ENV DATABASE_URL=postgresql://user:pwd@prod-db:5432/goapp
ENV REDIS_ADDR=prod-redis:6379
ENV ENVIRONMENT=production
```

**Using Environment Variables:**
```bash
export JWT_SECRET="your_production_secret"
export DATABASE_URL="postgresql://user:pwd@prod-db:5432/goapp"
export REDIS_ADDR="prod-redis:6379"
export ENVIRONMENT="production"

./app
```

**Using systemd service file:**
```ini
[Service]
Environment="JWT_SECRET=your_production_secret"
Environment="DATABASE_URL=postgresql://user:pwd@prod-db:5432/goapp"
Environment="REDIS_ADDR=prod-redis:6379"
Environment="ENVIRONMENT=production"
ExecStart=/path/to/app
```

**AWS Secrets Manager / Other Secret Vaults:**
```bash
# Fetch secrets from vault and set as env vars
export JWT_SECRET=$(aws secretsmanager get-secret-value --secret-id prod/jwt-secret --query SecretString --output text)
export DATABASE_URL=$(aws secretsmanager get-secret-value --secret-id prod/db-url --query SecretString --output text)
export REDIS_ADDR=$(aws secretsmanager get-secret-value --secret-id prod/redis-addr --query SecretString --output text)
export ENVIRONMENT="production"

./app
```

---

## Environment Variables Reference

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `JWT_SECRET` | ✅ Yes | - | Secret key for JWT signing (min 32 chars) |
| `DATABASE_URL` | ✅ Yes | - | PostgreSQL connection URL |
| `REDIS_ADDR` | ❌ No | `localhost:6379` | Redis server address |
| `ENVIRONMENT` | ❌ No | `development` | `development` or `production` |

---

## How Secrets Are Used

### 1. JWT Token Signing/Verification
```go
// Automatically loaded and used in:
// - middleware/auth.go: Token validation
// - services/user_service.go: Token generation
secret := config.GetJWTSecret()  // Returns []byte
```

### 2. Database Connections
```go
// Loaded with connection pooling for production:
// - 25 max open connections (production)
// - 5 idle connections (production)
connStr := config.GetDatabaseURL()
```

### 3. Redis Connections
```go
// Loaded with optimized timeouts:
addr := config.GetRedisAddr()
```

---

## Security Best Practices

### ✅ DO's
- [ ] Generate strong JWT secrets: `openssl rand -base64 32`
- [ ] Use different secrets for each environment
- [ ] Keep `.env` file out of version control (already in `.gitignore`)
- [ ] Rotate secrets periodically
- [ ] Use HTTPS in production
- [ ] Set `ENVIRONMENT=production` in production

### ❌ DON'Ts
- [ ] Don't commit `.env` file to git
- [ ] Don't use weak secrets
- [ ] Don't hardcode secrets in code
- [ ] Don't share secrets via Slack/Email
- [ ] Don't use `.env` files in production

---

## Generating Secure Secrets

**JWT Secret (minimum 32 characters):**
```bash
# Using OpenSSL
openssl rand -base64 32

# Using Go
go run -c 'package main; import ("crypto/rand"; "fmt"; "encoding/base64"); func main() { b := make([]byte, 32); rand.Read(b); fmt.Println(base64.StdEncoding.EncodeToString(b)) }'
```

---

## Troubleshooting

### Error: "JWT_SECRET environment variable is not set"
- Make sure `.env` file exists and contains `JWT_SECRET`
- In production, set the environment variable: `export JWT_SECRET=...`

### Error: "DATABASE_URL environment variable is not set"
- Update `.env` file with correct PostgreSQL connection URL
- Check database credentials and host

### Error: "Failed to connect to Redis"
- Ensure Redis is running: `redis-cli ping`
- Update `REDIS_ADDR` if Redis is on different host/port
- Check firewall rules

### App runs fine in dev but fails in production
- Verify `ENVIRONMENT=production` is set
- Check all required env vars are set (use `env | grep -E "JWT|DATABASE|REDIS"`)
- Review connection pool settings in `db.go`

---

## Implementation Details

### Secrets Loading Flow
```
1. Check if ENVIRONMENT != "production"
2. If dev, load .env file (optional)
3. Validate all required secrets exist
4. Store in AppSecrets singleton
5. Ready to use via config.GetJWTSecret(), etc.
```

### Files Modified
- `internal/config/secrets.go` - NEW: Centralized secret management
- `internal/config/db.go` - UPDATED: Uses secrets, adds connection pooling
- `internal/config/redis.go` - UPDATED: Uses secrets, adds timeouts
- `internal/middleware/auth.go` - UPDATED: Uses config.GetJWTSecret()
- `internal/services/user_service.go` - UPDATED: Uses config.GetJWTSecret()
- `cmd/main.go` - UPDATED: Calls config.LoadSecrets() first

---

## Next Steps (Optional Enhancements)

1. **Secret Rotation:** Implement periodic JWT secret rotation
2. **Audit Logging:** Log all secret access attempts
3. **Encryption at Rest:** Encrypt sensitive data in database
4. **Rate Limiting:** Add rate limiting to auth endpoints
5. **2FA:** Implement two-factor authentication
6. **Security Headers:** Add CORS and security headers middleware
