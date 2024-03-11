CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

create table if not exists users(
    id uuid not null default uuid_generate_v4() primary key,
    username varchar not null unique,
    "name" varchar not null,
    "password" varchar not null,
    created_at timestamptz not null default current_timestamp,
    updated_at timestamptz not null default current_timestamp
)