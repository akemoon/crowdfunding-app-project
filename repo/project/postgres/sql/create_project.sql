with created as (
    insert into projects (
        user_id,
        category_id,
        name,
        description,
        currency_id,
        goal_amount,
        duration_days
    ) values ($1, $2, $3, $4, $5, $6, $7)
    returning id
)
insert into project_applications (project_id)
select id from created;
