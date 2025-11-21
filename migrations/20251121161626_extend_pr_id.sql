-- +goose Up
-- +goose StatementBegin
ALTER TABLE pull_requests
ALTER COLUMN id TYPE VARCHAR(10);

ALTER TABLE pull_requests_users
ALTER COLUMN pull_request_id TYPE VARCHAR(10);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE pull_requests
ALTER COLUMN id TYPE VARCHAR(5);

ALTER TABLE pull_requests_users
ALTER COLUMN pull_request_id TYPE VARCHAR(5);
-- +goose StatementEnd
