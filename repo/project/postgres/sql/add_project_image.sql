insert into project_images (project_id, storage_key)
values ($1, $2)
returning id
