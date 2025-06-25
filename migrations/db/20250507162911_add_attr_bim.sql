-- +goose Up
-- +goose StatementBegin
insert into attrs(al, nm, so)
select q.al, q.nm, -2
from (
	select 'bim' al, 'Big user image' nm
) q
where not exists(
	select 1
	from attrs c
	where c.al = q.al
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
