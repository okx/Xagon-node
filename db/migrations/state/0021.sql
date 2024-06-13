-- +migrate Up
ALTER TABLE state.batch
    ADD COLUMN high_reserved_counters JSONB;

-- +migrate Down
ALTER TABLE state.batch
    DROP COLUMN high_reserved_counters;
    