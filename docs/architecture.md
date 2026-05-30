# Petverse Backend Architecture

## Positioning

This scaffold is a modular monolith. That is the right starting point for this project because splitting auth, media, community, discussion, and chat into microservices before the domain is stable would add operational cost without proving value.

## Boundaries

- `cmd/api`: executable entrypoint only.
- `internal/app`: application bootstrapping and dependency wiring.
- `internal/config`: environment-driven configuration.
- `internal/platform`: infrastructure adapters and shared runtime concerns.
- `internal/modules`: product modules with their own domain, use cases, ports, and adapters.
- `internal/shared`: only for concepts reused by multiple modules.
- `migrations`: database schema changes.
- `api`: API contracts such as OpenAPI.
- `schema`: database design drafts and diagrams.

## Dependency Rule

Product modules may depend inward on their own domain and ports. Infrastructure depends on module ports, not the other way around.

Handlers should be thin:

1. Decode and validate transport input.
2. Call an application use case.
3. Map the result to an HTTP response.

Use cases should own transaction boundaries, authorization checks, and orchestration. Repositories should only persist and fetch data.

## Initial Module Roadmap

1. Auth: users, credentials, sessions, refresh tokens, password reset, email verification.
2. Authorization: roles, permissions, ownership checks.
3. Media: multiple image upload, metadata, object storage, moderation hooks.
4. Communities: community creation, membership, moderation.
5. Discussions: questions, answers, comments, voting, accepted answers.
6. Chat: conversations, participants, messages, delivery/read state.

## What Not To Build Yet

- Microservices.
- Event sourcing.
- Generic repository abstractions.
- A custom dependency injection framework.
- Premature background worker topology.

Those choices may become valid later. Right now they would mostly slow learning and hide the fundamentals.
