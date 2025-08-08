2025-08-08: Implemented graceful shutdown using http.Server with timeouts and SIGINT/SIGTERM handling. Ran tests and lints (all passing).
2025-08-08: Switched to structured logging via log/slog in server startup and shutdown paths. Ran tests and lints (all passing).
2025-08-08: Added Makefile with lint/test targets and GitHub Actions CI workflow running vet and tests. All passing locally.
2025-08-08: Added Dockerfile (multi-stage, distroless) and .dockerignore. Built image locally.
2025-08-08: Added minimal OpenAPI spec (openapi.yaml) documenting health and todo endpoints.
2025-08-08: Updated OpenAPI version to 3.1.0 to satisfy schema linter.
2025-08-08: Enabled CORS using github.com/rs/cors to allow Swagger UI access from a different origin. Tests and lints passing.
2025-08-08: Added support for /todos without trailing slash and tests. Formatted, vetted, and all tests passing.
2025-08-08: Added dev Docker setup with hot reload using Air (Dockerfile.dev, docker-compose.yml, .air.toml) and Makefile targets (dev-up/dev-down). Verified container starts and watches files.
2025-08-08: Added Swagger UI service to docker-compose, serving openapi.yaml at http://localhost:8081.
2025-08-08: Implemented JSON error responses and structured logging across handlers. Added request logging and panic recovery middleware. All tests and lints passing.
2025-08-08: Removed unnecessary comments from server, middleware, and todo package files. Formatted, vetted, and tests passing.
2025-08-08: Added TODO.md roadmap outlining discrete PR-sized tasks to productionize the API.
2025-08-08: Added PostgreSQL service to docker-compose with persistent volume, env defaults, and healthcheck.
2025-08-08: Added goose migrations with Makefile targets (migrate-up/down). Created initial todos table migration and applied successfully to local Postgres.
2025-08-08: Implemented Postgres-backed repository and env-based repo selection (DB_DSN). Added integration test gated by TEST_DB_DSN. Tests and lints passing.
2025-08-08: Added Makefile target `test-integration` that boots DB, runs migrations, and executes tests with TEST_DB_DSN.
2025-08-08: Added /readyz endpoint that pings DB if configured (else returns OK). Added unit test for no-DB case. All tests and lints passing.
2025-08-08: Centralized configuration in internal/config (PORT, DB_DSN, LOG_LEVEL, ALLOWED_ORIGINS) with validation. Refactored server to use config. Added config unit tests. All tests and lints passing.
2025-08-08: Added request ID middleware (X-Request-ID) with propagation into logs and JSON error responses. Updated handlers and recovery. All tests and lints passing.
2025-08-08: Restricted CORS by environment via config (ENV). Disallow wildcard origins in prod. Added config tests. All tests and lints passing.
2025-08-08: Added request validation for JSON endpoints: enforce Content-Type, 1MB body limit, and unknown fields rejection. Added tests. All tests and lints passing.
2025-08-09: Standardized JSON error model {code, string, message, request_id, status}. Updated handlers and OpenAPI responses. Tests and lints passing.
