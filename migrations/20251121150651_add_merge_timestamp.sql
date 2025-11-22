-- +goose Up
-- +goose StatementBegin
ALTER TABLE pull_requests
ADD COLUMN merged_at TIMESTAMP WITH TIME ZONE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE pull_requests
DROP COLUMN merged_at;
-- +goose StatementEnd
