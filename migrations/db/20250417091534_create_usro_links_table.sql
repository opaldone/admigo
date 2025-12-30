-- +goose Up
-- +goose StatementBegin
create sequence sq_usro_links_id
    start with 1
    increment by 1
    no minvalue
    no maxvalue
    cache 1;

CREATE TABLE usro_links (
    id integer not null default nextval('sq_usro_links_id'::regclass),
	user_id integer NOT NULL,
    role_al varchar(50) NOT NULL,
	constraint pk_usro_links primary key (id),
    constraint uq_usro_links_usro unique (user_id, role_al),
	constraint fk_usro_links_users foreign key (user_id) references users(id) on delete cascade,
	constraint fk_usro_links_roles foreign key (role_al) references roles(al) on delete cascade
);

create index ix_usro_links_user_id on usro_links using btree (user_id);
create index ix_usro_links_role_al on usro_links using btree (role_al);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
