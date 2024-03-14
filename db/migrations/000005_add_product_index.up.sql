/*
add unique index on product name combine with user id
*/

alter table products add constraint unique_name unique (name, user_id);
