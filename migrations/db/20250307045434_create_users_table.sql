-- +goose Up
-- +goose StatementBegin
create sequence public.sq_users_id
    start with 1
    increment by 1
    no minvalue
    no maxvalue
    cache 1;

create table users (
	id integer not null default nextval('public.sq_users_id'::regclass),
    ac varchar(255) not null,
	em varchar(255) not null,
	pas varchar(255) not null,
    created_at timestamp not null default now(),
	confirmed smallint not null default 0,
	constraint uq_users_ac unique (ac),
	constraint uq_users_em unique (em),
	constraint pk_users primary key (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
