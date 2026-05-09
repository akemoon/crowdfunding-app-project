select id, storage_key
from project_images
where project_id = $1
order by created_at
