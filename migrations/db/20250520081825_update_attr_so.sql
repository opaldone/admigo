-- +goose Up
-- +goose StatementBegin
update attrs qu set
so = (
	case
		when qu.al = 'bim' then -2
		when qu.al = 'th' then -1
		when qu.al = 'ln' then 1
		when qu.al = 'fn' then 2
		when qu.al = 'mn' then 3
		when qu.al = 'mph' then 4
		when qu.al = 'wph' then 5
	end
)
where (qu.al = 'bim' and qu.so != -2)
or (qu.al = 'th' and qu.so != -1)
or (qu.al = 'ln' and qu.so != 1)
or (qu.al = 'fn' and qu.so != 2)
or (qu.al = 'mn' and qu.so != 3)
or (qu.al = 'mph' and qu.so != 4)
or (qu.al = 'wph' and qu.so != 5);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
