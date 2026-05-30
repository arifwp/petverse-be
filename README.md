# Petverse Backend

Go backend foundation for Petverse.

## Current Scope

This repository is intentionally a foundation, not a feature implementation. It includes:

- Go module and API binary entrypoint.
- Environment configuration.
- PostgreSQL connection foundation with connection-pool settings and readiness checks.
- HTTP server with health and readiness endpoints.
- Structured logging.
- Graceful shutdown.
- Modular folder structure for auth, authorization, media upload, communities, discussions, and chat.
- Development docs, OpenAPI placeholder, Docker assets, and migration/schema folders.

## Run

```sh
make run
```

Then check:

```sh
curl http://localhost:8080/healthz
curl http://localhost:8080/readyz
```

For local PostgreSQL, Redis, and MinIO:

```sh
docker compose -f deployments/docker/docker-compose.yml up -d
```

## Structure

```text
cmd/api                 API executable entrypoint
internal/app            Application bootstrap and wiring
internal/config         Environment configuration
internal/platform       Runtime infrastructure and adapters
internal/modules        Product modules
internal/shared         Carefully limited shared code
api                     API contracts
docs                    Architecture and development notes
migrations              Database migrations
schema                  Schema drafts
deployments/docker      Docker development/deployment assets
test/integration        Integration test space
```

## Next Implementation Step

Implement one vertical slice first. Auth login or registration is the best candidate because it forces you to learn:

- HTTP request/response handling.
- Validation.
- Domain modeling.
- Password hashing.
- Persistence.
- Transaction boundaries.
- Tests.

Do not start by filling every module with interfaces. That looks organized but usually produces weak abstractions.
