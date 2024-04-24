-- +migrate Up

-- This migration will delete all empty blocks
DELETE FROM state.block
WHERE NOT EXISTS (SELECT *
              FROM state.virtual_batch 
              WHERE state.virtual_batch.block_num = state.block.block_num) 
      AND NOT EXISTS (SELECT *
              FROM state.verified_batch  
              WHERE state.verified_batch.block_num = state.block.block_num)
      AND NOT EXISTS (SELECT *
              FROM state.forced_batch 
              WHERE state.forced_batch.block_num = state.block.block_num)
      AND NOT EXISTS (SELECT *
              FROM state.exit_root 
              WHERE state.exit_root.block_num = state.block.block_num)
      AND NOT EXISTS (SELECT *
              FROM state.monitored_txs 
              WHERE state.monitored_txs.block_num = state.block.block_num)
      AND NOT EXISTS (SELECT *
              FROM state.fork_id 
              WHERE state.fork_id.block_num = state.block.block_num);
              


-- +migrate Down

-- no action is needed, the data must remain deleted as it is useless