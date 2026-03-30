insert into projects (
    user_id,
    category_id,
    name,
    description,
    currency_id,
    goal_amount,
    duration_days
) values ($1, $2, $3, $4, $5, $6, $7)
returning
    id,
    user_id,
    category_id,
    name,
    description,
    currency_id,
    goal_amount,
    started_at,
    duration_days,
    status_id;
