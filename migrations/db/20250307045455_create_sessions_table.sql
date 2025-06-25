-- +goose Up
-- +goose StatementBegin
CREATE TABLE sessions (
    uid varchar(64) NOT NULL,
	user_id integer NOT NULL,
    created_at timestamp not null default now(),
	constraint pk_sessions primary key (uid),
	constraint fk_sessions_users foreign key (user_id) references users(id) on delete cascade
);

create index ix_sessions_user_id on sessions using btree (user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
