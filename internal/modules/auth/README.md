# Auth Module

Owns user authentication flows: registration, login, refresh tokens, password reset, email verification, and session revocation.

Suggested package split when you implement it:

- `domain`: entities, value objects, domain errors.
- `application`: use cases and transaction orchestration.
- `ports`: repository, token, mailer, and password-hasher interfaces.
- `adapter/http`: request/response handlers.
- `adapter/persistence`: database implementations.

Do not put password hashing or JWT signing directly in handlers. Keep those behind interfaces so tests stay cheap and security choices remain isolated.
