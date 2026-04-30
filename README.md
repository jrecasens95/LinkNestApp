# LinkNest

LinkNest is a private, self-hosted URL shortener for personal use, teams, newsrooms, events, and campaigns.

This repository is planned as a full-stack open-source project, but the current phase contains only the backend MVP.

## Stack

- Backend: Go + Fiber
- ORM: GORM
- Database: PostgreSQL
- Deployment: Docker Compose
- Frontend and auth: planned for later phases

## Backend MVP

Current features:

- PostgreSQL connection through GORM
- `ShortLink` model with automatic migration
- `POST /api/links` to create short links
- `GET /:code` to redirect short links
- `GET /api/health` health endpoint
- URL validation for `http://` and `https://`
- 6-character unique short codes
- Basic CORS
- Simple request logs

## Project Structure

```text
backend/
  cmd/server/main.go
  internal/config/
  internal/database/
  internal/models/
  internal/handlers/
  internal/services/
  internal/utils/
  .env.example
  Dockerfile
  go.mod
  go.sum
docker-compose.yml
```

## Run With Docker Compose

```bash
docker compose up --build
```

The API will be available at:

```text
http://localhost:4000
```

Health check:

```bash
curl http://localhost:4000/api/health
```

Create a short link:

```bash
curl -X POST http://localhost:4000/api/links \
  -H "Content-Type: application/json" \
  -d '{"original_url":"https://example.com/some/long/path","title":"Example"}'
```

Example response:

```json
{
  "code": "aB12xY",
  "short_url": "http://localhost:4000/aB12xY"
}
```

Open the returned `short_url` to redirect to the original URL.

## Run Backend Locally

Start PostgreSQL with Docker Compose:

```bash
docker compose up postgres
```

Configure the backend:

```bash
cp backend/.env.example backend/.env
```

Run the API:

```bash
cd backend
go mod tidy
go run ./cmd/server
```

## Environment Variables

The backend reads these values from environment variables or `backend/.env`:

| Name | Description | Example |
| --- | --- | --- |
| `DATABASE_URL` | PostgreSQL connection string | `postgres://linknest:linknest@localhost:5432/linknest?sslmode=disable` |
| `BASE_URL` | Public base URL used to build short URLs | `http://localhost:4000` |
| `PORT` | API port | `4000` |

## API

### `GET /api/health`

Returns:

```json
{
  "status": "ok"
}
```

### `POST /api/links`

Request body:

```json
{
  "original_url": "https://example.com",
  "title": "Optional title"
}
```

Validation:

- `original_url` is required.
- `original_url` must start with `http://` or `https://`.

### `GET /:code`

Redirects to the original URL.

Responses:

- `302` when the link is valid.
- `404` when the link does not exist.
- `410` when the link is inactive or expired.

## License

License to be decided.
