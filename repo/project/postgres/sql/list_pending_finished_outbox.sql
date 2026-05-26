with locked as (
    update finished_projects_outbox
    set locked_until = now() + interval '5 minutes'
    where id in (
        select o.id
        from finished_projects_outbox o
        where o.status_id = 1
          and (o.locked_until is null or o.locked_until < now())
        order by o.id
        limit $1
        for update skip locked
    )
    returning project_id
)
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
from locked l
join projects p on p.id = l.project_id;
