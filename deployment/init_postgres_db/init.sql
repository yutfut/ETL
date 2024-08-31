create schema client;

create table client.client
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
    create_at timestamp default now()
);
