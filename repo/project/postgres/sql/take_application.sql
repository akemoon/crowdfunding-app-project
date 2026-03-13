update project_applications
set status_id   = 2,
    assigned_to = $2,
    assigned_at = now(),
    updated_at  = now()
where project_id = $1
  and status_id = 1;
