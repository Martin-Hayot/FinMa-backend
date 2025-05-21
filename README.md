# About

**FinMa** is a personal finance manager built with Go. Track expenses, manage budgets, and gain insights into your financial health with a simple, secure, and efficient backend.

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


## MakeFile

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
