-- +goose Up

create table if not exists currencies
(
    id   smallint primary key,
    name char(3) not null unique
);

insert into currencies (id, name) values
(1, 'RUB'),
(2, 'USD');

create table if not exists project_categories
(
    id   smallint primary key,
    name text not null unique
);

-- TODO: add more categories
insert into project_categories (id, name) values
(1, 'science'),
(2, 'tech'),
(3, 'architecture_and_urban'),
(4, 'sport'),
(5, 'music');

create table if not exists project_statuses
(
    id   smallint primary key,
    name text not null unique
);

-- TODO: add status "on review"
insert into project_statuses (id, name) values
(1, 'active'),
(2, 'finished');

create table if not exists projects
(
    id             uuid        primary key default uuidv7(),
    user_id        uuid        not null,
    category_id    smallint    not null references project_categories(id),
    name           text        not null,
    description    text        not null,
    currency_id    smallint    not null references currencies(id),
    goal_amount    bigint      not null check (goal_amount > 0),
    current_amount bigint      not null default 0 check (current_amount >= 0),
    started_at     timestamptz not null default now(),
    duration_days  smallint    not null,
    status_id      smallint    not null default 1 references project_statuses(id),  -- TODO: default status is "on review"

    constraint projects_user_id_name_unique unique (user_id, name)
);

insert into projects (
    user_id,
    category_id,
    name,
    description,
    currency_id,
    goal_amount,
    current_amount,
    started_at,
    duration_days
) values
('019bfb8f-d27b-7089-8296-0584baf0408e', 1, 'Лаборатория по астрофизике', 'Сбор средств на оборудование и материалы для студенческой лаборатории.', 1, 150000, 125213, now() - interval '10 days', 20),
('019bfb8f-f71b-7abc-abbd-f8cbfe2a38d4', 2, 'Робот-помощник', 'Сбор средств на создание бытовых роботов-помощников.', 2, 10000, 9243, now() - interval '120 days', 100),
('019bfb90-1e61-7016-ae5b-73d7c439cccd', 5, 'Музыкальный альбом "Северный ветер"', 'Запись дебютного альбома с живыми инструментами.', 1, 25000, 19320, now() - interval '5 days', 15);

create table if not exists finished_projects_outbox_statuses
(
    id   smallint primary key,
    name text not null unique
);

insert into finished_projects_outbox_statuses (id, name) values
(1, 'pending'),
(2, 'sent');

create table if not exists finished_projects_outbox
(
    id         bigserial   primary key,
    project_id uuid        not null references projects(id),
    status_id  smallint    not null default 1 references finished_projects_outbox_statuses(id),
    created_at timestamptz not null default now()
);

-- +goose Down

drop table if exists finished_projects_outbox;
drop table if exists finished_projects_outbox_statuses;
drop table if exists projects;
drop table if exists project_statuses;
drop table if exists project_categories;
drop table if exists currencies;
