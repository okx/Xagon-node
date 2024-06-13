-- +migrate Up

-- +migrate Up
ALTER TABLE state.exit_root
    ADD COLUMN IF NOT EXISTS l1_info_tree_recursive_index BIGINT DEFAULT NULL UNIQUE;
CREATE INDEX IF NOT EXISTS idx_exit_root_l1_info_tree_recursive_index ON state.exit_root (l1_info_tree_recursive_index);

-- +migrate Down
ALTER TABLE state.exit_root
    DROP COLUMN IF EXISTS l1_info_tree_recursive_index;
DROP INDEX IF EXISTS state.idx_exit_root_l1_info_tree_recursive_index;

