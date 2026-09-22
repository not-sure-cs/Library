# Library Management & PDF Ingestion API Service

This is my attempt at writing a robust RESTful API built in Go. Essentially, it's my own rendering of a LibGen backend that helps one manage digital library collections, extract PDF metadata automatically, render cover images, control access by user roles, and store files using Cloudflare R2 and PostgreSQL.

---

## 1. Project Overview

The Library API service offers a complete backend for digital libraries and streaming assets. It removes the need for manual data entry by automatically checking uploaded files.

When an authorized user uploads a document:
1. **Deduplication Check**: Computes a SHA-256 fingerprint of the binary and queries the database to prevent duplicate document storage.
2. **Metadata Extraction**: Inspects PDF headers and catalogs to extract document title, author, producer software, page count, PDF version, and subject keywords.
3. **Automated Cover Page Rendering**: Isolates the first page of the PDF in memory, converts the vector PDF page to high-quality JPEG raster imagery, and stores it in object storage.
4. **Cloud Storage Synchronization**: Transmits both the full document and rendered cover art concurrently to Cloudflare R2 (S3-compatible object storage).
5. **Relational Consistency**: Persists normalized author, book, and relationship records within PostgreSQL via type-safe generated queries.

---

## 2. Architecture & Design Principles

```
                         ┌─────────────────────────────────────────┐
                         │               Client (HTTP)             │
                         └────────────────────┬────────────────────┘
                                              │
                                     (Cookie Session)
                                              ▼
                         ┌─────────────────────────────────────────┐
                         │         Go HTTP ServeMux Router         │
                         └──────┬──────────────┬─────────────┬─────┘
                                │              │             │
                    JSON Middleware      Auth Session    Role Check
                                │              │             │
                                └──────┬───────┴─────────────┘
                                       │
                                       ▼
                       ┌───────────────────────────────┐
                       │          API Handlers         │
                       │  (Signup, Login, Books CRUD)  │
                       └──────┬─────────────────┬──────┘
                              │                 │
              Extraction & Image Pipelines      │ Database Queries (sqlc)
                              │                 │
                              ▼                 ▼
                 ┌───────────────────────┐  ┌───────────────────────┐
                 │ Cloudflare R2 Storage │  │  PostgreSQL Database  │
                 │  (Books & Cover Art)  │  │ (Catalog, RBAC, Keys) │
                 └───────────────────────┘  └───────────────────────┘
```

The system embraces a clean layered architecture:
- **`cmd/server`**: The application composition root responsible for environment resolution, database connection pool provisioning, Cloudflare R2 client initialization, and graceful shutdown handling.
- **`internal/api`**: HTTP transport layer containing routing definitions, authentication middleware, role enforcement decorators, and JSON serialization helpers.
- **`internal/extraction`**: Low-level stream readers and conversion engines that perform MIME sniffing, PDF parsing, memory-buffered rendering, and hashing.
- **`internal/storage`**: Cloudflare R2 and AWS S3 SDK wrapper managing multipart uploads, presigned download link generation, and asset deletion.
- **`internal/database`**: Data access layer generated with `sqlc`, combining compile-time SQL verification with idiomatic Go interfaces.

---

## 3. Technology Stack

- **Language**: Go 1.22+ (utilizing standard library pattern matching in `net/http`)
- **Database**: PostgreSQL 14+ with pgx/v5 driver (`github.com/jackc/pgx/v5`)
- **Query Generation**: `sqlc` (type-safe, boilerplate-free SQL compilation)
- **Database Migrations**: `goose` (versioned SQL migrations)
- **Object Storage**: Cloudflare R2 via AWS SDK for Go v2 (`github.com/aws/aws-sdk-go-v2`)
- **Authentication & Sessions**: `github.com/gorilla/sessions` with cookie-backed encrypted storage
- **Password Security**: `golang.org/x/crypto/bcrypt` (salted adaptive hashing)
- **PDF & Image Processing**: Native in-memory PDF extraction and image buffering

---

## 4. Key Features & Capabilities

### Automated PDF Ingestion & Cover Extraction
Uploads bypass manual categorization. The ingestion pipeline reads stream bytes to resolve structural PDF metadata (such as Author and Title tags) and renders the cover page directly into JPEG format. File reader streams are automatically rewound at each processing boundary to prevent truncation.

### Multi-Tier Role-Based Access Control (RBAC)
The service enforces granular permissions through database enum types and cookie session validation:
- **`member`**: Can query catalog entries, inspect book details, and generate temporary presigned download URLs.
- **`moderator`**: Possesses member capabilities plus authorization to edit book information, change author links, and manage catalog metadata.
- **`admin`**: Retains full system privileges, including book deletion and storage purge operations.

