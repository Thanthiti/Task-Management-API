# Task Management API

A simple RESTful API for task management built with Go, GORM, and PostgreSQL (or SQLite for testing). Supports creating, reading, updating, and deleting tasks.

---

## Features

- User Registration & Login (with JWT Authentication)
- Task CRUD (Create, Read, Update, Delete)
- Middleware (Authentication, Logging, Error handling)
- PostgreSQL with GORM
- Environment-based config loading
- Dockerized for easy development & deployment
- Unit testing with testify + SQLite (for speed)

---

## Tech Stack

- **Go** (v1.21+)
- **Fiber** (Web framework)
- **GORM** (ORM for Go)
- **PostgreSQL** (Relational Database)
- **Docker** / **Docker Compose**
- **JWT** for secure user auth
- **Logrus** for structured logging
- **Validator** for input validation
- **SQLite** for lightweight test DB
- **Testify** for unit testing

---

## Project Structure (Clean Architecture)

your-project/
│
├── cmd/
│ └── server/ # Main entrypoint (main.go)
│
├── internal/
│ ├── user/
│ │ ├── handler/ # HTTP handlers
│ │ ├── model/ # Data structures
│ │ ├── repository/ # DB operations
│ │ └── usecase/ # Business logic
│ │
│ └── task/
│ ├── handler/
│ ├── model/
│ ├── repository/
│ └── usecase/
│
├── pkg/
│ ├── config/ # DB and environment setup
│ ├── logger/ # Logrus setup
│ ├── middleware/ # Fiber middlewares
│ ├── helper/ # Utilities
│ └── auth/ # JWT helpers
│
├── .env.example # Sample env file
├── Dockerfile # App Dockerfile
├── docker-compose.yml # Compose file for app + DB
├── go.mod
└── README.md 

---

## Setup Instructions

### 1.Clone the repo

```bash
git clone https://github.com/yourusername/task-api.git
cd task-api
```

### 2.Setup Environment Variables 

```bash
cp .env.example .env
```
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=yourpassword
DB_NAME=taskdb
DB_SSL=disable
JWT_SECRET=your_jwt_secret

### 3. Start the App with Docker Compose
```bash
docker-compose up --build
```
API Server 
http://localhost:8080

---
### Running Unit Tests

To run all unit tests in the project, use the following command:

```bash
GO_ENV=test go test ./...
```
---
