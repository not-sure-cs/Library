# Project Summary: Golang Library & Inventory API Service

## 1. Executive Overview
This project represents a high-performance, statically compiled, monolithic REST-style web service architected in Go (version 1.25+). Designed with the core tenets of the "Go Way"—simplicity, explicitness, and performance—the service provides a robust backend for managing a library's digital assets, physical inventory, and user ecosystem. 

Unlike many modern web applications that rely on heavy third-party frameworks like Gin, Echo, or Fiber, this service is built almost entirely on the Go standard library (`net/http`). This decision minimizes dependency bloat, reduces the attack surface for security vulnerabilities, and ensures long-term maintainability through standard-compliant code. The project serves as a demonstration of a modern layered architecture, distinctly separating the transport layer from the persistence layer while maintaining strict type-safety across the entire stack.

---

## 2. Architectural Philosophy: The Layered Approach
The codebase is organized into a clean, hierarchical structure that facilitates ease of testing, modular development, and clear separation of concerns.

### 2.1 Bootstrapping (`cmd/server/main.go`)
The application entry point is responsible for the "Wiring" of the service. It performs several critical startup tasks:
- **Environment Configuration:** Utilizing `godotenv` to load 12-factor app compliant variables such as `PORT` and `DB_URL`.
- **Database Initialization:** Establishing a connection pool via `sql.Open` using the `pgx/v5` driver, followed by a health check (`Ping`) to ensure service readiness.
- **Routing & Middleware Registration:** Defining the API surface using the modern `http.NewServeMux` (leveraging Go 1.22+ pattern matching) and wrapping it in a global middleware pipeline.
- **Lifecycle Management:** Setting aggressive timeouts (Read, Write, Idle) on the `http.Server` to prevent resource exhaustion attacks and ensure high availability.

### 2.2 Transport Layer (`internal/api/`)
The API layer handles the HTTP request-response cycle. It is composed of functional handlers that ingest JSON payloads, validate input, and coordinate with the database layer. 
- **JSON Standardization:** Utilities in `JSON.go` provide a consistent interface for success (`RespondWithJSON`) and failure (`RespondWithError`), ensuring all API responses adhere to a predictable structure.
- **Middleware Pipeline:** 
    - `JSONMiddleware`: Injects standard headers and ensures consistent content negotiation.
    - `AuthedMiddleware`: Provides a secure perimeter, validating session state before allowing access to protected resources.

### 2.3 Persistence Layer (`internal/database/`)
The database layer is built around `sqlc`, a tool that generates type-safe Go code from raw SQL queries. This approach offers several advantages over traditional ORMs:
- **Performance:** Zero reflection at runtime; queries are compiled into idiomatic Go.
- **Reliability:** SQL syntax errors are caught at build time rather than runtime.
- **Maintainability:** The "Source of Truth" remains standard SQL files, making it easy for DBAs and developers to collaborate.

---

## 3. Technology Stack & Tooling

### 3.1 Core Language & Runtime
- **Golang (v1.25.5):** Utilizing the latest features of the Go runtime, including enhanced routing patterns and optimized garbage collection.
- **Standard Library Focus:** Heavy reliance on `net/http`, `encoding/json`, and `context`.

### 3.2 Data Management
- **PostgreSQL:** Chosen for its industry-leading reliability, ACID compliance, and advanced support for relational mapping.
- **SQLC (v1.30.0):** Used to generate the `DBQueries` interface, which abstracts the database operations and allows for easy mocking in unit tests.
- **PGX (v5):** A high-performance PostgreSQL driver that supports native types and optimized connection pooling.

### 3.3 Security & Identity
- **Google UUID:** Utilizing version 4 UUIDs for all primary keys to prevent ID enumeration attacks and simplify data synchronization across distributed systems.
- **Gorilla Sessions:** Secure, cookie-based session management for maintaining user state.
- **Bcrypt:** Implementing the Blowfish-based hashing algorithm for password security, ensuring that even in the event of a database breach, user credentials remain protected.

---

## 4. Deep Dive: Core Domains & Features

### 4.1 User Ecosystem & Authentication
The service implements a complete authentication flow:
- **Sign-Up:** Collects user metadata (first name, last name, email, phone) and persists a secure bcrypt hash of the password into a separate `secrets` table to maintain data isolation.
- **Login:** Performs a look-up by email, verifies the bcrypt hash, and initializes a secure session via a `CookieStore`.
- **Session Management:** Sessions are configured with `HttpOnly`, `Secure`, and `SameSite` flags to mitigate XSS and CSRF risks.

