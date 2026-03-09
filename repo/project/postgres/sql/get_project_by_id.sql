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
    duration_days,
    boosted_until is not null and boosted_until > now() as is_boosted
from projects
where id = $1;
