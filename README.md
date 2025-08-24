# Auto Message

This project is a **service written in Go**.
It sends **pending messages** from database to a webhook in period of time.
It should not send same message again. Also it gives some endpoints for checking and reporting.

---

## Getting Started

### 1. Environment Variables (.env)

First copy the example env file:

```bash
cp env.example .env
```

Then open `.env` and change values if you need:

```env
HTTP_PORT=8081

POSTGRES_USER=auto
POSTGRES_PASSWORD=messager
POSTGRES_DB=automessager
POSTGRES_PORT=5432
POSTGRES_HOST=db

REDIS_HOST=redis
REDIS_PORT=6379
ADMINER_PORT=1110

PERIOD=120
BATCH_SIZE=2
EXTERNAL_API_URL=https://webhook.site/xxxx-xxxx-xxxx
```

---

### 2. Run Services

Build docker image and run all services:

```bash
docker compose --env-file .env up -d --build
```

Now you will see these services:

* **App** → [http://localhost:8081](http://localhost:8081)
* **Swagger UI** → [http://localhost:8081/swagger/index.html](http://localhost:8081/swagger/index.html)
* **Adminer** → [http://localhost:1110](http://localhost:1110)


## Development

### 1. sqlc

**sqlc** make Go code from SQL files. Safer and type-checked.

Install:

```bash
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```

Add to PATH:

```bash
export PATH="$(go env GOPATH)/bin:$PATH"
```

Generate:

```bash
sqlc generate
```

Code will be in `internal/storage`.

---

### 2. Swagger

**swag** create Swagger / OpenAPI doc from comments in Go handlers.

Install:

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

Add to PATH:

```bash
export PATH="$(go env GOPATH)/bin:$PATH"
```

Generate docs:

```bash
swag init --dir ./cmd,./internal --output ./docs
```

After run:

* `docs/docs.go` → generated code
* Swagger UI → [http://localhost:8081/swagger/index.html](http://localhost:8081/swagger/index.html)

---

## URLs

* Health check: [http://localhost:8081/ping](http://localhost:8081/ping)
* Start listener: `GET http://localhost:8081/api/start-listener`
* Stop listener: `GET http://localhost:8081/api/stop-listener`
* List Sent messages: `GET http://localhost:8081/api/messages/sent`
* Swagger: [http://localhost:8081/swagger/index.html](http://localhost:8081/swagger/index.html)
* Adminer: [http://localhost:1110](http://localhost:1110)
