update projects
set current_amount = current_amount + $2
where id = $1
returning current_amount;
