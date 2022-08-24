-- Write your migrate up statements here
create table users(
    id serial primary key,
    username varchar(155) not null unique,
    password varchar(255) not null,
    created_at timestamp not null default now(),
    update_at timestamp not null default now(),
    deleted_at timestamp
);

create table links(
    id serial primary key,
    owner_id int not null,
    name varchar(155) not null unique,
    link text not null,
    visits int default 0,
    created_at timestamp not null default now(),
    update_at timestamp not null default now(),
    expired_at timestamp,
    deleted_at timestamp,
    constraint fk_owner foreign key(owner_id) references users(id)
);

create table api_keys(
    id serial primary key,
    user_id int not null,
    key varchar(255) not null,
    created_at timestamp not null default now(),
    deleted_at timestamp,
    constraint fk_user foreign key(user_id) references users(id)
);

create table blacklist_apikeys(
    id serial primary key,
    key_id int not null,
    created_at timestamp not null default now(),
    updated_at timestamp not null default now(),
    deleted_at timestamp,
    constraint fk_api_key foreign key(key_id) references api_keys(id)

);
---- create above / drop below ----
drop table users;
drop table links;
drop table api_keys;
drop table blacklist_apikeys;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
