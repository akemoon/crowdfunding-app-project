with updated_app as (
    update project_applications
    set status_id = 3
    where project_id = $1
      and status_id = 1
    returning project_id
)
update projects
set status_id  = 2,
    started_at = now()
where id = (select project_id from updated_app);
