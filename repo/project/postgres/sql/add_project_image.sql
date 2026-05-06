insert into project_images (project_id, url, storage_key)
values ($1, $2, $3)
returning id
