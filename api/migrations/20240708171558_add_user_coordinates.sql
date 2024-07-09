-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS cube;
CREATE EXTENSION IF NOT EXISTS earthdistance;
ALTER TABLE users ADD COLUMN location_lat DOUBLE PRECISION, ADD COLUMN location_long DOUBLE PRECISION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN location_lat, DROP COLUMN location_long;
DROP EXTENSION IF EXISTS earthdistance;
DROP EXTENSION IF EXISTS cube;
-- +goose StatementEnd