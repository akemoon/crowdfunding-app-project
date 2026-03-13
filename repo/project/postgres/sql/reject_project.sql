update project_applications
set status_id    = 3,
    reject_reason = $2,
    processed_at = now(),
    updated_at   = now()
where project_id = $1
  and status_id = 2;
