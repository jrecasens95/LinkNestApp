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
- `GET /api/links` to list links
- `GET /api/links/:id` to inspect one link
- `GET /api/links/:id/stats` to inspect click statistics
- `PATCH /api/links/:id` to edit title or activation status
- `DELETE /api/links/:id` to remove a link
- `GET /:code` to redirect short links
- `GET /api/health` health endpoint
- URL validation for `http://` and `https://`
- Strong URL validation with private/internal address blocking
- Rate limit on `POST /api/links`
- Security headers middleware
- Optional private mode using `X-API-Key`
- 6-character unique short codes
- Basic CORS
- Simple request logs
- Responsive SaaS-style landing page
- Form to create short links
- Copy button for generated short URLs
- Private dashboard without auth for single-instance administration
- Link detail page with title editing, active/inactive toggle, clicks, and delete
- Click events recorded on redirect with anonymized IP, user agent, referer, and timestamp
- Basic stats view with total clicks, latest clicks, referers, and click dates

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

Frontend routes:

- `/` creates a new short link.
- `/dashboard` lists all links.
- `/dashboard/links/:id` opens the link detail page.

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
| `PRIVATE` | Require `X-API-Key` for creating links | `false` |
| `API_KEY` | API key used when `PRIVATE=true` | `change-me` |
| `MAX_URL_LENGTH` | Maximum accepted destination URL length | `2048` |
| `BLACKLISTED_DOMAINS` | Comma-separated blocked domains | `localhost,local,internal,metadata.google.internal` |

The frontend reads this value from `apps/web/.env`:

| Name | Description | Example |
| --- | --- | --- |
| `VITE_API_URL` | Backend API URL | `http://localhost:4000` |
| `VITE_API_KEY` | Optional key sent to `POST /api/links` for private deployments | `change-me` |

## Security

The API applies basic security controls:

- `POST /api/links` is rate limited.
- Destination URLs must be valid absolute `http://` or `https://` URLs.
- URLs with credentials, control characters, invalid ports, or excessive length are rejected.
- `localhost`, loopback addresses, private networks, link-local ranges, multicast, unspecified addresses, and configured blacklisted domains are blocked.
- New links resolve their hostname before being saved to reduce internal-address redirects.
- Existing links are validated again before redirecting.
- Security headers are set globally.

When `PRIVATE=true`, creating links requires:

```bash
X-API-Key: your-api-key
```

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
- `original_url` must not point to localhost, private/internal IP ranges, single-label hosts, or blacklisted domains.
- If `PRIVATE=true`, requests must include `X-API-Key`.

### `GET /:code`

Redirects to the original URL.

Responses:

- `302` when the link is valid.
- `404` when the link does not exist.
- `410` when the link is inactive or expired.

### `GET /api/links`

Lists all links ordered by newest first.

### `GET /api/links/:id`

Returns one link by internal ID.

### `GET /api/links/:id/stats`

Returns basic click statistics:

```json
{
  "total_clicks": 12,
  "recent_clicks": [
    {
      "id": 1,
      "user_agent": "Mozilla/5.0",
      "referer": "https://example.com",
      "ip_address": "192.168.1.0",
      "created_at": "2026-04-30T06:30:00Z"
    }
  ],
  "referers": [
    {
      "referer": "https://example.com",
      "count": 4
    }
  ]
}
```

### `PATCH /api/links/:id`

Updates editable fields:

```json
{
  "title": "Updated title",
  "is_active": true
}
```

### `DELETE /api/links/:id`

Deletes a link.

## License

MIT. See [LICENSE](LICENSE).
