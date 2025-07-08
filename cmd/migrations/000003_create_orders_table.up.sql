create table if not exists orders_service.orders
(
    id          bigserial PRIMARY KEY,
    user_id     bigint         not null,
    product_id  bigint         not null,
    quantity    integer        not null,
    total_price numeric(10, 2) not null,
    created_at  timestamptz    not null default now(),
    updated_at  timestamptz    not null default now(),
    foreign key (user_id) references orders_service.users (id),
    foreign key (product_id) references orders_service.products (id)
);