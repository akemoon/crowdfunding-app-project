update projects
set category_id   = $3,
    name          = $4,
    description   = $5,
    currency_id   = $6,
    goal_amount   = $7,
    duration_days = $8
where id        = $1
  and user_id   = $2
  and status_id = 4;
