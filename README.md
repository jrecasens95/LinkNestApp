# link-nest

Generated with `project-starter`.

## Stack

- Frontend: react-vite
- Backend: go-fiber
- Database: postgres
- Frontend modules: tailwind, react-query, router, radix-ui, shadcn, vitest
- Backend modules: gorm, auth-jwt, go-test
- Shared modules: docker, git, readme, env.example, vercel

## Development

Frontend app:

```bash
cd apps/web
npm install
npm run dev
```

Backend app:

```bash
cd apps/api
go mod tidy
go run ./cmd/server
```

## Environment

Each app owns its own env example:

```bash
cp apps/web/.env.example apps/web/.env
cp apps/api/.env.example apps/api/.env
```