### Storage Safety & Deletion Inversion
When removing catalog entries, the application purges binary assets from Cloudflare R2 before dropping database foreign keys. This guarantees that failed storage requests do not leave untracked orphan files occupying cloud storage.

---

## 5. Directory Structure

```
.
├── cmd/
│   └── server/
│       └── main.go                 # Server entry point & dependency wiring
├── internal/
│   ├── api/
│   │   ├── Authedmiddleware.go     # Session authentication verification
│   │   ├── Rolemiddelware.go       # RBAC role validation decorator
│   │   ├── JSONmiddleware.go       # Global Content-Type header injection
│   │   ├── JSON.go                 # Standardized JSON response formatting
│   │   ├── handleSignUp.go         # User registration handler
│   │   ├── handleLogging.go        # Session login & credential check
│   │   ├── handleCreateBooks.go    # PDF upload, ingestion & cover creation
│   │   ├── handleGetBook.go        # Book retrieval & presigned URL delivery
│   │   ├── handleUpdateBook.go     # Book & author metadata updates
│   │   ├── handleDeleteBook.go     # Cascaded asset & record deletion
│   │   └── handleStatus.go         # Health check & uptime monitor
│   ├── database/
│   │   ├── models.go               # sqlc generated data models
│   │   ├── db.go                   # Database query instance constructor
│   │   ├── interface.go            # DBQueries mockable interface definition
│   │   ├── struct.go               # API response transfer objects (UserBook, Stream)
│   │   ├── func.go                 # Storage helpers, hash utilities & custom queries
│   │   ├── books.sql.go            # Generated book mutation routines
│   │   ├── bookuse.sql.go          # Generated book query routines
│   │   ├── users.sql.go            # Generated user creation queries
│   │   └── usersuse.sql.go         # Generated user retrieval queries
│   ├── extraction/
│   │   ├── func.go                 # Metadata extraction & SHA256 file hashing
│   │   └── conversion.go           # PDF cover page rendering to JPEG
│   └── storage/
│       ├── r2.go                   # Cloudflare R2 S3-compatible client
│       └── secret.go               # Credentials and bucket configuration model
├── sql/
│   ├── queries/                    # Raw SQL query templates for sqlc
│   │   ├── books.sql
│   │   ├── bookuse.sql
│   │   ├── users.sql
│   │   └── usersuse.sql
│   └── schema/                     # Sequential Goose migration scripts
│       ├── 001_Books.sql
│       ├── 002_Users.sql
│       ├── ...
│       ├── 013_BookUpdate.sql
│       └── 014_UserEmailUnique.sql
├── sqlc.yaml                       # sqlc code generation configuration
├── go.mod                          # Go module dependencies
└── readme.md                       # Documentation
```

---

## 6. Environment Configuration

Create a `.env` file in the project root with the following variables:

```ini
PORT=8080
ENV=development
DB_URL=postgres://username:password@localhost:5432/library_db?sslmode=disable
KEY=your-32-byte-base64-session-key

# Cloudflare R2 Credentials
R2_ACCOUNT_ID=your_cloudflare_account_id
R2_ACCESS_KEY_ID=your_r2_access_key
R2_SECRET_ACCESS_KEY=your_r2_secret_key
R2_BUCKET_NAME=your_bucket_name
```

---

## 7. Setup & Installation

### Prerequisites
- Go 1.22 or higher
- PostgreSQL running locally or accessible via network
- `sqlc` CLI (`go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest`)
- `goose` migration tool (`go install github.com/pressly/goose/v3/cmd/goose@latest`)

### Step 1: Clone Repository & Install Dependencies
```bash
git clone https://github.com/knibirdgautam/library.git
cd library
go mod download
```

### Step 2: Run Database Migrations
Apply all schema changes in sequence:
```bash
goose -dir sql/schema postgres "$DB_URL" up
```

### Step 3: Compile SQL Queries (Optional)
If you modified files in `sql/queries/`:
```bash
sqlc generate
```

### Step 4: Build and Run the Server
```bash
go run ./cmd/server
```

The server initializes on port `8080` (or the configured `PORT` variable) and outputs:
```
Port: 8080
Successfully connected to the database
Starting Server on Port: 8080
```

---

## 8. API Reference

All requests and responses use `application/json`, except file upload which uses `multipart/form-data`.

### 1. System Health
- **Endpoint**: `GET /status`
- **Auth Required**: No
- **Response** `200 OK`:
  ```json
  {
    "Uptime": 1420583300
  }
  ```

---

### 2. User Registration
- **Endpoint**: `POST /user/signup`
- **Auth Required**: No
- **Request Body**:
  ```json
  {
    "firstname": "Jane",
    "lastname": "Doe",
    "email": "jane.doe@example.com",
    "phone": "+1234567890",
    "password": "securepassword123"
  }
  ```
