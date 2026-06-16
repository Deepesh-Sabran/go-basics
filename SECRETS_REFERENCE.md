# 🔑 Secrets Reference Card

## Load Secrets in Code

```go
// Secrets are loaded automatically on app startup
// in cmd/main.go by calling:
config.LoadSecrets()

// Get secrets in your code:
jwtSecret := config.GetJWTSecret()        // Returns []byte
dbUrl := config.GetDatabaseURL()          // Returns string
redisAddr := config.GetRedisAddr()        // Returns string
isProd := config.IsProduction()           // Returns bool
```

## Environment Variables

| Variable | Required | Type | Example |
|----------|----------|------|---------|
| `JWT_SECRET` | ✅ | string | 32+ chars, random bytes |
| `DATABASE_URL` | ✅ | string | `user=app password=pwd dbname=app` |
| `REDIS_ADDR` | ❌ | string | `localhost:6379` |
| `ENVIRONMENT` | ❌ | string | `development` or `production` |

## Generate Secure Secrets

```bash
# JWT Secret (32+ random characters)
openssl rand -base64 32

# Database password
openssl rand -base64 16

# Redis password
openssl rand -base64 16
```

## Development Setup

```bash
# 1. Copy template
cp .env.example .env

# 2. Edit with your values
# Generate strong JWT secret:
JWT_SECRET=$(openssl rand -base64 32)
echo "JWT_SECRET=$JWT_SECRET" >> .env

# 3. Update DATABASE_URL and REDIS_ADDR
# 4. Run
go run cmd/main.go
```

## Production Deployment

### Using Docker
```dockerfile
FROM golang:1.22 AS builder
WORKDIR /app
COPY . .
RUN go build -o app cmd/main.go

FROM alpine:3.19
RUN apk add --no-cache ca-certificates
COPY --from=builder /app/app .
# Set secrets as environment variables
ENV JWT_SECRET=${JWT_SECRET}
ENV DATABASE_URL=${DATABASE_URL}
ENV REDIS_ADDR=${REDIS_ADDR}
ENV ENVIRONMENT=production
CMD ["./app"]
```

### Using Docker Compose
```yaml
version: '3.8'
services:
  app:
    build: .
    environment:
      JWT_SECRET: ${JWT_SECRET}
      DATABASE_URL: ${DATABASE_URL}
      REDIS_ADDR: redis:6379
      ENVIRONMENT: production
    depends_on:
      - redis
  redis:
    image: redis:7-alpine
```

### Using Kubernetes
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: app-secrets
type: Opaque
stringData:
  jwt_secret: "your_jwt_secret_here"
  database_url: "postgresql://user:pwd@db:5432/app"
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: app
spec:
  template:
    spec:
      containers:
      - name: app
        image: myapp:latest
        env:
        - name: JWT_SECRET
          valueFrom:
            secretKeyRef:
              name: app-secrets
              key: jwt_secret
        - name: DATABASE_URL
          valueFrom:
            secretKeyRef:
              name: app-secrets
              key: database_url
        - name: REDIS_ADDR
          value: "redis:6379"
        - name: ENVIRONMENT
          value: "production"
```

### Using SystemD
```ini
[Service]
Type=simple
WorkingDirectory=/opt/app
ExecStart=/opt/app/app
Restart=on-failure

# Set secrets as environment variables
EnvironmentFile=/etc/app/secrets.env
Environment="ENVIRONMENT=production"

User=app
Group=app
```

## Accessing Secrets in Handlers/Services

```go
package handlers

import "github.com/Deepesh-Sabran/go-basics/internal/config"

func SomeHandler(w http.ResponseWriter, r *http.Request) {
    // Access secrets
    secret := config.GetJWTSecret()
    
    // Use for JWT operations
    token, _ := jwt.NewWithClaims(
        jwt.SigningMethodHS256, 
        claims,
    ).SignedString(secret)
}
```

## Rotating Secrets

### JWT Secret Rotation
```go
// Don't do this in production without a plan!
// Better approach: support multiple keys and deprecate old ones

// 1. Generate new secret
newSecret := "new_secret_from_vault"

// 2. Update ENVIRONMENT
export JWT_SECRET=newSecret

// 3. Restart app
# Deploy new version

// 4. Existing tokens with old secret still work briefly
// but new tokens use new secret
```

## Troubleshooting

| Error | Cause | Fix |
|-------|-------|-----|
| `JWT_SECRET is not set` | Missing env var | `export JWT_SECRET=...` or add to .env |
| `DATABASE_URL is not set` | Missing env var | `export DATABASE_URL=...` or add to .env |
| `Failed to connect to Redis` | Redis offline | Start Redis: `redis-server` |
| `Invalid token` | Secret mismatch | Same secret used for signing/verification |

## Security Checklist

- [ ] .env is in .gitignore (verify: `git check-ignore .env`)
- [ ] Never commit secrets to git
- [ ] Use different secrets per environment
- [ ] Rotate secrets periodically
- [ ] Use strong, random secrets (min 32 chars for JWT)
- [ ] Restrict access to secret files/environment
- [ ] Monitor secret access in logs
- [ ] Use secret management system in production (Vault, K8s Secrets, etc.)

## API Endpoints

Test your authentication setup:

```bash
# Login
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"name":"testuser","password":"testpass"}'

# Use returned token
curl -X GET http://localhost:8080/me \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"

# Refresh token
curl -X POST http://localhost:8080/refresh \
  -H "Content-Type: application/json" \
  -d '{"refresh_token":"YOUR_REFRESH_TOKEN"}'
```

---

**For complete details, see `SECRET_MANAGEMENT.md`**
