-- +goose Up
-- +goose StatementBegin
alter table roles add prot smallint default 0 not null;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
