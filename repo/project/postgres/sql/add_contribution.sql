update projects
set current_amount = current_amount + $2
where id = $1
  and status_id = 1;