### 4.2 Library Inventory (Books & Authors)
The core of the application revolves around the relationship between books and their creators.
- **Many-to-Many Relationships:** Handled via a junction table `book_authors`, allowing a single book to have multiple authors and vice versa.
- **Intelligent Upserts:** The `handleCreateBooks` logic includes a sophisticated "Check-then-Create" pattern. If a book creation request references an author not currently in the database, the system transparently creates the author record first before linking them.
- **Referential Integrity:** SQL-level constraints (`ON DELETE CASCADE`) ensure that when a book or author is removed, the junction table remains clean, preventing orphaned records.

### 4.3 Asset & File Handling
Unique to this implementation is the ability to handle physical/digital asset uploads:
- **Multipart Form Processing:** The service parses `multipart/form-data` to ingest both JSON metadata and binary file content (e.g., PDF or EPUB files).
- **Concurrency in Action:** During book creation, the service spawns a goroutine to concurrently count existing books while the file is being processed, demonstrating Go's powerful concurrency primitives (channels and goroutines).
- **Storage Strategy:** Files are saved to a version-controlled path (`Public/Books/`) with unique, timestamp-based filenames to prevent collisions.

---

## 5. API Reference & Request Lifecycle

### 5.1 Public Endpoints
- `GET /status`: Returns the system uptime and server health metrics.
- `POST /user/signup`: Registers a new user account.
- `POST /user/login`: Authenticates a user and establishes a session.

### 5.2 Protected Endpoints (Requires Session)
- `POST /book`: Ingests metadata and a file to create a new book entry.
- `GET /book/{id}`: Retrieves detailed metadata for a specific book by its UUID.
- `PUT /book/{id}`: Updates existing book information (supports partial overrides).
- `DELETE /book/{id}`: Removes a book and its associated author links from the system.
- `GET /author/{id}/books`: Retrieves a listing of all books associated with a specific author ID.

### 5.3 The Request Journey
1. **Entry:** A client sends a request. The `wrappedMux` (JSON Middleware) catches it, setting the `Content-Type: application/json` header.
2. **Auth Check:** If the route is protected, `AuthedMiddleware` checks for a valid session. If missing, it returns a `401 Unauthorized` via `RespondWithError`.
3. **Routing:** `http.NewServeMux` matches the path and method (e.g., `PUT /book/{id}`).
4. **Processing:** The handler decodes the request body, performs business logic, and interacts with the database via the `Queries` interface.
5. **Exit:** The handler calls `RespondWithJSON`, which encodes the response and flushes the buffer to the client.

---

## 6. Security Implementation Details
- **Password Salting:** Bcrypt handles salting automatically, preventing rainbow table attacks.
- **SQL Injection Prevention:** By using `sqlc` and prepared statements, the application is inherently protected against SQL injection attacks.
- **Data Sanitization:** Input is rigorously decoded into strongly typed structs, ensuring that malformed JSON or unexpected fields are rejected early in the request lifecycle.
- **Session Protection:** All session cookies are signed with an environment-provided secret key and restricted to `SameSite` strict mode.

---

## 7. Testing & Quality Assurance Strategy
The project is designed for high testability:
- **Interface-Based Database Access:** The `database.DBQueries` interface allows the API handlers to be tested in isolation by injecting a mock database.
- **Standard Tooling:** Utilizing `net/http/httptest` for end-to-end handler verification.
- **Coverage Goal:** The architecture aims for >80% code coverage, focusing on critical paths such as authentication, author-book linking, and file upload validation.

---

## 8. Future Roadmap
The current architecture is highly extensible, with plans for:
- **Circulation Management:** Adding tables and logic for book loans, returns, and waitlists.
- **Advanced Search:** Implementing Full-Text Search (FTS) using PostgreSQL's native capabilities for searching through book titles and descriptions.
- **Dockerization:** Standardizing the deployment environment via multi-stage Docker builds.
- **API Documentation:** Integrating Swagger/OpenAPI for interactive documentation.

---
*Last Updated: May 2026*
*This summary serves as the primary technical documentation for the Library API Service.*

