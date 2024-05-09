-- +migrate Up
ALTER TABLE state.block
    ADD COLUMN IF NOT EXISTS checked BOOL NOT NULL DEFAULT FALSE;

-- set block.checked to true for all blocks below max - 100
UPDATE state.block SET checked = true WHERE block_num <= (SELECT MAX(block_num) - 1000 FROM state.block);

-- +migrate Down
ALTER TABLE state.block
    DROP COLUMN IF EXISTS checked;

