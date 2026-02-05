# Transaction Backend API (Golang)

Project ini dibuat sebagai bagian dari **Technical Test Backend Developer (Golang)**.

---

## Tech Stack

**Backend**

- Go
- Gin (HTTP framework)
- GORM (ORM)
- MySQL (database utama)

**Testing**

- Go testing
- Testify
- SQLite

**Logging**

- zap

**Tools**

- Postman

---

## Setup Project

### 1. Buat Database

Pastikan MySQL sudah berjalan, lalu buat database berikut:

```sql
CREATE DATABASE transactions_db;
```

---

### 2. Clone Repository

```bash
git clone https://github.com/azdharsyahputra/transaction-technical-test.git
cd transaction-technical-test
```

---

### 3. Install Dependency

```bash
go mod tidy
```

---

### 4. Set Environment Variable Pada database.go

Contoh (Linux / macOS):

```bash
dsn := fmt.Sprintf(
    "%s:%s@tcp(%s:%s)/%s?parseTime=true",
    getEnv("DB_USER", "root"),
    getEnv("DB_PASSWORD", ""),
    getEnv("DB_HOST", "localhost"),
    getEnv("DB_PORT", "3306"),
    getEnv("DB_NAME", "transactions_db"),
)
```

Sesuaikan dengan konfigurasi MySQL di device masing-masing.

---

### 5. Jalankan Aplikasi

```bash
go run cmd/api/main.go
```

Server akan berjalan di:

```
http://localhost:8080
```

Database table akan otomatis dibuat menggunakan GORM AutoMigrate saat aplikasi dijalankan.

---

## Testing

Untuk menjalankan test:

```bash
go test ./internal/... -coverprofile=coverage.out
```

```bash
go tool cover -func=coverage.out
```

---

## API Documentation (Postman)

Seluruh endpoint API telah diuji menggunakan Postman.  
Postman collection sudah diexport dan tersedia di repository ini.

File :

```
technical test.postman_collection.json
```

---
