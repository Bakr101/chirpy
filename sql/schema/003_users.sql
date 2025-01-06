-- +goose Up
ALTER TABLE users
ADD COLUMN hashed_passwords TEXT DEFAULT 'unset';

ALTER TABLE users
ALTER COLUMN hashed_passwords SET NOT NULL;

-- +goose Down
ALTER TABLE users
DROP COLUMN hashed_passwords;