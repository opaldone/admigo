-- +goose Up
-- +goose StatementBegin
insert into role_links(child, parent)
select q.child, q.parent
from (
	select 'adm' child, '' parent union
	select 'admr' child, 'adm' parent union
	select 'admu' child, 'adm' parent
) q
where not exists(
	select 1
	from role_links qq
	where qq.child = q.child
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
