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
    status_id,
    boosted_until is not null and boosted_until > now() as is_boosted,
    boosted_until,
    cover_key
from projects
where id = $1;
