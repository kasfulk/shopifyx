create table if not exists payments(
    id uuid not null default uuid_generate_v4() primary key,
    product_id uuid not null,
    product_name varchar not null,
    product_image_url varchar not null,
    product_qty int not null check(product_qty > 0),
    product_price int not null check(product_price >= 0),
    user_id uuid not null,
    buyer_username varchar not null,
    buyer_name varchar not null,
    bank_name varchar not null,
    bank_account_name varchar not null,
    bank_account_number varchar not null,
    created_at timestamptz not null default current_timestamp,
    updated_at timestamptz not null default current_timestamp
)