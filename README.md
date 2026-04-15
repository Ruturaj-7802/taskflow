# TaskFlow

## 1. Overview
Task management REST API. Go + chi + PostgreSQL. JWT authentication.
Backend-only submission — includes Postman collection instead of frontend.

Tech stack: Go 1.22, chi, pgx/v5, golang-migrate, PostgreSQL 16, Docker

## 2. Architecture Decisions
- **chi over gin**: handlers stay as standard http.HandlerFunc — portable and testable
- **pgx/v5 over GORM**: raw SQL, every query is explicit and reviewable, no ORM magic
- **golang-migrate**: versioned SQL files, tracked in git, auto-run on startup
- **Layered architecture** (handler → service → repository): each layer has one job
- **401 vs 403**: 401 = no valid token, 403 = valid token but wrong permissions — not conflated
- **Postgres ENUMs** for status/priority: enforced at DB level, invalid values rejected
- **TIMESTAMPTZ**: all timestamps stored in UTC

## 3. Running Locally
```bash
git clone https://github.com/Ruturaj-7802/taskflow
cd taskflow
cp .env.example .env
docker compose up
# API available at http://localhost:8080
```

## 4. Running Migrations
Migrations run automatically on container start via golang-migrate in main.go.
No manual steps required.

## 5. Test Credentials
Email:    test@example.com
Password: password123

## 6. API Reference

## 7. What I'd Do With More Time
- Wrap project delete in a DB transaction (currently two queries — tasks then project)
- Add refresh tokens with revocation table
- Rate limit the login endpoint against brute force
- Integration tests using testcontainers-go
- Add updated_at trigger in Postgres instead of setting it in Go code