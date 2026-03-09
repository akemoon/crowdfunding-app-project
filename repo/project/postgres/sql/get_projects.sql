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
        (
            -(current_date - (
                case
                    when status_id = 2 then coalesce(finished_at, started_at + duration_days * interval '1 day')::date
                    else started_at::date
                end
            ))
            + case
                when status_id = 1
                 and boosted_until is not null
                 and boosted_until > now()
                then 10  -- boost score constant
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
    is_boosted
from scored
order by
    case when $6 = 'default' then score end desc nulls last,
    case when $6 = 'date' then
        case when status_id = 2 then finished_at else started_at end
    end desc nulls last,
    case when $6 = 'date' then score end desc nulls last,
    id desc
limit $4 offset $5;
