/*
add unique index on product name
*/

alter table products add constraint unique_name unique (name);
