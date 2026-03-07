select
    id,
    user_id,
    category_id,
    name,
    description,
    currency_id,
    goal_amount,
    current_amount,
    started_at,
    duration_days
from projects
where id = $1;
