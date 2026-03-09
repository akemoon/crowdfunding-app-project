with reset_app as (
    update project_applications
    set status_id     = 1,
        reject_reason = ''
    where project_id = $1
      and status_id  = 2
    returning project_id
)
update projects
set category_id   = $3,
    name          = $4,
    description   = $5,
    currency_id   = $6,
    goal_amount   = $7,
    duration_days = $8
where id      = (select project_id from reset_app)
  and user_id = $2
  and status_id = 1;
