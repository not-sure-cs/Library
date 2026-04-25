### Project Summary: Golang Library API Service

**Architecture Overview**
This project evaluates as a statically compiled, monolithic REST-style web service API architected in Go. It employs a traditional layered architecture, distinctly separating the HTTP transport layer (`internal/api`) from the database execution layer (`internal/database`). The application bootstrapper (`cmd/server/main.go`) handles dependency injection, initializing the database connection pool, loading environment configurations, and wrapping the request multiplexer with global middleware.

**Technology Stack**
*   **Language:** Go (Golang). Leveraging modern Go routing syntax introduced in Go 1.22+ (e.g., pattern matching `GET /book/{id}`).
*   **HTTP Server:** The native `net/http` standard library, eliminating the need for heavyweight third-party web frameworks like Gin or Fiber.
*   **Datastore:** PostgreSQL, acting as the primary relational persistence layer. 
*   **Database Driver:** `github.com/jackc/pgx/v5/stdlib`, a high-performance PostgreSQL driver and toolkit for Go, favored for its native advanced type support.
*   **ORM & Query Builder:** `sqlc` (v1.30.0). Instead of a runtime reflection-heavy ORM like GORM, the project utilizes `sqlc` to parse SQL query files and dynamically generate type-safe idiomatic Go data-access code resulting in performant, predictable database interactions free of interface/type-assertion clutter.
*   **Configuration:** 12-Factor App compliant using `github.com/joho/godotenv` to map `.env` file configurations onto system environment variables (specifically for `PORT` and `DB_URL` mapping).
*   **Identifier Standard:** `github.com/google/uuid` utilizing UUID v4 to maintain globally unique pseudo-random 128-bit identifiers, preventing integer collision and sequential ID discovery attacks.
*   **Testing Strategy:** Test-driven methodologies utilizing the `net/http/httptest` package. The database layer is obfuscated via a strict `database.DBQueries` interface, permitting full manual mock struct overrides. This guarantees rapid, isolated unit testing across all HTTP handler closures without a running testing database, currently achieving 83% line code coverage within the API domain.

**Current Scope & Feature Set**
At present, the service domain handles fundamental inventory and metadata manipulation rather than circulation state. The scope focuses on atomic transactions ensuring database normalization.

*   **RESTful Mutators (Books):** Supports creating, reading, updating (PUT overrides), and deleting book records via UUID mapping parameters. Incoming payloads are marshaled from JSON via `json.NewDecoder`. Content is actively protected against unexpected malformation.
*   **Author Graph Linking:** Incorporates advanced association logic. When creating a book referencing an unknown author string, the API transparently performs an upsert execution, generating a new `Author` entity and constructing an entry in the aggregate many-to-many intersection table (`book_authors`). Returns composite nested JSON objects encapsulating the entire normalized relationship graph (`database.Linked`).
*   **Referential Integrity Constraints:** System deletes (`DELETE /book/{id}`) correctly fire a prerequisite `UnlinkBook` transaction against the `book_authors` junction schema to prevent orphaned records and cascading constraint faults. 
*   **Query Operations:** Capable of fulfilling complex read operations such as fetching all associated book items specifically belonging to a single UUID Author entity.
*   **Middleware Pipeline:** Injects global `Content-Type: application/json` headers preceding the ServeHTTP call chain, guaranteeing standardized client-side response parsing. Included generic error-formatter utilities abstract conditional error logic (400 Client error vs. 500 internal server masking) downstream.
