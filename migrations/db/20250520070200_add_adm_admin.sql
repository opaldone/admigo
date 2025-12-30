-- +goose Up
-- +goose StatementBegin
insert into usro_links(user_id, role_al)
select q.user_id, q.role_al
from (
	select
	(select u.id from users u where u.ac = 'admin') user_id,
	(select r.al from roles r where r.al = 'adm') role_al
) q
where not exists(
	select 1
	from usro_links qq
	where qq.user_id = q.user_id
	and qq.role_al = q.role_al
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
