# API Notes

Only system endpoints exist in the scaffold:

- `GET /healthz`: process health.
- `GET /readyz`: readiness check, including PostgreSQL status when configured.

Feature endpoints should be added only when their use cases exist. Avoid creating route files full of empty handlers because that gives a false sense of progress.
