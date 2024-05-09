-- +migrate Up
ALTER TABLE state.receipt
    ADD COLUMN IF NOT EXISTS im_state_root BYTEA;

UPDATE state.receipt SET im_state_root = post_state WHERE block_num >= (SELECT MIN(block_num) FROM state.l2block WHERE batch_num >= (SELECT from_batch_num FROM state.fork_id WHERE fork_id = 7));

-- +migrate Down
ALTER TABLE state.receipt
    DROP COLUMN IF EXISTS im_state_root;
