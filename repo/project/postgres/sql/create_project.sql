insert into projects (
    user_id,
    category_id,
    name,
    description,
    currency_id,
    goal_amount,
    duration_days,
    status_id
) values ($1, $2, $3, $4, $5, $6, $7, 4)
returning id
