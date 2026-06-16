# 🚀 Quick Start: Secret Management

## 5-Minute Setup

### 1️⃣ Create Your `.env` File

```bash
# Copy the example
cp .env.example .env

# Edit with your actual secrets
nano .env
```

Update these values:

```env
JWT_SECRET=your_super_secret_key_32_chars_minimum
DATABASE_URL=user=apple password=yourdb dbname=goapp sslmode=disable
REDIS_ADDR=localhost:6379
ENVIRONMENT=development
```

### 2️⃣ Verify Setup

```bash
# Check .env file exists
ls -la .env

# Build the project
go build -o app cmd/main.go

# Run the app
./app
```

You should see:

```
✅ Secrets loaded successfully
✅ Successfully connected to DB !!
✅ Successfully connected to Redis !!
Server listens on port:8080
```

---

## 📋 Checklist

- [ ] `.env` file created and filled with real values
- [ ] `.env` is in `.gitignore` (verify: `git check-ignore .env`)
- [ ] `JWT_SECRET` is at least 32 characters
- [ ] Database URL is correct
- [ ] Redis is running on specified address
- [ ] Application builds without errors
- [ ] Application starts successfully
- [ ] Tokens are signed/verified correctly

---

## 🔐 Production Deployment

Set these in your production environment:

```bash
export JWT_SECRET="your_production_secret_from_vault"
export DATABASE_URL="postgresql://prod:pass@prod-db:5432/goapp"
export REDIS_ADDR="prod-redis:6379"
export ENVIRONMENT="production"

./app
```

**Never use `.env` files in production!**

---

## 🆘 Common Issues

**"JWT_SECRET is not set"**
→ Make sure `.env` file exists and contains `JWT_SECRET=...`

**"DB connection error"**
→ Check `DATABASE_URL` is correct and database is running

**"Failed to connect to Redis"**
→ Verify Redis is running: `redis-cli ping`

**See `SECRET_MANAGEMENT.md` for complete guide →**
