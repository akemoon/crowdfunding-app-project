select id, url, storage_key
from project_images
where id = $1 and project_id = $2
