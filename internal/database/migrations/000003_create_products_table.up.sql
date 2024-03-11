create table if not exists products(
    id uuid not null default uuid_generate_v4() primary key,
    user_id uuid not null references users(id) on delete cascade,
    "name" varchar not null,
    price int not null default 0 check(price >= 0),
    image_url varchar not null,
    stock int not null default 0 check(stock >= 0),
    condition varchar not null,
    tags varchar[] not null default array[]::varchar[],
    is_purchaseable boolean not null,
    created_at timestamptz not null default current_timestamp,
    updated_at timestamptz not null default current_timestamp
)