-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS statuses
(
    id serial PRIMARY KEY,
    name varchar
);

CREATE TABLE IF NOT EXISTS pull_requests
(
    id varchar PRIMARY KEY,
    name varchar,
    author_id varchar REFERENCES users(id),
    status_id INTEGER REFERENCES statuses(id),
    create_at timestamp,
    merged_at timestamp
);

CREATE TABLE IF NOT EXISTS reviewers
(
    pull_request_id varchar,
    reviewer_id varchar REFERENCES users(id)
);

INSERT INTO statuses (name)
values ('OPEN'),
       ('MERGED');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS reviewers, pull_requests, statuses;
-- +goose StatementEnd
