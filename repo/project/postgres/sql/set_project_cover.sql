update projects
set cover_key = $2
where id = $1 and status_id = 4;
