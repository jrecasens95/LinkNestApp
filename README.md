# LinkNest

LinkNest is a private, self-hosted URL shortener for personal use, teams, newsrooms, events, and campaigns.

## Stack

- Backend: Go + Fiber
- ORM: GORM
- Database: PostgreSQL
- Frontend: React + Vite + Tailwind
- Deployment: Docker Compose
- Auth: planned for later phases

## MVP

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
- Responsive SaaS-style landing page
- Form to create short links
- Copy button for generated short URLs

## Project Structure

```text
apps/api/
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
apps/web/
  src/
    api/
    components/
    pages/
    App.tsx
    main.tsx
  .env.example
docker-compose.yml
```

## Run With Docker Compose

```bash
docker compose up --build
```

The web app will be available at:

```text
http://localhost:3000
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
cp apps/api/.env.example apps/api/.env
```

Run the API:

```bash
cd apps/api
go mod tidy
go run ./cmd/server
```

## Run Frontend Locally

Configure the frontend:

```bash
cp apps/web/.env.example apps/web/.env
```

Run the web app:

```bash
cd apps/web
npm install
npm run dev
```

## Environment Variables

The backend reads these values from environment variables or `apps/api/.env`:

| Name | Description | Example |
| --- | --- | --- |
| `DATABASE_URL` | PostgreSQL connection string | `postgres://linknest:linknest@localhost:5432/linknest?sslmode=disable` |
| `BASE_URL` | Public base URL used to build short URLs | `http://localhost:4000` |
| `PORT` | API port | `4000` |

The frontend reads this value from `apps/web/.env`:

| Name | Description | Example |
| --- | --- | --- |
| `VITE_API_URL` | Backend API URL | `http://localhost:4000` |

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
