-- +goose Up
-- +goose StatementBegin
CREATE TYPE status AS ENUM ('OPEN', 'MERGED');

CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(5) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS teams (
    name VARCHAR(255) PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS pull_requests (
    id VARCHAR(5) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    author VARCHAR(5) NOT NULL,
    status status NOT NULL DEFAULT 'OPEN',
    FOREIGN KEY (author) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS teams_users (
    team_name VARCHAR(255) NOT NULL,
    user_id VARCHAR(5) NOT NULL,
    PRIMARY KEY (team_name, user_id),
    FOREIGN KEY (team_name) REFERENCES teams(name),
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS pull_requests_users (
    pull_request_id VARCHAR(5) NOT NULL,
    user_id VARCHAR(5) NOT NULL,
    PRIMARY KEY (pull_request_id, user_id),
    FOREIGN KEY (pull_request_id) REFERENCES pull_requests(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS pull_requests_users CASCADE;
DROP TABLE IF EXISTS teams_users CASCADE;
DROP TABLE IF EXISTS pull_requests CASCADE;
DROP TABLE IF EXISTS teams CASCADE;
DROP TABLE IF EXISTS users CASCADE;
DROP TYPE IF EXISTS status CASCADE;
-- +goose StatementEnd
