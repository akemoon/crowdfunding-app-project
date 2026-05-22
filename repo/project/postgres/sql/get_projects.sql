with scored as (
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
        boosted_until is not null and boosted_until > now() as is_boosted,
        cover_key,
        (
            -- freshness: -0.5 per day since start
            -0.5 * (current_date - started_at::date)
            -- popularity: progress toward goal, capped at 100%
            + least(current_amount::float / nullif(goal_amount, 0), 1.0) * 20
            -- urgency: small bonus for projects close to deadline
            + case
                when (started_at + duration_days * interval '1 day')::date - current_date < 7
                then 5
                else 0
              end
            -- boost: always applied
            + case
                when boosted_until is not null and boosted_until > now()
                then 10
                else 0
              end
        ) as score,
        status_id,
        finished_at
    from projects
    where status_id = $1
      and ($2::smallint is null or category_id = $2::smallint)
      and ($3::text is null or name ilike '%' || $3 || '%')
)
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
    is_boosted,
    cover_key
from scored
order by
    case when $6 = 'default' then score end desc nulls last,
    case when $6 = 'date' then
        coalesce(finished_at, started_at)
    end desc nulls last,
    case when $6 = 'date' then score end desc nulls last,
    id desc
limit $4 offset $5;
