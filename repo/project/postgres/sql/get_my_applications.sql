select
    pa.status_id,
    pa.reject_reason,
    pa.created_at,
    pa.assigned_at,
    pa.processed_at,
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
where pa.assigned_to = $1
  and pa.status_id = 2
order by pa.assigned_at asc;
