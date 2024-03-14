/*
drop unique index on product name if exist
*/

alter table products drop constraint if exists unique_name;

