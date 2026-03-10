update projects
set boosted_until = greatest(coalesce(boosted_until, now()), now()) + ($3 * interval '1 day')
where id = $1
  and user_id = $2
  and status_id = 2;
