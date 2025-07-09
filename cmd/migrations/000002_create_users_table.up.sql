create table if not exists orders_service.users
(
    id         bigserial PRIMARY KEY,
    name       varchar(255) not null,
    password   bytea        not null,
    email      varchar(255) not null,
    created_at timestamptz  not null default now()
)