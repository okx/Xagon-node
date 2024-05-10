-- +migrate Up
ALTER TABLE state.batch
    ADD COLUMN IF NOT EXISTS checked BOOLEAN NOT NULL DEFAULT TRUE;

-- +migrate Down
ALTER TABLE state.batch
    DROP COLUMN IF EXISTS checked;