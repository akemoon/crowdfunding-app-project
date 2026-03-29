select
    id,
    user_id,
    category_id,
    name,
    description,
    currency_id,
    goal_amount,
    started_at,
    duration_days,
    status_id
from projects
order by id desc;
