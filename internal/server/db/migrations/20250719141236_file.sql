-- +goose Up
-- +goose StatementBegin
CREATE TYPE availability AS ENUM ('AVAILABLE', 'PROCESSING', 'RESERVED');
CREATE TABLE IF NOT EXISTS files (
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY
    file_name TEXT
    file_status availability
    owner_id BIGINT REFERENCES users(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS files;
-- +goose StatementEnd
