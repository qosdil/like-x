# like-x backend: User Service

A Go microservice, uses Fiber (HTTP), pgx (PostgreSQL), and a clean repository-service-transport layered architecture.

## Features

- `POST /v1/users/sign-up`
- Validation for `full_name` and `password`
- PostgreSQL persistence via `pgx` and `pgxpool`
- Structured packages:
  - `model` (domain models)
  - `repository` (DB repository)
  - `service` (business logic)
  - `transport/http` (HTTP handlers)

## Requirements

- Go 1.26 (or latest supported)
- PostgreSQL (DB setup via environment variables)

## Environment variables

- `POSTGRES_USER`
- `POSTGRES_PASSWORD`
- `POSTGRES_HOST`
- `POSTGRES_PORT`
- `POSTGRES_DB_NAME`
- `POSTGRES_SSL_MODE` (e.g. `disable`)
- `HTTP_SERVER_PORT` (e.g. `3000`)

## Run

```bash
cd backend/user
go run .
```

If using `.env` file, existing code loads it via `github.com/joho/godotenv`.

## API

Request:

```bash
curl -X POST http://localhost:3000/v1/users/sign-up \
  -H 'Content-Type: application/json' \
  -d '{"full_name":"John Doe","password":"secret123"}'
```

Response:

```json
{"id":"<public-id>"}
```

## Docker

Build the image:

```bash
docker build -t like-x-user .
```

Run the container:

```bash
docker run --rm --env-file .env -p 3000:3000 --name like-x-post-latest like-x-user:latest
```

## Tests

```bash
go test -cover ./... -count=1 -v
```

## Simple, Local Load Testing

Run the HTTP server with constrained resources for a lightweight load test:

```sh
GOMAXPROCS=0.05 GOMEMLIMIT=128MiB go run main.go
```

### Install vegeta

Let's use `vegeta` for this load testing. Install via Homebrew (macOS/Linux), Go install, or download a binary:

```sh
# Homebrew (recommended on macOS)
brew install vegeta

# Go install (with Go 1.20+)
go install github.com/tsenart/vegeta/v12@latest

# Verify installation
vegeta -version
```

Generate 1,000 RPS for 30s with unique JSON input per request:

```bash
seq 1 30000 \
  | xargs -I {} jq --arg full_name "John Doe "{} --arg password "mysecret"{} -ncM '{
    method: "POST", url: "http://localhost:3000/v1/users/sign-up",
    body: {"full_name": $full_name, "password": $password} | @base64,
    header: {"Content-Type": ["application/json"]}
  }' \
  | vegeta attack -format=json -rate=1000 -duration=30s -timeout=60s \
  | vegeta report
```
