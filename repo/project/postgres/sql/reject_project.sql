with rejected as (
    update project_applications
    set status_id     = 3,
        reject_reason = $2,
        processed_at  = now(),
        updated_at    = now()
    where project_id = $1
      and status_id  = 2
    returning project_id
)
update projects
set status_id = 4
from rejected
where projects.id = rejected.project_id
  and projects.status_id = 1;
