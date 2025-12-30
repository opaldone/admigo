-- +goose Up
-- +goose StatementBegin
create table attrs (
    al varchar(50) not null,
    nm varchar(255) not null,
    constraint pk_attrs primary key (al)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
