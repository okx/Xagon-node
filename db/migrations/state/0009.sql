-- +migrate Up
ALTER TABLE IF EXISTS state.fork_id DROP CONSTRAINT IF EXISTS fork_id_block_num_fkey;

-- +migrate Down
DELETE FROM state.fork_id f
 WHERE NOT EXISTS(SELECT 1 FROM state.block b WHERE b.block_num = f.block_num);

ALTER TABLE IF EXISTS state.fork_id ADD CONSTRAINT fork_id_block_num_fkey
FOREIGN KEY(block_num) REFERENCES state.block (block_num) ON DELETE CASCADE;
