# 🚀 GoCMS

A modern Content Management System built with **Go**, featuring a clean architecture, robust database management, and comprehensive testing.

## ✨ Features

- 🏗️ **Clean Architecture** - Well-organized layered structure (adapter, domain, usecase, repository)
- 🗄️ **Database Migrations** - Automated migrations using Goose
- 🧪 **Comprehensive Testing** - Unit tests for all critical components
- 🐳 **Docker Support** - Easy deployment with Docker Compose
- 🔄 **RESTful API** - HTTP handlers with proper routing

## 📋 Project Structure

```
goCMS/
├── cmd/                 # Application entry points
│   └── server/         # Server startup
├── internal/
│   ├── adapter/        # External integrations (HTTP handlers, routers)
│   ├── app/            # Application setup & initialization
│   ├── domain/         # Business domain entities
│   ├── repository/     # Data access layer
│   ├── usecase/        # Business logic & use cases
│   └── utils/          # Utilities (config, database)
├── migrations/         # Database migrations
└── docker-compose.yaml # Docker configuration
```

## 🛠️ Prerequisites

- Go 1.21+
- MySQL 8.0+
- Docker & Docker Compose (optional)

## 🚀 Quick Start

### 1. Setup Goose (Database Migration Tool)

Install Goose:
```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
# or on macOS
brew install goose
```

### 2. Configure Environment Variables

```bash
export GOOSE_DRIVER=mysql
export GOOSE_DBSTRING="huri:huri1234@tcp(localhost:3306)/goCMS"
export GOOSE_MIGRATION_DIR="migrations"
```

### 3. Database Migrations

```bash
# Create a new migration
goose create add_posts_tables sql

# Run all pending migrations
goose up

# Check migration status
goose status

# Rollback all migrations
goose reset
```

## 🐳 Docker Setup

Start the application with Docker Compose:

```bash
docker-compose up -d
```

## 🧪 Testing

Run all tests:

```bash
bash cicd/unit_test.sh
```

## 📝 License

MIT