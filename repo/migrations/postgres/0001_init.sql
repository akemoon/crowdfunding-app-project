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

insert into project_statuses (id, name) values
(1, 'review'),
(2, 'active'),
(3, 'finished');

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
    started_at     timestamptz,
    duration_days  smallint    not null,
    status_id      smallint    not null default 1 references project_statuses(id),
    boosted_until  timestamptz,
    finished_at    timestamptz,

    constraint projects_user_id_name_unique unique (user_id, name)
);

create table if not exists project_application_statuses
(
    id   smallint primary key,
    name text not null unique
);

insert into project_application_statuses (id, name) values
(1, 'pending'),
(2, 'rejected'),
(3, 'approved');

create table if not exists project_applications
(
    project_id    uuid        primary key references projects(id),
    status_id     smallint    not null default 1 references project_application_statuses(id),
    reject_reason text        not null default '',
    created_at    timestamptz not null default now()
);

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
drop table if exists project_applications;
drop table if exists project_application_statuses;
drop table if exists projects;
drop table if exists project_statuses;
drop table if exists project_categories;
drop table if exists currencies;
