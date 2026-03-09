select
    pa.status_id,
    pa.reject_reason,
    pa.created_at,
    p.id,
    p.user_id,
    p.category_id,
    p.name,
    p.description,
    p.currency_id,
    p.goal_amount,
    p.current_amount,
    p.duration_days,
    p.status_id,
    p.boosted_until is not null and p.boosted_until > now() as is_boosted
from project_applications pa
join projects p on p.id = pa.project_id
where pa.status_id = 1
order by pa.created_at asc;
