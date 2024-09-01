create schema if not exists client;

create table if not exists client.client
(
    id integer PRIMARY KEY GENERATED ALWAYS AS IDENTITY unique,
    name varchar(128),
    settlement varchar(10) not null, -- в какой валюте храним деньги пользователя
    margin_algorithm int2 not null, -- алгоритм расчета маржи -- сделать кастомный тип для дениса
    gateway boolean not null,
    vendor boolean not null,
    is_active boolean default true,
    is_pro boolean, -- deprecated убрать нахуй!!!!!
    is_interbank boolean, -- только для payIN
    create_at timestamp default now(),
    update_at timestamp default '1970-01-01 00:00:00'
);

create table if not exists client.etl
(
    id integer PRIMARY KEY GENERATED ALWAYS AS IDENTITY unique,
    last_insert_id integer,
    last_update_at timestamp,
    update_at timestamp default now()
);

insert into client.etl
(
    last_insert_id,
    last_update_at
) values
(
    0,
    '1970-01-01 00:00:00'
)