# Media Module

Owns image upload workflows, including multiple image upload, validation, metadata persistence, and object storage.

Implementation boundaries:

- Validate size, MIME type, and extension before storage.
- Store binary files outside the database.
- Persist metadata and ownership in the database.
- Use background jobs only when thumbnailing or moderation becomes expensive enough to need it.
