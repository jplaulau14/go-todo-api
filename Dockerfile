# syntax=docker/dockerfile:1

FROM golang:1.23-alpine AS builder
WORKDIR /app
ENV CGO_ENABLED=0

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -ldflags="-s -w" -o /out/server ./cmd/server

FROM gcr.io/distroless/static-debian12:nonroot
ENV PORT=8080
EXPOSE 8080
COPY --from=builder /out/server /server
ENTRYPOINT ["/server"]


