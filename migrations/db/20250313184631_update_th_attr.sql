-- +goose Up
-- +goose StatementBegin
update attrs set
nm = 'Image for avatar'
where al = 'th'
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
