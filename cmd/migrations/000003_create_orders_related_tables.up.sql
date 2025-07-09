create table if not exists orders_service.deliveries
(
    id         bigserial PRIMARY KEY,
    name       varchar(255) not null,
    phone      varchar(255) not null,
    zip        varchar(255) not null,
    city       varchar(255) not null,
    address    varchar(255) not null,
    region     varchar(255) not null,
    email      varchar(255) not null,
    created_at timestamptz  not null default now()
);

create table if not exists orders_service.payments
(
    id            bigserial PRIMARY KEY,
    transaction   varchar(255) not null,
    request_id    varchar(255), -- Can be empty
    currency      varchar(255) not null,
    provider      varchar(255) not null,
    amount        int          not null,
    payment_dt    bigint       not null,
    bank          varchar(255) not null,
    delivery_cost int          not null,
    goods_total   int          not null,
    custom_fee    int          not null,
    created_at    timestamptz  not null default now()
);

create table if not exists orders_service.orders
(
    id                 bigserial PRIMARY KEY,
    order_uid          uuid         not null unique,
    track_number       varchar(255) not null unique,
    entry              varchar(255) not null,
    delivery_data_id   bigint       not null,
    payment_data_id    bigint       not null,
    locale             varchar(255) not null,
    internal_signature varchar(255),
    customer_id        varchar(255) not null,
    delivery_service   varchar(255) not null,
    shardkey           varchar(255) not null,
    sm_id              int          not null,
    date_created       timestamptz  not null default now(),
    oof_shard          varchar(255) not null,

    foreign key (delivery_data_id) references orders_service.deliveries (id) on delete cascade,
    foreign key (payment_data_id) references orders_service.payments (id) on delete cascade
);

create table if not exists orders_service.items
(
    id           bigserial PRIMARY KEY,
    chrt_id      int          not null,
    track_number varchar(255) not null,
    price        int          not null,
    rid          varchar(255) not null,
    name         varchar(255) not null,
    sale         int          not null,
    size         varchar(255) not null,
    total_price  int          not null,
    nm_id        int          not null,
    brand        varchar(255) not null,
    status       int          not null,
    created_at   timestamptz  not null default now(),
    order_id     bigint       not null,

    foreign key (order_id) references orders_service.orders (id) on delete cascade
);