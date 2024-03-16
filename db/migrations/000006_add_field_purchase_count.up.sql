/*
add new field: purchase_count on products table
*/

alter table products add column purchase_count int not null default 0;

