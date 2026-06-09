create table if not exists users (
    id uuid primary key,
    name text,
    username text not null unique,
    email text not null unique,
    avatar_url text,
    created_at timestamp with time zone not null default now(),
    updated_at timestamp with time zone not null default now(),
    last_login_at timestamp with time zone,
    email_verified boolean not null default false,
    password text not null
);

create index if not exists idx_users_email on users (email);
create index if not exists idx_users_username on users (username);
