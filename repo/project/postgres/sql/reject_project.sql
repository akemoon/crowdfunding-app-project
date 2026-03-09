update project_applications
set status_id     = 2,
    reject_reason = $2
where project_id = $1
  and status_id = 1;
