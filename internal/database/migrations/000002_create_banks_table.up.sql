create table if not exists banks(
    id bigserial primary key,
    user_id bigserial references users(id) on delete cascade,
    bank_name varchar not null,
    bank_account_number varchar not null,
    bank_account_name varchar not null,
    created_at timestamptz not null default current_timestamp,
    updated_at timestamptz not null default current_timestamp
)