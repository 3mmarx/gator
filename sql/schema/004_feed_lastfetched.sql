-- +goose Up
ALTER TABLE feeds add column last_fetched_at TIMESTAMP;


-- +goose Down
ALTER TABLE feeds drop column last_fetched_at;