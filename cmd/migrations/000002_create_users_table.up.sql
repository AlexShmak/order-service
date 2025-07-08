create table orders_service.users
(
    id            bigserial PRIMARY KEY,
    name          varchar(255) not null,
    password_hash bytea        not null,
    email         varchar(255) not null,
    created_at    timestamptz  not null default now()
)