# 📚 Belajar Golang RESTful API

> Catatan belajar membangun RESTful API menggunakan Go (Golang) dengan arsitektur **Repository-Service-Controller Pattern**, database MySQL, dan HTTP Router.

---

## 📋 Daftar Isi

- [Overview](#-overview)
- [Tech Stack & Dependencies](#-tech-stack--dependencies)
- [Arsitektur & Struktur Project](#-arsitektur--struktur-project)
- [Diagram Alur Request](#-diagram-alur-request)
- [Setup Database](#-setup-database)
- [Penjelasan Setiap Layer](#-penjelasan-setiap-layer)
  - [1. Model Layer](#1-model-layer)
  - [2. Repository Layer](#2-repository-layer)
  - [3. Service Layer](#3-service-layer)
  - [4. Controller Layer](#4-controller-layer)
  - [5. Helper Functions](#5-helper-functions)
  - [6. Middleware](#6-middleware)
  - [7. Exception Handling](#7-exception-handling)
  - [8. App (Database & Router)](#8-app-database--router)
  - [9. Main Entry Point](#9-main-entry-point)
- [API Specification](#-api-specification)
- [Testing](#-testing)
- [Cara Menjalankan](#-cara-menjalankan)
- [Pelajaran Penting (Key Takeaways)](#-pelajaran-penting-key-takeaways)

---

## 🔍 Overview

Project ini adalah implementasi **CRUD RESTful API** untuk entity `Category` menggunakan **pure Go** (tanpa framework berat seperti Gin/Echo). Tujuannya adalah memahami bagaimana membangun API dari nol dengan:

- **HTTP Router** (`httprouter`) sebagai multiplexer
- **Database MySQL** dengan `database/sql` (tanpa ORM)
- **Manual Dependency Injection** (tanpa wire/dig)
- **Layered Architecture**: Repository → Service → Controller
- **Middleware** untuk autentikasi API Key
- **Centralized Error Handling** via `PanicHandler`
- **Request Validation** menggunakan `go-playground/validator`
- **Integration Testing** menggunakan `httptest`

---

## 🛠 Tech Stack & Dependencies

| Package | Kegunaan |
|---|---|
| `net/http` | Standard library HTTP server |
| `database/sql` | Standard library database driver |
| [`github.com/julienschmidt/httprouter`](https://github.com/julienschmidt/httprouter) | Lightweight HTTP request router |
| [`github.com/go-sql-driver/mysql`](https://github.com/go-sql-driver/mysql) | MySQL driver untuk `database/sql` |
| [`github.com/go-playground/validator/v10`](https://github.com/go-playground/validator) | Struct validation (required, min, max) |
| [`github.com/stretchr/testify`](https://github.com/stretchr/testify) | Assertion library untuk testing |

**Go Version:** `1.24.0`

---

## 🏗 Arsitektur & Struktur Project

```
belajar-golang-restful-api/
├── main.go                          # Entry point aplikasi
├── go.mod                           # Module definition & dependencies
├── go.sum                           # Dependency checksums
├── apispec.yml                      # OpenAPI 3.0 specification
├── test.http                        # HTTP client test file (REST Client)
│
├── app/                             # Application setup
│   ├── database.go                  # Konfigurasi koneksi MySQL
│   └── router.go                    # Route definitions (URL mapping)
│
├── model/                           # Data models
│   ├── domain/
│   │   └── category.go              # Domain entity (representasi tabel DB)
│   └── web/
│       ├── web_response.go          # Generic API response wrapper
│       ├── category_create_request.go  # DTO untuk create request
│       ├── category_update_request.go  # DTO untuk update request
│       └── category_response.go     # DTO untuk response ke client
│
├── repository/                      # Data Access Layer
│   ├── category_repository.go       # Interface repository
│   └── category_repository_impl.go  # Implementasi repository (SQL queries)
│
├── service/                         # Business Logic Layer
│   ├── category_service.go          # Interface service
│   └── category_service_impl.go     # Implementasi service (validasi + transaksi)
│
├── controller/                      # Presentation Layer
│   ├── category_controller.go       # Interface controller
│   └── category_controller_impl.go  # Implementasi controller (HTTP handler)
│
├── helper/                          # Utility functions
│   ├── error.go                     # Panic error helper
│   ├── json.go                      # JSON encode/decode helper
│   ├── model.go                     # Domain-to-Response converter
│   └── tx.go                        # Transaction commit/rollback helper
│
├── middleware/                       # HTTP Middleware
│   └── auth_middleware.go           # API Key authentication
│
├── exception/                       # Error handling
│   ├── error_handler.go             # Centralized panic/error handler
│   └── not_found_error.go           # Custom NotFoundError type
│
└── test/                            # Integration tests
    └── category_controller_test.go  # End-to-end API tests
```

---

## 🔄 Diagram Alur Request

Berikut adalah alur lengkap sebuah HTTP request melewati setiap layer:

```
┌──────────────────────────────────────────────────────────────────────┐
│                        HTTP Client (Postman/cURL)                    │
└─────────────────────────────────┬────────────────────────────────────┘
                                  │
                                  ▼
┌──────────────────────────────────────────────────────────────────────┐
│  1. http.Server (localhost:3000)                                      │
│     - Menerima incoming HTTP request                                 │
└─────────────────────────────────┬────────────────────────────────────┘
                                  │
                                  ▼
┌──────────────────────────────────────────────────────────────────────┐
│  2. AuthMiddleware                                                   │
│     - Cek header "X-API-Key"                                         │
│     - Jika VALID → lanjut ke router                                  │
│     - Jika INVALID → return 401 Unauthorized                         │
└─────────────────────────────────┬────────────────────────────────────┘
                                  │
                                  ▼
┌──────────────────────────────────────────────────────────────────────┐
│  3. httprouter.Router                                                │
│     - Match URL pattern ke handler function                          │
│     - PanicHandler → exception.ErrorHandler (catch semua panic)      │
└─────────────────────────────────┬────────────────────────────────────┘
                                  │
                                  ▼
┌──────────────────────────────────────────────────────────────────────┐
│  4. Controller                                                       │
│     - Parse HTTP request (body, params)                              │
│     - Panggil Service layer                                          │
│     - Format response JSON                                           │
└─────────────────────────────────┬────────────────────────────────────┘
                                  │
                                  ▼
┌──────────────────────────────────────────────────────────────────────┐
│  5. Service                                                          │
│     - Validasi request (validator)                                   │
│     - Begin transaction (db.Begin)                                   │
│     - Panggil Repository layer                                       │
│     - defer CommitOrRollback                                         │
└─────────────────────────────────┬────────────────────────────────────┘
                                  │
                                  ▼
┌──────────────────────────────────────────────────────────────────────┐
│  6. Repository                                                       │
│     - Execute SQL query (INSERT/UPDATE/DELETE/SELECT)                 │
│     - Return domain entity                                           │
└─────────────────────────────────┬────────────────────────────────────┘
                                  │
                                  ▼
┌──────────────────────────────────────────────────────────────────────┐
│  7. MySQL Database                                                   │
│     - Tabel: category (id INT AUTO_INCREMENT, name VARCHAR)          │
└──────────────────────────────────────────────────────────────────────┘
```

---

## 💾 Setup Database

### 1. Buat Database untuk Development

```sql
CREATE DATABASE belajar_golang_restful_api;
```

### 2. Buat Database untuk Testing

```sql
CREATE DATABASE belajar_golang_restful_api_test;
```

### 3. Buat Tabel `category`

Jalankan di **kedua database** (development & test):

```sql
USE belajar_golang_restful_api; -- atau belajar_golang_restful_api_test

CREATE TABLE category (
    id   INT NOT NULL AUTO_INCREMENT,
    name VARCHAR(200) NOT NULL,
    PRIMARY KEY (id)
) ENGINE = InnoDB;
```

### 4. Verifikasi Koneksi

Pastikan MySQL berjalan di `localhost:3306` dengan user `root` tanpa password (sesuai konfigurasi di `app/database.go`).

```
DSN: root@tcp(localhost:3306)/belajar_golang_restful_api
```

---

## 📖 Penjelasan Setiap Layer

### 1. Model Layer

Model layer berisi **struct** yang merepresentasikan data. Dibagi jadi dua sub-package:

#### a. `model/domain` — Domain Entity

```go
// model/domain/category.go
type Category struct {
    Id   int
    Name string
}
```

**Tujuan:** Merepresentasikan data di database (tabel `category`). Struct ini **tidak punya JSON tag** karena hanya digunakan secara internal.

#### b. `model/web` — DTO (Data Transfer Object)

DTO berfungsi sebagai kontrak data antara client dan server:

| File | Struct | Fungsi |
|---|---|---|
| `category_create_request.go` | `CategoryCreateRequest` | Menerima data create (hanya `name`) |
| `category_update_request.go` | `CategoryUpdateRequest` | Menerima data update (`id` + `name`) |
| `category_response.go` | `CategoryResponse` | Format data yang dikirim ke client |
| `web_response.go` | `WebResponse` | Wrapper generic untuk semua response |

**Contoh validation tag:**

```go
type CategoryCreateRequest struct {
    Name string `validate:"required,max=200,min=1" json:"name"`
}
```

- `validate:"required"` → Field wajib diisi
- `validate:"max=200"` → Maksimal 200 karakter
- `validate:"min=1"` → Minimal 1 karakter
- `json:"name"` → Key JSON saat serialize/deserialize

**Kenapa pakai DTO terpisah?**
- **Separation of Concerns**: Domain entity tidak terikat dengan format request/response HTTP
- **Security**: Client tidak bisa inject field yang tidak diinginkan (misalnya `Id` saat create)
- **Validation**: Bisa pasang validation rules berbeda per operasi

---

### 2. Repository Layer

Repository layer bertanggung jawab untuk **akses data ke database**. Ini adalah layer paling bawah yang langsung berinteraksi dengan SQL.

#### a. Interface (`category_repository.go`)

```go
type CategoryRepository interface {
    Save(ctx context.Context, tx *sql.Tx, category domain.Category) domain.Category
    Update(ctx context.Context, tx *sql.Tx, category domain.Category) domain.Category
    Delete(ctx context.Context, tx *sql.Tx, category domain.Category)
    FindById(ctx context.Context, tx *sql.Tx, categoryId int) (domain.Category, error)
    FindAll(ctx context.Context, tx *sql.Tx) []domain.Category
}
```

**Kenapa pakai interface?**
- Memungkinkan **dependency injection** — implementasi bisa diganti (mock) saat testing
- Mengikuti prinsip **Dependency Inversion** (SOLID)

**Kenapa parameter `*sql.Tx` (transaction), bukan `*sql.DB`?**
- Agar repository **tidak mengelola transaksi sendiri**
- Transaksi dikelola oleh **Service layer** → memastikan satu operasi bisnis = satu transaksi

#### b. Implementasi (`category_repository_impl.go`)

Setiap method menjalankan raw SQL query:

```go
// INSERT
func (c CategoryRepositoryImpl) Save(ctx context.Context, tx *sql.Tx, Category domain.Category) domain.Category {
    SQL := "INSERT INTO category(name) VALUES (?);"
    result, err := tx.ExecContext(ctx, SQL, Category.Name)
    helper.PanicIfError(err)

    id, err := result.LastInsertId()
    helper.PanicIfError(err)

    Category.Id = int(id)
    return Category
}

// SELECT by ID
func (c CategoryRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, CategoryId int) (domain.Category, error) {
    SQL := "SELECT id, name FROM category WHERE id = ?;"
    rows, err := tx.QueryContext(ctx, SQL, CategoryId)
    helper.PanicIfError(err)
    defer rows.Close()

    Category := domain.Category{}
    if rows.Next() {
        err = rows.Scan(&Category.Id, &Category.Name)
        helper.PanicIfError(err)
        return Category, nil
    } else {
        return Category, errors.New("category is not found")
    }
}
```

**Pola penting:**
- Menggunakan **parameterized query** (`?`) → mencegah **SQL Injection**
- `defer rows.Close()` → memastikan rows ditutup setelah selesai (mencegah connection leak)
- `FindById` return `error` jika data tidak ditemukan → dihandle sebagai `NotFoundError` di Service

---

### 3. Service Layer

Service layer berisi **business logic**. Di sinilah validasi dan manajemen transaksi terjadi.

#### a. Interface (`category_service.go`)

```go
type CategoryService interface {
    Create(ctx context.Context, request web.CategoryCreateRequest) web.CategoryResponse
    Update(ctx context.Context, request web.CategoryUpdateRequest) web.CategoryResponse
    Delete(ctx context.Context, categoryId int)
    FindById(ctx context.Context, categoryId int) web.CategoryResponse
    FindAll(ctx context.Context) []web.CategoryResponse
}
```

**Perbedaan dengan Repository interface:**
- Service menerima **DTO (web model)**, bukan domain model
- Service tidak menerima `*sql.Tx` — karena transaksi dikelola di sini

#### b. Implementasi (`category_service_impl.go`)

```go
type CategoryServiceImpl struct {
    CategoryRepository repository.CategoryRepository
    DB                 *sql.DB
    Validate           *validator.Validate
}
```

**Dependencies yang di-inject:**
1. `CategoryRepository` — akses data
2. `DB` (*sql.DB) — untuk membuka transaksi baru
3. `Validate` — validator instance

**Alur tiap method:**

```go
func (service CategoryServiceImpl) Create(ctx context.Context, request web.CategoryCreateRequest) web.CategoryResponse {
    // 1. Validasi request
    err := service.Validate.Struct(request)
    helper.PanicIfError(err)  // panic → ditangkap ErrorHandler → return 400

    // 2. Begin transaction
    tx, err := service.DB.Begin()
    helper.PanicIfError(err)
    defer helper.CommitOrRollback(tx)  // auto commit/rollback

    // 3. Convert DTO → Domain
    category := domain.Category{
        Name: request.Name,
    }

    // 4. Simpan via repository
    category = service.CategoryRepository.Save(ctx, tx, category)

    // 5. Convert Domain → Response DTO
    return helper.ToCategoryResponse(category)
}
```

**Pattern penting: `defer helper.CommitOrRollback(tx)`**

```go
func CommitOrRollback(tx *sql.Tx) {
    err := recover()        // tangkap panic (jika ada)
    if err != nil {
        tx.Rollback()       // ada error → rollback
        panic(err)          // re-panic agar ditangkap ErrorHandler
    } else {
        tx.Commit()         // tidak ada error → commit
    }
}
```

Ini adalah **mekanisme auto-transaction**:
- Jika semua berjalan lancar → **Commit**
- Jika ada panic (validation error, SQL error, not found) → **Rollback** lalu re-panic
- Re-panic akan ditangkap oleh `httprouter.PanicHandler` → `exception.ErrorHandler`

---

### 4. Controller Layer

Controller layer bertanggung jawab untuk **menerima HTTP request** dan **mengirim HTTP response**. Tidak ada business logic di sini.

#### a. Interface (`category_controller.go`)

```go
type CategoryController interface {
    Create(w http.ResponseWriter, r *http.Request, params httprouter.Params)
    Update(w http.ResponseWriter, r *http.Request, params httprouter.Params)
    Delete(w http.ResponseWriter, r *http.Request, params httprouter.Params)
    FindById(w http.ResponseWriter, r *http.Request, params httprouter.Params)
    FindAll(w http.ResponseWriter, r *http.Request, params httprouter.Params)
}
```

Method signature mengikuti format `httprouter.Handle`.

#### b. Implementasi (`category_controller_impl.go`)

**Contoh Create handler:**

```go
func (c CategoryControllerImpl) Create(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
    // 1. Decode JSON body → DTO
    categoryCreateRequest := web.CategoryCreateRequest{}
    helper.ReadFromRequestBody(r, &categoryCreateRequest)

    // 2. Panggil service
    categoryResponse := c.CategoryService.Create(r.Context(), categoryCreateRequest)

    // 3. Bungkus dalam WebResponse
    webResponse := web.WebResponse{
        Code:   200,
        Status: "OK",
        Data:   categoryResponse,
    }

    // 4. Encode ke JSON dan kirim
    helper.WriteToResponseBody(w, webResponse)
}
```

**Contoh Update handler (dengan path parameter):**

```go
func (c CategoryControllerImpl) Update(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
    categoryUpdateRequest := web.CategoryUpdateRequest{}
    helper.ReadFromRequestBody(r, &categoryUpdateRequest)

    // Ambil :categoryId dari URL path
    categoryId := params.ByName("categoryId")
    id, err := strconv.Atoi(categoryId)  // Convert string → int
    helper.PanicIfError(err)

    categoryUpdateRequest.Id = id  // Set ID dari URL

    categoryResponse := c.CategoryService.Update(r.Context(), categoryUpdateRequest)
    // ... format response
}
```

**Catatan:** `Id` pada `CategoryUpdateRequest` datang dari **URL path**, bukan dari request body. Ini common pattern di REST API.

---

### 5. Helper Functions

Helper berisi utility functions yang digunakan di seluruh project:

#### a. `helper/error.go` — Panic Helper

```go
func PanicIfError(err error) {
    if err != nil {
        panic(err)
    }
}
```

Digunakan sebagai **shortcut** error handling. Semua panic akan ditangkap oleh `ErrorHandler`.

> **Catatan:** Pola ini valid untuk project learning. Di production, sebaiknya gunakan explicit error return.

#### b. `helper/json.go` — JSON Encode/Decode

```go
func ReadFromRequestBody(request *http.Request, result interface{}) {
    decoder := json.NewDecoder(request.Body)
    err := decoder.Decode(result)
    PanicIfError(err)
}

func WriteToResponseBody(writer http.ResponseWriter, response interface{}) {
    writer.Header().Add("Content-Type", "application/json")
    encoder := json.NewEncoder(writer)
    err := encoder.Encode(response)
    PanicIfError(err)
}
```

- `json.NewDecoder` → streaming decoder (lebih efisien daripada `json.Unmarshal` untuk HTTP body)
- `json.NewEncoder` → streaming encoder langsung ke `ResponseWriter`

#### c. `helper/model.go` — Model Converter

```go
func ToCategoryResponse(category domain.Category) web.CategoryResponse {
    return web.CategoryResponse{
        Id:   category.Id,
        Name: category.Name,
    }
}

func ToCategoryResponses(categories []domain.Category) []web.CategoryResponse {
    categoryResponses := []web.CategoryResponse{}
    for _, category := range categories {
        categoryResponses = append(categoryResponses, ToCategoryResponse(category))
    }
    return categoryResponses
}
```

Ini adalah **mapper** dari domain entity ke response DTO. Memastikan hanya data yang boleh dikirim ke client yang ter-expose.

#### d. `helper/tx.go` — Transaction Helper

```go
func CommitOrRollback(tx *sql.Tx) {
    err := recover()
    if err != nil {
        errorRollback := tx.Rollback()
        PanicIfError(errorRollback)
        panic(err)
    } else {
        errorCommit := tx.Commit()
        PanicIfError(errorCommit)
    }
}
```

Digunakan dengan `defer` di Service layer. Mekanisme:
1. `recover()` menangkap panic yang terjadi di dalam function
2. Jika ada panic → **Rollback** transaksi, lalu **re-panic**
3. Jika tidak ada panic → **Commit** transaksi

---

### 6. Middleware

#### `middleware/auth_middleware.go` — API Key Authentication

```go
type AuthMiddleware struct {
    Handler http.Handler
}

func (middleware *AuthMiddleware) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
    if "RAHASIA" == request.Header.Get("X-API-Key") {
        middleware.Handler.ServeHTTP(writer, request)  // ✅ Lanjut
    } else {
        writer.WriteHeader(http.StatusUnauthorized)
        webResponse := web.WebResponse{
            Code:   http.StatusUnauthorized,
            Status: "UNAUTHORIZED",
        }
        helper.WriteToResponseBody(writer, webResponse)  // ❌ Tolak
    }
}
```

**Konsep:**
- `AuthMiddleware` implement `http.Handler` interface (method `ServeHTTP`)
- Membungkus router sebagai inner handler
- Request dengan header `X-API-Key: RAHASIA` akan diteruskan
- Request tanpa/salah API key akan di-reject dengan status `401`

**Cara pasang middleware (di `main.go`):**

```go
server := http.Server{
    Addr:    "localhost:3000",
    Handler: middleware.NewAuthMiddleware(router),  // router dibungkus middleware
}
```

---

### 7. Exception Handling

#### a. `exception/not_found_error.go` — Custom Error Type

```go
type NotFoundError struct {
    Error string
}

func NewNotFoundError(error string) NotFoundError {
    return NotFoundError{Error: error}
}
```

Custom error type ini digunakan untuk membedakan "not found" error dari error lainnya melalui **type assertion**.

#### b. `exception/error_handler.go` — Centralized Error Handler

```go
func ErrorHandler(writer http.ResponseWriter, request *http.Request, err interface{}) {
    if notFoundError(writer, request, err) {
        return  // 404
    }
    if validationErrors(writer, request, err) {
        return  // 400
    }
    internalServerError(writer, request, err)  // 500
}
```

**Mekanisme:**
1. Di-register sebagai `router.PanicHandler` di `app/router.go`
2. Setiap kali ada `panic()` dalam handler → function ini dipanggil
3. Menggunakan **type assertion** untuk menentukan tipe error:
   - `NotFoundError` → HTTP 404 Not Found
   - `validator.ValidationErrors` → HTTP 400 Bad Request
   - Lainnya → HTTP 500 Internal Server Error

**Flow error handling:**

```
panic(exception.NewNotFoundError("..."))
    ↓
httprouter.PanicHandler catches panic
    ↓
exception.ErrorHandler(writer, request, err)
    ↓
Type assertion: err.(NotFoundError) → true
    ↓
Return JSON: {"code": 404, "status": "NOT FOUND", "data": "category is not found"}
```

---

### 8. App (Database & Router)

#### a. `app/database.go` — Konfigurasi Database

```go
func NewDB() *sql.DB {
    db, err := sql.Open("mysql", "root@tcp(localhost:3306)/belajar_golang_restful_api")
    helper.PanicIfError(err)

    // setup database pooling
    db.SetMaxIdleConns(5)                   // minimal jumlah koneksi standby (idle)
    db.SetMaxOpenConns(20)                  // maksimal jumlah koneksi yang bisa dibuka
    db.SetConnMaxLifetime(60 * time.Minute) // berapa lama koneksi boleh digunakan sebelum direfresh
    db.SetConnMaxIdleTime(10 * time.Minute) // berapa lama koneksi idle boleh bertahan sebelum dihapus

    return db
}
```

**Connection Pool Settings:**

| Setting | Value | Penjelasan |
|---|---|---|
| `MaxIdleConns` | 5 | Minimal jumlah koneksi standby yang siap dilayani sewaktu-waktu |
| `MaxOpenConns` | 20 | Batas maksimal koneksi paralel yang dapat dibuat dari aplikasi ke DB |
| `ConnMaxLifetime` | 60 min | Durasi mutlak koneksi dapat digunakan. Menghindari stale connection dari MySQL server (refresh per jam). |
| `ConnMaxIdleTime` | 10 min | Waktu tenggang sebelum koneksi idle/menganggur ditutup untuk menghemat resource |

#### b. `app/router.go` — Route Definitions

```go
func NewRouter(categoryController controller.CategoryController) *httprouter.Router {
    router := httprouter.New()

    router.GET("/api/categories", categoryController.FindAll)
    router.GET("/api/categories/:categoryId", categoryController.FindById)
    router.POST("/api/categories", categoryController.Create)
    router.PUT("/api/categories/:categoryId", categoryController.Update)
    router.DELETE("/api/categories/:categoryId", categoryController.Delete)

    router.PanicHandler = exception.ErrorHandler

    return router
}
```

**Route mapping:**

| Method | Path | Handler | Keterangan |
|---|---|---|---|
| `GET` | `/api/categories` | `FindAll` | List semua category |
| `GET` | `/api/categories/:categoryId` | `FindById` | Get category by ID |
| `POST` | `/api/categories` | `Create` | Buat category baru |
| `PUT` | `/api/categories/:categoryId` | `Update` | Update category |
| `DELETE` | `/api/categories/:categoryId` | `Delete` | Hapus category |

**`:categoryId`** adalah **named parameter** dari httprouter, diakses via `params.ByName("categoryId")`.

---

### 9. Main Entry Point

```go
func main() {
    // 1. Setup database connection
    db := app.NewDB()

    // 2. Setup validator
    validate := validator.New()

    // 3. Manual Dependency Injection (bottom-up)
    categoryRepository := repository.NewCategoryRepository()
    categoryService := service.NewCategoryService(categoryRepository, db, validate)
    categoryController := controller.NewCategoryController(categoryService)

    // 4. Setup router
    router := app.NewRouter(categoryController)

    // 5. Setup HTTP server dengan middleware
    server := http.Server{
        Addr:    "localhost:3000",
        Handler: middleware.NewAuthMiddleware(router),
    }

    // 6. Start server
    err := server.ListenAndServe()
    helper.PanicIfError(err)
}
```

**Dependency Injection Flow:**

```
Repository (tanpa dependency)
    ↓ injected ke
Service (+ DB + Validator)
    ↓ injected ke
Controller
    ↓ injected ke
Router
    ↓ dibungkus oleh
AuthMiddleware
    ↓ dipasang ke
http.Server
```

Ini adalah **manual/poor man's dependency injection** — semua wiring dilakukan secara eksplisit di `main()`. Untuk project yang lebih besar, bisa menggunakan library seperti [Wire](https://github.com/google/wire) atau [Dig](https://github.com/uber-go/dig).

---

## 📡 API Specification

API menggunakan format **OpenAPI 3.0** (lihat `apispec.yml`).

### Authentication

Semua endpoint memerlukan header:

```
X-API-Key: RAHASIA
```

### Endpoints

#### 1. Create Category

```http
POST /api/categories
Content-Type: application/json
X-API-Key: RAHASIA

{
  "name": "Gadget"
}
```

**Response (200):**
```json
{
  "code": 200,
  "status": "OK",
  "data": {
    "id": 1,
    "name": "Gadget"
  }
}
```

#### 2. List All Categories

```http
GET /api/categories
X-API-Key: RAHASIA
```

**Response (200):**
```json
{
  "code": 200,
  "status": "OK",
  "data": [
    { "id": 1, "name": "Gadget" },
    { "id": 2, "name": "Fashion" }
  ]
}
```

#### 3. Get Category by ID

```http
GET /api/categories/1
X-API-Key: RAHASIA
```

#### 4. Update Category

```http
PUT /api/categories/1
Content-Type: application/json
X-API-Key: RAHASIA

{
  "name": "Updated Name"
}
```

#### 5. Delete Category

```http
DELETE /api/categories/1
X-API-Key: RAHASIA
```

### Error Responses

| Kode | Status | Kapan Terjadi |
|---|---|---|
| `400` | BAD REQUEST | Validasi gagal (name kosong, terlalu panjang) |
| `401` | UNAUTHORIZED | API Key salah atau tidak ada |
| `404` | NOT FOUND | Category dengan ID tersebut tidak ditemukan |
| `500` | INTERNAL SERVER ERROR | Error tidak terduga (DB down, dll) |

---

## 🧪 Testing

Project ini menggunakan **Integration Test** — test langsung menghit full stack (controller → service → repository → database).

### Setup Test

Test menggunakan **database terpisah** (`belajar_golang_restful_api_test`) agar tidak mengganggu data development:

```go
func setupTestDB() *sql.DB {
    db, err := sql.Open("mysql", "root@tcp(localhost:3306)/belajar_golang_restful_api_test")
    // ... connection pool settings
    return db
}
```

Setiap test di-awali dengan `truncateCategory(db)` untuk memastikan **clean state**.

### Test menggunakan `httptest`

```go
func TestCreateCategorySuccess(t *testing.T) {
    db := setupTestDB()
    truncateCategory(db)
    router := setupRouter(db)

    // 1. Buat request
    requestBody := strings.NewReader(`{"name" : "Gadget"}`)
    request := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/categories", requestBody)
    request.Header.Add("Content-Type", "application/json")
    request.Header.Add("X-API-Key", "RAHASIA")

    // 2. Record response
    recorder := httptest.NewRecorder()
    router.ServeHTTP(recorder, request)

    // 3. Assert
    response := recorder.Result()
    assert.Equal(t, 200, response.StatusCode)

    body, _ := io.ReadAll(response.Body)
    var responseBody map[string]interface{}
    json.Unmarshal(body, &responseBody)

    assert.Equal(t, 200, int(responseBody["code"].(float64)))
    assert.Equal(t, "OK", responseBody["status"])
    assert.Equal(t, "Gadget", responseBody["data"].(map[string]interface{})["name"])
}
```

### Daftar Test Cases

| Test Function | Skenario | Expected |
|---|---|---|
| `TestCreateCategorySuccess` | Create dengan name valid | 200 OK |
| `TestCreateCategoryFailed` | Create dengan name kosong | 400 Bad Request |
| `TestUpdateCategorySuccess` | Update category yang ada | 200 OK |
| `TestUpdateCategoryFailed` | Update dengan name kosong | 400 Bad Request |
| `TestGetCategorySuccess` | Get category yang ada | 200 OK |
| `TestGetCategoryFailed` | Get category yang tidak ada | 404 Not Found |
| `TestDeleteCategorySuccess` | Delete category yang ada | 200 OK |
| `TestDeleteCategoryFailed` | Delete category yang tidak ada | 404 Not Found |
| `TestListCategoriesSuccess` | List multiple categories | 200 OK + array |
| `TestUnauthorized` | Request tanpa API Key valid | 401 Unauthorized |

### Menjalankan Test

```bash
# Jalankan semua test
go test -v ./test/

# Jalankan test spesifik
go test -v -run TestCreateCategorySuccess ./test/
```

---

## 🚀 Cara Menjalankan

### Prerequisites

- Go 1.24+
- MySQL running di `localhost:3306`
- Database dan tabel sudah dibuat (lihat [Setup Database](#-setup-database))

### Langkah

```bash
# 1. Clone / masuk ke directory project
cd belajar-golang-restful-api

# 2. Download dependencies
go mod tidy

# 3. Jalankan server
go run main.go

# Server berjalan di http://localhost:3000
```

### Test API dengan cURL

```bash
# List all categories
curl -X GET http://localhost:3000/api/categories \
  -H "X-API-Key: RAHASIA"

# Create category
curl -X POST http://localhost:3000/api/categories \
  -H "X-API-Key: RAHASIA" \
  -H "Content-Type: application/json" \
  -d '{"name": "Gadget"}'

# Get by ID
curl -X GET http://localhost:3000/api/categories/1 \
  -H "X-API-Key: RAHASIA"

# Update
curl -X PUT http://localhost:3000/api/categories/1 \
  -H "X-API-Key: RAHASIA" \
  -H "Content-Type: application/json" \
  -d '{"name": "Electronics"}'

# Delete
curl -X DELETE http://localhost:3000/api/categories/1 \
  -H "X-API-Key: RAHASIA"
```

---

## 💡 Pelajaran Penting (Key Takeaways)

### 1. Layered Architecture Pattern

```
Controller (HTTP) → Service (Business Logic) → Repository (Data Access) → Database
```

- Setiap layer punya **tanggung jawab tunggal** (Single Responsibility)
- Komunikasi antar layer melalui **interface** (Dependency Inversion)
- Perubahan di satu layer **tidak mempengaruhi** layer lain

### 2. Interface-Based Design di Go

Go tidak punya class/inheritance. Sebagai gantinya:
- Definisikan **interface** untuk kontrak
- Buat **struct** sebagai implementasi
- **Constructor function** (`New...()`) return interface, bukan struct

```go
// Interface
type CategoryRepository interface { ... }

// Struct (implementasi)
type CategoryRepositoryImpl struct {}

// Constructor return interface
func NewCategoryRepository() CategoryRepository {
    return &CategoryRepositoryImpl{}
}
```

### 3. Error Handling Strategy: Panic + Recover

Project ini menggunakan pola **panic-recover** untuk error handling:

```
Error terjadi → panic() → recover() di CommitOrRollback → re-panic → PanicHandler → ErrorHandler
```

| Kelebihan | Kekurangan |
|---|---|
| Kode lebih bersih (tidak perlu `if err != nil` berulang) | Tidak idiomatic Go |
| Centralized error handling | Sulit di-debug jika stack trace panjang |
| Automatic rollback transaction | Performa sedikit lebih lambat |

> **Best Practice:** Untuk production, lebih disarankan menggunakan explicit error return dan custom error types.

### 4. Manual Dependency Injection

Semua dependency di-wire secara manual di `main.go`:

```go
repository → service → controller → router → middleware → server
```

Ini membuat **dependency graph** terlihat jelas dan mudah dipahami.

### 5. Integration Test dengan `httptest`

- `httptest.NewRequest` → membuat fake HTTP request
- `httptest.NewRecorder` → menangkap HTTP response tanpa network call
- Test berjalan **tanpa** menjalankan HTTP server sesungguhnya
- Menggunakan database test terpisah → isolasi data

### 6. Transaction Management Pattern

```go
tx, err := service.DB.Begin()
helper.PanicIfError(err)
defer helper.CommitOrRollback(tx)  // ← otomatis commit/rollback
```

`defer` + `recover()` digunakan untuk memastikan transaksi **selalu** di-commit atau di-rollback, tanpa menulis try-catch pattern secara manual.

---

## 📝 Catatan Tambahan

- **API Spec** tersedia di `apispec.yml` (OpenAPI 3.0) — bisa di-import ke Swagger UI atau Postman
- **HTTP Test File** (`test.http`) bisa dijalankan langsung di IDE yang mendukung REST Client (GoLand, VS Code + REST Client extension)
- Server binary sudah ter-compile dan tersedia di file `server` (root directory)

---

> **Last Updated:** Maret 2026
