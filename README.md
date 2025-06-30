# About

**FinMa** is a personal finance manager built with Go. Track expenses, manage budgets, and gain insights into your financial health with a simple, secure, and efficient backend.
There is also a frontend for this project: [FinMa Frontend](https://github.com/Martin-Hayot/FinMa-frontend.git) 

---

## 🚀 Getting Started

Follow these steps to run **FinMa** locally on your machine.

### ✅ Prerequisites

Make sure you have the following installed:

* [Go](https://go.dev/doc/install) (v1.20+ recommended)
* [Docker](https://www.docker.com/products/docker-desktop/)
* [Docker Compose](https://docs.docker.com/compose/install/)

### 🛠️ Installation

1. **Clone the repository**

```bash
git clone https://github.com/Martin-Hayot/FinMa-backend.git
cd FinMa-backend
```

2. **Copy the content of the .env.example into .env**

2. **Start the database using Docker Compose**

```bash
docker-compose up -d
```

This will start the database and any required services.

3. **Run the Go application**

```bash
go run main.go
```

4. **Access the API**

The server will run at `http://localhost:PORT` (port configured in `.env` file).

---
## **Project Structure**
The FinMa backend is a Go-based financial management API built with clean architecture principles. Here's an overview of the key directories and their purposes:

```
FinMa-backend/
├── config/             # Application configuration 
├── constants/          # Application-wide constants
├── dto/                # Data Transfer Objects for API requests/responses
├── internal/           # Core application code (not exposed as packages)
│   ├── api/            # API layer (HTTP server, routes, handlers)
│   │   ├── handlers/   # HTTP request handlers
│   │   ├── middleware/ # HTTP middleware
│   ├── domain/         # Domain entities and business logic
│   ├── repository/     # Data access layer
│   │   ├── postgres/   # PostgreSQL implementation
│   ├── service/        # Business logic layer
├── pkg/                # External APIs
│   ├── gocardless/     # GoCardless Client and API interactions
├── utils/              # Utility functions
├── main.go             # Application entry point
```

### Key Features
- Clean Architecture: Separation of concerns with layers for API, services, and repositories
- Authentication: JWT-based authentication with refresh tokens
- Database: PostgreSQL with GORM for ORM functionality
- API: RESTful API built with Fiber framework
- Validation: Request validation using go-playground/validator
- GoCardless Integration: Bank account data API
- Email Services: Transactional emails with Resend
### Important Files
- main.go: Application entry point
- config.go: Configuration management
- server.go: HTTP server setup
- routes.go: API route definitions
- models.go: Domain entities
- auth.go: Authentication business logic

### Technologies
- Go 1.23
- Fiber web framework
- GORM ORM
- PostgreSQL
- JWT authentication
- GoCardless Bank Account API
- Docker & Docker Compose

---

## MakeFile Commands

run all make commands 
```bash
make all build
```

build the application
```bash
make build
```

run the application
```bash
make run
```

Create DB container
```bash
make docker-run
```

Shutdown DB container
```bash
make docker-down
```

live reload the application
```bash
make watch
```

run the test suite
```bash
make test
```

clean up binary from the last build
```bash
make clean
```
