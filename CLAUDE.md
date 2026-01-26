# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Gophish is an open-source phishing toolkit designed for security awareness training and penetration testing. This fork adds e-learning course functionality for security awareness training.

## Build & Run Commands

```bash
# Build the Go binary
go build

# Run the server (default: all modes)
./gophish

# Run in specific mode
./gophish --mode admin    # Admin interface only
./gophish --mode phish    # Phishing server only
./gophish --config /path/to/config.json  # Custom config path

# Run tests
go test ./...

# Run tests for a specific package
go test ./models/...
go test ./controllers/...

# Build frontend assets (requires npm/yarn)
npm install  # or yarn
gulp build   # Builds vendor.js, app scripts, and CSS
gulp scripts # JS only
gulp styles  # CSS only
```

## Architecture

### Two-Server Design
Gophish runs two separate HTTP servers:

1. **AdminServer** (`controllers/route.go`) - Port 3333 (default)
   - Dashboard, campaign management, user management
   - REST API under `/api/` prefix
   - Session-based authentication with CSRF protection
   - Routes require `mid.RequireLogin` middleware

2. **PhishingServer** (`controllers/phish.go`) - Port 80 (default)
   - Handles phishing campaign events (email opens, link clicks, form submissions)
   - Tracks results via unique RID (Result ID) in URLs
   - Serves landing pages to targets

### Key Packages

- **models/** - Database models using GORM ORM. Core entities: Campaign, Group, Template, Page, Result, User, Course
- **controllers/api/** - REST API handlers. `server.go` registers all API routes
- **worker/** - Background worker that polls for scheduled campaign emails and sends them via the mailer
- **mailer/** - Email sending logic with SMTP support
- **middleware/** - Authentication, CSRF, rate limiting, RBAC permission checks
- **auth/** - Password hashing, API key generation, session management

### Database
- Supports SQLite3 (default) and MySQL
- Migrations in `db/db_sqlite3/migrations/` and `db/db_mysql/migrations/`
- Uses goose for migrations, automatically run on startup via `models.Setup()`

### Configuration
- `config.json` - Main configuration (server addresses, TLS, database)
- `.env` - E-learning notification SMTP settings (optional feature)

### Frontend
- Templates in `templates/` (Go html/template)
- Static assets in `static/` (JS, CSS, images)
- Gulp builds/minifies JS and CSS to `static/js/dist/` and `static/css/dist/`

### E-Learning Extension (feature/elearning branch)
This fork adds security awareness training courses:
- **models/course.go** - Course, Module, Quiz, Enrollment models
- **controllers/api/course.go** - Course management API
- **controllers/api/user_dashboard.go** - User-facing course APIs
- Routes: `/courses`, `/user/courses`, `/api/courses/`, `/api/user/courses/`

### Role-Based Access Control
- Admin (RoleID=1): Full access, can impersonate users
- User (RoleID=2): Limited to personal dashboard and assigned courses
- Permission checks via `mid.RequirePermission(models.PermissionModifySystem)`

### Request Flow
1. Request hits AdminServer or PhishingServer
2. Middleware chain: ProxyHeaders -> GZIP -> Logging -> CSRF -> Auth -> Handler
3. API routes require API key header (`Authorization`) or session auth
4. User-specific endpoints (`/api/user/*`) use session auth, not API keys
