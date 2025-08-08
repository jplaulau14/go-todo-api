2025-08-08: Implemented graceful shutdown using http.Server with timeouts and SIGINT/SIGTERM handling. Ran tests and lints (all passing).
2025-08-08: Switched to structured logging via log/slog in server startup and shutdown paths. Ran tests and lints (all passing).
2025-08-08: Added Makefile with lint/test targets and GitHub Actions CI workflow running vet and tests. All passing locally.
2025-08-08: Added Dockerfile (multi-stage, distroless) and .dockerignore. Built image locally.
2025-08-08: Added minimal OpenAPI spec (openapi.yaml) documenting health and todo endpoints.
2025-08-08: Updated OpenAPI version to 3.1.0 to satisfy schema linter.
