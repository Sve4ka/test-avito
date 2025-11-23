-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users(
    id varchar PRIMARY KEY,
    username varchar,
    team_name varchar,
    is_active boolean
);

CREATE INDEX IF NOT EXISTS users_team_name_idx
    ON users (team_name);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS users_team_name_idx;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
