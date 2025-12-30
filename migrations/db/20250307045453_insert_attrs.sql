-- +goose Up
-- +goose StatementBegin
insert into attrs(al, nm)
select q.al, q.nm
from (
	select 'th' al, 'Image for thumb' nm union
	select 'ln' al, 'Last name' nm union
	select 'fn' al, 'First name' nm union
	select 'mn' al, 'Middle name' nm union
	select 'wph' al, 'Work phone' nm union
	select 'mph' al, 'Mobile phone' nm
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
