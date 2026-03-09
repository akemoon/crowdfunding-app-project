with updated as (
    update projects p
    set status_id = 3,
        finished_at = now()
    where p.status_id = 2
      and p.started_at + (p.duration_days * interval '1 day') <= now()
    returning
        p.id
)
insert into finished_projects_outbox (project_id, status_id)
select
    u.id,
    1
from updated u;
