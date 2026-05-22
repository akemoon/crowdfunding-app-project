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
    cover_key
from projects
where user_id = $1
  and ($2 or status_id in (2, 3))
order by started_at desc nulls last, id desc;
