-- +goose Up
-- +goose StatementBegin
create sequence sq_user_attrs_id
    start with 1
    increment by 1
    no minvalue
    no maxvalue
    cache 1;

CREATE TABLE user_attrs (
    id integer not null default nextval('sq_user_attrs_id'::regclass),
	user_id integer NOT NULL,
    attr varchar(50) NOT NULL,
	val varchar(255) NOT NULL,
	constraint pk_user_attrs primary key (id),
    constraint uq_user_attrs_usat unique (user_id, attr),
	constraint fk_user_attrs_users foreign key (user_id) references users(id) on delete cascade,
	constraint fk_user_attrs_attrs foreign key (attr) references attrs(al) on delete cascade
);

create index ix_user_attrs_user_id on user_attrs using btree (user_id);
create index ix_user_attrs_attr on user_attrs using btree (attr);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
