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
    -- TODO: cover_url text (MinIO presigned upload, minio-go)

    constraint projects_user_id_name_unique unique (user_id, name)
);

create table if not exists project_application_statuses
(
    id   smallint primary key,
    name text not null unique
);

insert into project_application_statuses (id, name) values
(1, 'pending'),
(2, 'review'),
(3, 'rejected'),
(4, 'approved');

-- project_applications is a separate entity (not just a status on projects) because:
--   1. it carries manager-specific fields (assigned_to, reject_reason, timestamps)
--   2. in the future it should become 1:many — one project may have multiple applications
--      (e.g. re-submission after rejection), each with its own lifecycle and history
-- TODO: migrate to 1:many by adding id bigserial PK and removing project_id as PK
--       (keep project_id as FK + index); manager actions will then reference application id
create table if not exists project_applications
(
    project_id    uuid        primary key references projects(id),
    status_id     smallint    not null default 1 references project_application_statuses(id),
    assigned_to   uuid,
    reject_reason text        not null default '',
    created_at    timestamptz not null default now(),
    assigned_at   timestamptz,
    processed_at  timestamptz,
    updated_at    timestamptz not null default now()
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
