-- +goose Up
-- +goose StatementBegin
insert into users(ac, em, pas, confirmed, prot)
select q.ac, q.em, q.pas, q.confirmed, q.prot
from (
	select 'admin' ac, 'admin' em, '111' pas, 1 confirmed, 1 prot
) q
where not exists(
	select 1
	from users qq
	where qq.ac = q.ac
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
