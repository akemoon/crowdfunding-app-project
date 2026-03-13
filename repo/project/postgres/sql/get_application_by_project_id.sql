select
    pa.status_id,
    pa.reject_reason,
    pa.created_at
from project_applications pa
join projects p on p.id = pa.project_id
where pa.project_id = $1
  and p.user_id = $2;
