CREATE DATABASE IF NOT EXISTS minecraft_store;
USE minecraft_store;
create table if not exists products
(
    id          int auto_increment,
    name        varchar(128) null,
    description varchar(500) null,
    price       smallint     null,
    command     varchar(128) null,
    primary key (id)
);

create table if not exists purchases
(
    id             int auto_increment,
    product_id     int                   not null,
    nickname       varchar(64)           not null,
    email          varchar(128)          null,
    transaction_id boolean default false null,
    primary key (id)
);

create table if not exists transactions
(
    id     varchar(36),
    status varchar(20),
    primary key (id)
);




