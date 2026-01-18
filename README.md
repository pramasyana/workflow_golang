# Workflow Approval System

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.19+-00ADD8?style=for-the-badge&logo=go" alt="Go Version">
  <img src="https://img.shields.io/badge/GoFiber-v2-blue?style=for-the-badge&logo=go" alt="GoFiber">
  <img src="https://img.shields.io/badge/MySQL-8.0-4479A1?style=for-the-badge&logo=mysql" alt="MySQL">
  <img src="https://img.shields.io/badge/GORM-v2-black?style=for-the-badge" alt="GORM">
</p>

<p align="center">
  REST API Service built with Golang for managing workflow approval systems with double-layer concurrency control.
</p>

---

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Project Structure](#project-structure)
- [Technology Stack](#technology-stack)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
  - [Configuration](#configuration)
  - [Running the Application](#running-the-application)
- [Database Schema](#database-schema)
- [API Endpoints](#api-endpoints)
  - [Authentication](#authentication)
  - [Actors](#actors)
  - [Workflows](#workflows)
  - [Workflow Steps](#workflow-steps)
  - [Requests](#requests)
  - [Users](#users)
  - [Profile](#profile)
- [Architecture](#architecture)
  - [Clean Architecture](#clean-architecture)
  - [Concurrency Control](#concurrency-control)
  - [Authentication Flow](#authentication-flow)
- [Workflow Process](#workflow-process)
- [Testing](#testing)
- [Docker](#docker)
- [License](#license)

---

## Overview

Sistem Workflow Approval ini adalah REST API service yang dibangun dengan Golang untuk mengelola proses persetujuan workflow. Sistem ini dirancang untuk menangani berbagai jenis pengajuan seperti dokumen, purchase requests, atau submission internal.

**Konsep Utama:**
- **Actor/Role** - Peran atau jabatan yang dapat melakukan approval (Manager, Director, Finance, dll)
- **Workflow** - Template workflow yang terdiri dari beberapa langkah approval
- **Workflow Step** - Langkah-langkah approval dengan kondisi dan actor yang berbeda
- **Request** - Pengajuan yang mengikuti workflow tertentu dengan approval bertahap

---

## Features

- **RESTful API** dengan GoFiber framework
- **JWT Authentication** dengan middleware protection
- **Actor-based Authorization** - Approver harus memiliki actor_id yang sesuai dengan step
- **Workflow Management** dengan multiple approval steps
- **Request Submission & Approval** dengan proses bertahap
- **Double-layer Concurrency Control** (Mutex + SELECT FOR UPDATE)
- **Pagination & Filtering** untuk list endpoints
- **Approval History** - Pelacakan lengkap siapa yang approve/reject
- **MySQL Database** dengan GORM ORM
- **Clean Architecture** design pattern
- **Environment Variables** support dengan YAML config
- **Default Admin Account** untuk easy setup

---

## Project Structure

```
workflow-golang/
├── config/
│   ├── config.yaml          # Configuration file
│   └── config.go            # Configuration loader
├── docs/
│   ├── openapi.yaml         # Swagger/OpenAPI documentation
│   └── postman_collection.json
├── framework/
│   ├── middleware/
│   │   └── jwt.go           # JWT authentication middleware
│   └── router/
│       └── router.go        # Route configuration
├── package/
│   ├── actor/               # Actor/Role management
│   │   ├── domain/          # Actor entity
│   │   ├── usecase/         # Actor business logic
│   │   ├── handler/         # HTTP handlers
│   │   ├── repository/      # Database operations
│   │   └── ports/           # Interface definitions
│   ├── auth/                # Authentication
│   │   ├── domain/
│   │   ├── usecase/
│   │   ├── handler/
│   │   ├── repository/
│   │   └── ports/
│   ├── user/                # User management
│   │   ├── domain/          # User entity
│   │   ├── usecase/
│   │   ├── handler/
│   │   ├── repository/
│   │   └── ports/
│   ├── workflow/            # Workflow management
│   │   ├── domain/          # Workflow entity
│   │   ├── usecase/
│   │   ├── handler/
│   │   ├── repository/
│   │   └── ports/
│   ├── workflow_step/       # Workflow steps
│   │   ├── domain/          # WorkflowStep entity
│   │   ├── usecase/
│   │   ├── handler/
│   │   ├── repository/
│   │   └── ports/
│   ├── request/             # Request & approval
│   │   ├── domain/          # Request entity
│   │   ├── usecase/         # Approval logic with mutex
│   │   ├── handler/
│   │   ├── repository/      # SELECT FOR UPDATE
│   │   └── ports/
│   ├── approval_history/    # Approval tracking
│       ├── domain/          # ApprovalHistory entity
│       ├── usecase/
│       ├── handler/
│       ├── repository/
│       └── ports/
├── utils/
│   ├── utils.go             # Utility functions
│   └── jwthelper/           # JWT helper
├── main.go
├── Dockerfile
├── go.mod
├── go.sum
└── README.md
```

---

## Technology Stack

| Component | Technology |
|-----------|------------|
| Language | Go 1.19+ |
| Web Framework | GoFiber v2 |
| ORM | GORM v2 |
| Database | MySQL 8.0 |
| Authentication | JWT (HMAC-SHA256) |
| Password Hashing | bcrypt |
| Concurrency | golang.org/x/sync (sync.Map + Mutex) |
| Config | gopkg.in/yaml.v3 + godotenv |

---

## Getting Started

### Prerequisites

- Go 1.19 atau lebih tinggi
- MySQL 8.0 atau lebih tinggi

### Installation

1. Clone repository:
```bash
git clone <repository-url>
cd workflow-golang
```

2. Install dependencies:
```bash
go mod tidy
```

3. Configure database (see Configuration section)

### Configuration

#### Using config.yaml (Default)

Edit `config/config.yaml` untuk konfigurasi aplikasi:

```yaml
# Application Configuration
app:
  host: "0.0.0.0"
  port: 8080
  name: "Workflow Approval System"
  env: "development"

# Database Configuration (MySQL)
database:
  host: "localhost"
  port: 3306
  username: "root"
  password: ""
  name: "workflow_approval"
  charset: "utf8"
  parse_time: true
  max_open_conns: 25
  max_idle_conns: 5
  conn_max_lifetime: 300

# JWT Configuration
jwt:
  secret: "your-super-secret-key-change-in-production"
  expiration: 24
  issuer: "workflow-approval-system"

# Logging Configuration
logging:
  level: "debug"
  format: "json"
```

#### Default Admin Account

Setelah migration berjalan, admin default akan otomatis dibuat:

| Field | Value |
|-------|-------|
| Email | administrator@gmail.com |
| Password | password123 |

### Running the Application

#### Option 1: Using Go (Development)

Jalankan aplikasi langsung menggunakan Go:

```bash
# Run dengan go run
go run main.go

# Atau build terlebih dahulu
go build -o workflow-approval .
./workflow-approval
```

Server akan berjalan di `http://localhost:8080`

**Catatan:** Pastikan MySQL sudah berjalan dan terkonfigurasi di `config/config.yaml` atau environment variables.

#### Option 2: Using Docker Compose (Production/Development)

Cara termudah untuk menjalankan aplikasi dengan database adalah menggunakan Docker Compose:

```bash
# Build dan jalankan semua services
docker-compose up -d --build
```

Server akan berjalan di `http://localhost:8080`

### Database Migration

Aplikasi akan otomatis menjalankan migration saat startup menggunakan GORM. Tables yang dibuat:

- `actors` - Actor/Role table
- `users` - User accounts
- `workflows` - Workflow templates
- `workflow_steps` - Workflow steps
- `requests` - Approval requests
- `approval_history` - Approval/rejection history

---

## Database Schema

### Actors (Peran/Jabatan)

| Column | Type | Description |
|--------|------|-------------|
| id | VARCHAR(36) | Primary key (UUID) |
| name | VARCHAR(255) | Actor name (e.g., "Manager", "Director") |
| code | VARCHAR(50) | Unique code (e.g., "manager", "director") |
| created_at | DATETIME | Creation time |
| updated_at | DATETIME | Last update time |

### Users

| Column | Type | Description |
|--------|------|-------------|
| id | VARCHAR(36) | Primary key (UUID) |
| email | VARCHAR(255) | Unique, not null |
| password | VARCHAR(255) | Hashed password (bcrypt) |
| name | VARCHAR(255) | User's full name |
| is_admin | BOOLEAN | Admin flag (default: FALSE) |
| actor_id | VARCHAR(36) | Optional actor association |
| created_at | DATETIME | Creation time |
| updated_at | DATETIME | Last update time |

### Workflows

| Column | Type | Description |
|--------|------|-------------|
| id | VARCHAR(36) | Primary key (UUID) |
| name | VARCHAR(255) | Workflow name |
| created_at | DATETIME | Creation time |
| updated_at | DATETIME | Last update time |

### Workflow Steps

| Column | Type | Description |
|--------|------|-------------|
| id | VARCHAR(36) | Primary key (UUID) |
| workflow_id | VARCHAR(36) | Foreign key to workflow |
| level | INT | Step level (unique per workflow, starting from 1) |
| actor_id | VARCHAR(36) | Required actor for this step |
| conditions | JSON | Step conditions (min_amount, max_amount, roles) |
| description | VARCHAR(500) | Optional step description |
| created_at | DATETIME | Creation time |
| updated_at | DATETIME | Last update time |

### Requests

| Column | Type | Description |
|--------|------|-------------|
| id | VARCHAR(36) | Primary key (UUID) |
| workflow_id | VARCHAR(36) | Foreign key to workflow |
| requester_id | VARCHAR(36) | Foreign key to user |
| current_step | INT | Current approval step (default: 1) |
| status | VARCHAR(20) | PENDING, APPROVED, REJECTED |
| amount | DECIMAL(15,2) | Request amount |
| title | VARCHAR(255) | Request title |
| description | TEXT | Request description |
| version | INT | Optimistic locking version |
| created_at | DATETIME | Creation time |
| updated_at | DATETIME | Last update time |

### Approval History

| Column | Type | Description |
|--------|------|-------------|
| id | VARCHAR(36) | Primary key (UUID) |
| request_id | VARCHAR(36) | Foreign key to request |
| workflow_id | VARCHAR(36) | Foreign key to workflow |
| step_level | INT | Step level of action |
| actor_id | VARCHAR(36) | Actor who performed action |
| user_id | VARCHAR(36) | User who performed action |
| action | VARCHAR(20) | APPROVE or REJECT |
| comment | TEXT | Comment or rejection reason |
| created_at | DATETIME | Creation time |

---

## API Endpoints

### Authentication

#### Login

**Endpoint:** `POST /auth/login`

Login menggunakan JSON body dengan email dan password.

```http
POST /auth/login
Content-Type: application/json

{
    "email": "user@example.com",
    "password": "password123"
}
```

**Response (200 OK):**
```json
{
    "success": true,
    "data": {
        "user": {
            "id": "550e8400-e29b-41d4-a716-446655440000",
            "email": "user@example.com",
            "name": "John Doe",
            "role": "user",
            "is_admin": false,
            "actor_id": null
        },
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
    },
    "error": null
}
```

#### Refresh Token

**Endpoint:** `POST /auth/refresh`

Generate token baru menggunakan token yang masih valid.

```http
POST /auth/refresh
Authorization: Bearer <token>
```

#### Logout

**Endpoint:** `POST /auth/logout`

Logout user (invalidasi token).

```http
POST /auth/logout
Authorization: Bearer <token>
```

---

### Actors

Actors merepresentasikan peran/jabatan yang dapat melakukan approval (Manager, Director, Finance, dll).

#### List All Actors

```http
GET /api/actors
Authorization: Bearer <token>
```

#### Create Actor

```http
POST /api/actors
Authorization: Bearer <token>
Content-Type: application/json

{
    "name": "Manager",
    "code": "manager"
}
```

#### Get Actor

```http
GET /api/actors/{id}
Authorization: Bearer <token>
```

#### Update Actor

```http
PUT /api/actors/{id}
Authorization: Bearer <token>
Content-Type: application/json

{
    "name": "Senior Manager"
}
```

#### Delete Actor

```http
DELETE /api/actors/{id}
Authorization: Bearer <token>
```

---

### Workflows

#### List Workflows (with pagination)

```http
GET /api/workflows?page=1&limit=10
Authorization: Bearer <token>
```

#### Create Workflow

```http
POST /api/workflows
Authorization: Bearer <token>
Content-Type: application/json

{
    "name": "Purchase Approval"
}
```

#### Get Workflow

```http
GET /api/workflows/{id}
Authorization: Bearer <token>
```

#### Update Workflow

```http
PUT /api/workflows/{id}
Authorization: Bearer <token>
Content-Type: application/json

{
    "name": "Updated Purchase Approval"
}
```

#### Delete Workflow

```http
DELETE /api/workflows/{id}
Authorization: Bearer <token>
```

---

### Workflow Steps

Workflow steps mendefinisikan langkah-langkah approval dalam sebuah workflow dengan kondisi dan actor yang berbeda.

#### List Workflow Steps

```http
GET /api/workflows/{id}/steps
Authorization: Bearer <token>
```

#### Create Workflow Step

```http
POST /api/workflows/{id}/steps
Authorization: Bearer <token>
Content-Type: application/json

{
    "level": 1,
    "actor_id": "550e8400-e29b-41d4-a716-446655440001",
    "conditions": {
        "min_amount": 1000000
    },
    "description": "Manager approval for purchases above 1M"
}
```

#### Get Workflow Step

```http
GET /api/workflows/{id}/steps/{stepId}
Authorization: Bearer <token>
```

#### Update Workflow Step

```http
PUT /api/workflows/{id}/steps/{stepId}
Authorization: Bearer <token>
Content-Type: application/json

{
    "level": 2,
    "actor_id": "550e8400-e29b-41d4-a716-446655440002",
    "conditions": {
        "min_amount": 5000000
    }
}
```

#### Delete Workflow Step

```http
DELETE /api/workflows/{id}/steps/{stepId}
Authorization: Bearer <token>
```

---

### Requests

Requests adalah pengajuan yang mengikuti workflow tertentu dengan proses approval bertahap.

#### List Requests (with pagination & filtering)

```http
GET /api/requests?page=1&limit=10&status=PENDING
Authorization: Bearer <token>
```

#### Create Request

```http
POST /api/requests
Authorization: Bearer <token>
Content-Type: application/json

{
    "workflow_id": "550e8400-e29b-41d4-a716-446655440001",
    "amount": 1500000,
    "title": "Office Supplies Purchase",
    "description": "Monthly office supplies for Q1"
}
```

#### Get Request

```http
GET /api/requests/{id}
Authorization: Bearer <token>
```

#### Update Request

Hanya bisa dilakukan jika status masih PENDING.

```http
PUT /api/requests/{id}
Authorization: Bearer <token>
Content-Type: application/json

{
    "amount": 2000000,
    "title": "Updated Office Supplies Purchase"
}
```

#### Approve Request

Approve request saat ini dan lanjut ke step berikutnya. Approver harus memiliki actor_id yang sesuai dengan step saat ini.

```http
POST /api/requests/{id}/approve
Authorization: Bearer <token>
```

**Business Rules:**
- Request harus dalam status PENDING
- Amount harus memenuhi min_amount condition dari step saat ini
- User's actor_id harus sesuai dengan step's actor_id (kecuali admin)
- Jika tidak ada step berikutnya, status menjadi APPROVED

#### Reject Request

Reject request dengan alasan.

```http
POST /api/requests/{id}/reject
Authorization: Bearer <token>
Content-Type: application/json

{
    "reason": "Insufficient documentation"
}
```

**Business Rules:**
- Request harus dalam status PENDING
- Reason harus diisi
- User's actor_id harus sesuai dengan step's actor_id (kecuali admin)

#### Get Request Approval History

Mendapatkan riwayat lengkap approval/rejection untuk sebuah request.

```http
GET /api/requests/{id}/history
Authorization: Bearer <token>
```

**Response:**
```json
{
    "success": true,
    "data": {
        "request_id": "550e8400-e29b-41d4-a716-446655440002",
        "history": [
            {
                "id": "550e8400-e29b-41d4-a716-446655440003",
                "step_level": 1,
                "actor_id": "550e8400-e29b-41d4-a716-446655440001",
                "user_id": "550e8400-e29b-41d4-a716-446655440004",
                "action": "APPROVE",
                "comment": "",
                "created_at": "2025-01-15T10:30:00Z"
            }
        ]
    },
    "error": null
}
```

#### Delete Request

```http
DELETE /api/requests/{id}
Authorization: Bearer <token>
```

---

### Users

#### Create User

```http
POST /api/users
Authorization: Bearer <token>
Content-Type: application/json

{
    "email": "newuser@example.com",
    "password": "password123",
    "name": "New User",
    "is_admin": false,
    "actor_id": "550e8400-e29b-41d4-a716-446655440001"
}
```

---

### Profile

#### Get Profile

```http
GET /api/profile
Authorization: Bearer <token>
```

#### Update Profile

```http
PUT /api/profile
Authorization: Bearer <token>
Content-Type: application/json

{
    "name": "Updated Name"
}
```

---

## JWT Token Validation Errors

Semua endpoint yang membutuhkan JWT akan mengembalikan error codes yang spesifik:

| Error Code | HTTP Status | Description |
|------------|-------------|-------------|
| `MISSING_AUTH_HEADER` | 401 | Authorization header tidak ada |
| `INVALID_AUTH_FORMAT` | 401 | Format authorization salah |
| `EMPTY_TOKEN` | 401 | Token kosong |
| `TOKEN_EXPIRED` | 401 | Token sudah expired |
| `TOKEN_NOT_YET_VALID` | 401 | Token belum aktif |
| `INVALID_SIGNATURE` | 401 | Signature token tidak valid |
| `MALFORMED_TOKEN` | 401 | Format token tidak valid |
| `INVALID_SIGNING_METHOD` | 401 | Method signing tidak valid |
| `INVALID_CLAIMS` | 401 | Claims token tidak valid |
| `INVALID_TOKEN` | 401 | Token tidak valid |

**Example Error Response:**
```json
{
    "success": false,
    "data": null,
    "error": "Token has expired. Please refresh your token or login again.",
    "code": "TOKEN_EXPIRED"
}
```

---

## Architecture

### Clean Architecture

Aplikasi ini mengikuti prinsip Clean Architecture dengan pemisahan yang jelas:

```
┌─────────────────────────────────────────────────────────────┐
│                    HANDLERS (HTTP Layer)                    │
│              HTTP Request/Response Handling                  │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                     USECASES (Business Logic)                │
│              Contains business rules and logic               │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                      PORTS (Interfaces)                      │
│              Defines contracts between layers                │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                   REPOSITORIES (Data Access)                 │
│              Database operations implementation              │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                       DOMAIN (Entities)                      │
│              Business entities and models                    │
└─────────────────────────────────────────────────────────────┘
```

### Concurrency Control

Sistem ini mengimplementasikan **double-layer concurrency control** untuk mencegah double approval dan race conditions.

#### Layer 1: In-Memory Mutex (sync.Map + sync.Mutex)

```go
// Per-request mutex untuk mencegah concurrent approval
// dalam instance aplikasi yang sama
func (s *RequestServiceImpl) LockRequest(requestID string) func() {
    mutex, _ := s.mutexMap.LoadOrStore(requestID, &sync.Mutex{})
    mu := mutex.(*sync.Mutex)
    mu.Lock()
    return func() { mu.Unlock() }
}
```

**Keunggulan:**
- Sangat cepat (operasi in-memory dalam nanoseconds)
- Tidak memerlukan dependencies eksternal
- Mencegah race conditions dalam satu instance aplikasi
- Granular (lock per-request, bukan global)

#### Layer 2: Database SELECT FOR UPDATE + Optimistic Locking

Menggunakan version field dan SELECT FOR UPDATE untuk perlindungan antar-instance.

```go
// Version field untuk optimistic locking
type Request struct {
    Version int `json:"version" gorm:"not null;default:1"`
}

// Increment version saat update
request.Version++
if err := s.requestRepo.Update(ctx, request); err != nil {
    return nil, err // Will fail if version mismatch
}
```

#### Mengapa Double-Layer?

| Layer | Kegunaan | Keterbatasan |
|-------|----------|--------------|
| **Mutex** | Cepat, sederhana, mencegah race conditions dalam instance yang sama | Tidak bekerja antar instance |
| **SELECT FOR UPDATE** | Jaminan database level antar instance | Lebih lambat dari mutex |
| **Version Field** | Mendeteksi concurrent modifications | Hanya mendeteksi, tidak mencegah |

Kombinasi ini memastikan double approval tidak mungkin terjadi bahkan dalam deployment terdistribusi.

### Authentication Flow

```
┌──────────┐     POST /auth/login      ┌─────────────┐
│  Client  │ ───────────────────────▶  │   Server    │
│          │                           │             │
│          │   {email, password}       │  Validate   │
│          │                           │  Credentials│
│          │                           │             │
│          │  {token, user_info} ◀──── │  Generate   │
│          │                           │  JWT Token  │
└──────────┘                           └─────────────┘
                                              │
                                              ▼
┌──────────┐     Authorization: Bearer ...  ┌──────────────────┐
│  Client  │ ───────────────────────▶       │   Server         │
│          │                                │                  │
│          │                                │Validate JWT      │
│          │                                │Check Permissions │
│          │                                │                  │
│          │                                │Process Request   │
│          │                                │                  │
│          │         Response         ◀──── │                  │
└──────────┘                                └──────────────────┘
```

---

## Workflow Process

### Alur Approval Biasa

```
┌─────────────────────────────────────────────────────────────────┐
│                    Workflow: Purchase Approval                   │
├─────────────────────────────────────────────────────────────────┤
│  Step 1: Manager (min_amount: 1,000,000)                        │
│  Step 2: Director (min_amount: 5,000,000)                       │
│  Step 3: Finance (all amounts)                                  │
└─────────────────────────────────────────────────────────────────┘

Request Created (Amount: 2,000,000)
        │
        ▼
┌───────────────────┐
│  Step 1: PENDING  │  Manager must approve
└───────────────────┘
        │
        │ Manager approves
        ▼
┌───────────────────────┐
│  Step 2: PENDING      │  Director must approve
└───────────────────────┘
        │
        │ Director approves
        ▼
┌─────────────────────┐
│  APPROVED            │  Final state
└─────────────────────┘
```

### Actor-based Authorization

Setiap workflow step memiliki `actor_id` yang menentukan peran yang dapat melakukan approval:

```
Workflow Step:
{
    "level": 1,
    "actor_id": "550e8400-e29b-41d4-a716-446655440001",  // Manager
    "conditions": {
        "min_amount": 1000000
    }
}

User's JWT Token:
{
    "user_id": "...",
    "actor_id": "550e8400-e29b-41d4-a716-446655440001",  // Must match!
    "is_admin": false
}

→ User with actor_id "550e8400-e29b-41d4-a716-446655440001" can approve
→ Admin can approve any step regardless of actor_id
```

---

## Testing

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./package/request/... -v
```

---

## Docker

### Build dan Run dengan Docker Compose (Recommended)

Cara termudah untuk menjalankan aplikasi adalah menggunakan Docker Compose yang akan menjalankan aplikasi dan database secara bersamaan.

```bash
# Build dan jalankan semua services
docker-compose up -d --build

# Lihat status container
docker-compose ps

# Lihat logs
docker-compose logs -f app

# Stop semua services
docker-compose down
```

**Akses Aplikasi:**
- URL: http://localhost:8082 (atau port yang dikonfigurasi di docker-compose.yml)
- Health check: http://localhost:8082/health

### Build dan Run dengan Docker (Manual)

Jika hanya ingin build image tanpa database:

```bash
# Build Docker image
docker build -t workflow-approval .

# Run container dengan koneksi ke database eksternal
docker run -p 8080:8080 \
  -e DB_HOST=host.docker.internal \
  -e DB_PORT=3306 \
  -e DB_USERNAME=root \
  -e DB_PASSWORD=your_password \
  -e DB_NAME=workflow_approval \
  -e JWT_SECRET=your-super-secret-key \
  workflow-approval
```

**Environment Variables yang tersedia:**
| Variable | Default | Description |
|----------|---------|-------------|
| `APP_HOST` | 0.0.0.0 | Application host |
| `APP_PORT` | 8080 | Application port |
| `APP_ENV` | production | Environment mode |
| `DB_HOST` | localhost | Database host |
| `DB_PORT` | 3306 | Database port |
| `DB_USERNAME` | root | Database username |
| `DB_PASSWORD` | - | Database password |
| `DB_NAME` | workflow_approval | Database name |
| `DB_CHARSET` | utf8mb4 | Database charset |
| `JWT_SECRET` | - | JWT secret key |
| `JWT_EXPIRATION` | 24 | Token expiration (hours) |
