-- +goose Up
-- +goose StatementBegin
create sequence sq_role_links_id
    start with 1
    increment by 1
    no minvalue
    no maxvalue
    cache 1;

create table role_links (
    id integer not null default nextval('sq_role_links_id'::regclass),
    child varchar(50) not null,
    parent varchar(50) not null default '',
    constraint pk_role_links primary key (id),
    constraint uq_role_links_chpa unique (child, parent),
	constraint fk_role_links_roles_ch foreign key (child) references roles(al) on delete cascade,
	constraint fk_role_links_roles_pa foreign key (parent) references roles(al) on delete set null
);

create index ix_role_links_ch on role_links using btree (child);
create index ix_role_links_pa on role_links using btree (parent);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
