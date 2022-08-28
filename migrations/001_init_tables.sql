-- Write your migrate up statements here
create table links(
    id serial primary key,
    owner_id int,
    name varchar(155) not null unique,
    link text not null,
    visits int default 0,
    created_at timestamp not null default now(),
    updated_at timestamp not null default now(),
    expired_at timestamp,
    deleted_at timestamp
);
---- create above / drop below ----
drop table links;
-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