- **Responses**:
  - `201 Created`: User entity created.
  - `400 Bad Request`: Missing mandatory fields.
  - `409 Conflict`: User with this email already registered.

---

### 3. User Authentication (Login)
- **Endpoint**: `POST /user/login`
- **Auth Required**: No
- **Request Body**:
  ```json
  {
    "email": "jane.doe@example.com",
    "password": "securepassword123"
  }
  ```
- **Responses**:
  - `200 OK`: Session cookie set (`Set-Cookie: user-session=...`).
  - `401 Unauthorized`: Invalid credentials.

---

### 4. Upload & Ingest Book
- **Endpoint**: `POST /book`
- **Auth Required**: Yes (`member`, `moderator`, `admin`)
- **Headers**: `Content-Type: multipart/form-data`
- **Form Fields**:
  - `uploadFile`: Binary PDF file payload.
- **Responses**:
  - `201 Created`:
    ```json
    {
      "author": {
        "ID": "4e3caab6-e636-4b81-b443-397a6be358c8",
        "Name": "K. A. Stroud"
      },
      "book": {
        "ID": "e09395f9-2318-4677-95f3-1f946a0f15b0",
        "Name": "Advanced Engineering Mathematics",
        "FilePath": "Assets/Books/Type-Book-18bb6a05.pdf",
        "CoverPath": "Assets/Covers/Type-Book-18bb6a05.jpg",
        "PageCount": { "Int32": 1248, "Valid": true },
        "PdfVersion": { "String": "1.4", "Valid": true }
      }
    }
    ```
  - `400 Bad Request`: Unreadable MIME or invalid upload part.
  - `409 Conflict`: Binary already exists (duplicate SHA-256 match).

---

### 5. Retrieve Book by ID
- **Endpoint**: `GET /book/{id}`
- **Auth Required**: Yes
- **Responses**:
  - `200 OK`:
    ```json
    {
      "book": {
        "id": "e09395f9-2318-4677-95f3-1f946a0f15b0",
        "book": "Advanced Engineering Mathematics",
        "author": "K. A. Stroud",
        "cover_path": "Assets/Covers/Type-Book-18bb6a05.jpg",
        "page_count": { "Int32": 1248, "Valid": true }
      },
      "link": "https://library-assets.r2.cloudflarestorage.com/Assets/Books/Type-Book-18bb6a05.pdf?X-Amz-Signature=..."
    }
    ```
  - `404 Not Found`: Book ID does not exist.
  - `422 Unprocessable Entity`: Malformed UUID string.

---

### 6. Update Book Metadata
- **Endpoint**: `PUT /book/{id}`
- **Auth Required**: Yes (`moderator`, `admin`)
- **Request Body**:
  ```json
  {
    "title": "Advanced Engineering Mathematics - 8th Edition",
    "author": "K. A. Stroud, Dexter Booth",
    "isbn": "978-1352010275",
    "subject": "Engineering Mathematics",
    "producer": "Macmillan",
    "page_count": 1250
  }
  ```
- **Responses**:
  - `200 OK`: Updated book record.
  - `403 Forbidden`: Insufficient role permissions.
  - `404 Not Found`: Book does not exist.

---

### 7. Delete Book
- **Endpoint**: `DELETE /book/{id}`
- **Auth Required**: Yes (`moderator`, `admin`)
- **Responses**:
  - `200 OK`: `{"message": "Book deleted"}`
  - `403 Forbidden`: Insufficient role permissions.
  - `404 Not Found`: Target record not found.

---

## 9. Testing & Verification Guide

Test the end-to-end flow from terminal using `curl`:

```bash
# 1. Register a test user
curl -s -X POST http://localhost:8080/user/signup \
  -H "Content-Type: application/json" \
  -d '{
    "firstname": "Alex",
    "lastname": "Smith",
    "email": "alex.smith@example.com",
    "password": "Password123!"
  }'

# 2. Login and capture session cookies
curl -s -X POST http://localhost:8080/user/login \
  -H "Content-Type: application/json" \
  -c cookies.txt \
  -d '{
    "email": "alex.smith@example.com",
    "password": "Password123!"
  }'

# 3. Upload a PDF document
curl -s -X POST http://localhost:8080/book \
  -b cookies.txt \
  -F "uploadFile=@sample.pdf"

# 4. Fetch the book and presigned stream URL (substitute returned UUID)
curl -s http://localhost:8080/book/<BOOK_UUID> \
  -b cookies.txt
```

---

## 10. License & Maintenance

Developed by not-sure-cs. Contributions are welcome via pull requests. I am still learning and request everyone to engage with the project and suggest improvements and new features to implement. I am a willing learner, and any criticism or critique is also Welcome.
