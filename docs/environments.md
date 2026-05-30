# Environments

The supported `APP_ENV` values are:

- `local`: laptop or local container development.
- `development`: deployed development environment shared by the team.
- `staging`: production-like pre-release environment.
- `production`: live customer environment.

`local` is allowed to start without external infrastructure so you can iterate on isolated code. Every non-local environment requires:

- `DATABASE_URL`
- `JWT_ACCESS_SECRET`
- `JWT_REFRESH_SECRET`

## Environment Files

The repository includes example files only:

- `.env.local.example`
- `.env.development.example`
- `.env.staging.example`
- `.env.production.example`

Do not commit real `.env` files or secrets.

## PostgreSQL

PostgreSQL is the primary database. The runtime uses `database/sql` with the `pgx` driver.

Local default:

```text
postgres://petverse:petverse@localhost:5432/petverse?sslmode=disable
```

Development, staging, and production examples use `sslmode=require`. Keep that unless your actual deployment platform terminates TLS in a controlled private network and documents why database TLS is not needed.
