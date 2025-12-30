-- +goose Up
-- +goose StatementBegin
insert into roles(al, nm, prot)
select q.al, q.nm, q.prot
from (
	select '' al, 'empty' nm, 1 prot union
	select 'adm' al, 'Admin' nm, 1 prot union
	select 'admr' al, 'Roles admin' nm, 1 prot union
	select 'admu' al, 'Users admin' nm, 1 prot
) q
where not exists(
	select 1
	from roles qq
	where qq.al = q.al
)
order by al;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
