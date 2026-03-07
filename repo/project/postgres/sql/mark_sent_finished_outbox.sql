update finished_projects_outbox
set status_id = 2
where project_id = $1 and status_id = 1;
