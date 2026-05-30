# pgAdmin 4 Setup on macOS Without Docker

This guide creates a local PostgreSQL database for Petverse and connects it to pgAdmin 4.

The target local credentials are:

```text
Host: localhost
Port: 5432
Database: petverse
Username: petverse
Password: petverse
```

That matches `.env.local.example`:

```text
postgres://petverse:petverse@localhost:5432/petverse?sslmode=disable
```

## 1. Install PostgreSQL Locally

Use Homebrew:

```sh
brew install postgresql@17
```

Start PostgreSQL:

```sh
brew services start postgresql@17
```

Check that PostgreSQL is running:

```sh
brew services list
```

If `psql` is not found, add PostgreSQL to your shell path:

```sh
echo 'export PATH="/opt/homebrew/opt/postgresql@17/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc
```

For Intel Macs, Homebrew may use `/usr/local` instead:

```sh
echo 'export PATH="/usr/local/opt/postgresql@17/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc
```

## 2. Create The Petverse User And Database

Open the default local PostgreSQL shell:

```sh
psql postgres
```

Inside `psql`, run:

```sql
create role petverse with login password 'petverse';
create database petverse owner petverse;
grant all privileges on database petverse to petverse;
```

Exit:

```sql
\q
```

If the role or database already exists, PostgreSQL will show an error. That is not automatically fatal; it usually means you already created it.

## 3. Verify From Terminal

Run:

```sh
psql "postgres://petverse:petverse@localhost:5432/petverse?sslmode=disable"
```

Inside `psql`, test:

```sql
select current_database(), current_user, now();
```

Expected:

```text
current_database: petverse
current_user: petverse
```

Exit:

```sql
\q
```

Do not skip this step. If terminal access fails, pgAdmin will fail too.

## 4. Register The Server In pgAdmin 4

Open pgAdmin 4.

In the left sidebar:

1. Right-click `Servers`.
2. Select `Register`.
3. Select `Server...`.

In the `General` tab:

```text
Name: Petverse Local
```

In the `Connection` tab:

```text
Host name/address: localhost
Port: 5432
Maintenance database: petverse
Username: petverse
Password: petverse
```

Enable `Save password` for local development if you do not want to type it every time.

Click `Save`.

## 5. Find The Database In pgAdmin

After connecting:

1. Expand `Servers`.
2. Expand `Petverse Local`.
3. Expand `Databases`.
4. Open `petverse`.

The database will be empty for now because this project does not have migrations yet. That is expected.

## 6. Test With pgAdmin Query Tool

Right-click the `petverse` database, then select `Query Tool`.

Run:

```sql
select current_database(), current_user, now();
```

Expected:

```text
current_database: petverse
current_user: petverse
```

## 7. Connect The Go App To PostgreSQL

Create your local environment file:

```sh
cp .env.local.example .env
```

This project reads environment variables directly. It does not automatically load `.env`.

For a quick local run:

```sh
export DATABASE_URL='postgres://petverse:petverse@localhost:5432/petverse?sslmode=disable'
make run
```

Then check readiness:

```sh
curl http://localhost:8080/readyz
```

When PostgreSQL is reachable, the response should include:

```json
{
  "database": "ok",
  "status": "ready"
}
```

## Troubleshooting

### psql: command not found

Your PostgreSQL binary path is missing from `PATH`.

Apple Silicon:

```sh
echo 'export PATH="/opt/homebrew/opt/postgresql@17/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc
```

Intel:

```sh
echo 'export PATH="/usr/local/opt/postgresql@17/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc
```

### Connection refused

PostgreSQL is not running.

```sh
brew services start postgresql@17
brew services list
```

### role "petverse" does not exist

You skipped role creation or created it in the wrong PostgreSQL instance.

Run:

```sh
psql postgres
```

Then:

```sql
create role petverse with login password 'petverse';
```

### database "petverse" does not exist

Create it:

```sh
psql postgres
```

Then:

```sql
create database petverse owner petverse;
```

### password authentication failed

Reset the local password:

```sh
psql postgres
```

Then:

```sql
alter role petverse with password 'petverse';
```

### Port 5432 already in use

You likely have another PostgreSQL installation running. Do not guess. Check:

```sh
lsof -i :5432
```

If you have multiple PostgreSQL installations, decide which one you want to keep. Running multiple local PostgreSQL servers without understanding ports and data directories is a good way to waste time.

## Important Discipline

Local manual database setup is fine for learning. For shared development, staging, and production, do not click databases into existence through pgAdmin. Use versioned migrations and repeatable infrastructure setup.
