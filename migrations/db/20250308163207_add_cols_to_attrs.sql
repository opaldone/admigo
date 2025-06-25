-- +goose Up
-- +goose StatementBegin
alter table attrs add tp varchar(50) default 'users' not null;
alter table attrs add so smallint default 1 not null;

create index ix_attrs_tp on attrs using btree (tp);
create index ix_attrs_so on attrs using btree (so);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
