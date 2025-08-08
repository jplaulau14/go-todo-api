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
