select
    p.id,
    p.user_id,
    p.category_id,
    p.name,
    p.description,
    p.currency_id,
    p.goal_amount,
    p.current_amount,
    p.started_at,
    p.duration_days,
    p.status_id,
    p.boosted_until is not null and p.boosted_until > now() as is_boosted
from finished_projects_outbox o
join projects p on p.id = o.project_id
where o.status_id = 1
order by o.id
limit $1;
