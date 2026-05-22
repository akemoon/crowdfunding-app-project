with submitted as (
    update projects
    set status_id = 1
    where id = $1 and user_id = $2 and status_id = 4
    returning id
)
insert into project_applications (project_id)
select id from submitted
on conflict (project_id) do update
set status_id     = 1,
    reject_reason = '',
    assigned_to   = null,
    assigned_at   = null,
    processed_at  = null,
    updated_at    = now();
