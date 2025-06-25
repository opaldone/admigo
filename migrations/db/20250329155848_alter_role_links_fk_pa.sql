-- +goose Up
-- +goose StatementBegin
ALTER TABLE role_links DROP CONSTRAINT fk_role_links_roles_pa;
ALTER TABLE role_links ADD CONSTRAINT fk_role_links_roles_pa FOREIGN KEY (parent) REFERENCES roles(al) ON DELETE SET DEFAULT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
