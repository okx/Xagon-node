-- +migrate Up

-- the update below fix the wrong receipt TX indexes
WITH map_fix_tx_index AS (
    SELECT t.l2_block_num AS block_num
         , t.hash AS tx_hash
         , r.tx_index AS current_index
         , (ROW_NUMBER() OVER (PARTITION BY t.l2_block_num ORDER BY r.tx_index))-1 AS correct_index
      FROM state.receipt r
     INNER JOIN state."transaction" t
        ON t.hash = r.tx_hash
)
UPDATE state.receipt AS r
   SET tx_index = m.correct_index
  FROM map_fix_tx_index m
 WHERE m.block_num = r.block_num
   AND m.tx_hash = r.tx_hash
   AND m.current_index = r.tx_index
   AND m.current_index != m.correct_index;


-- +migrate Down

-- no action is needed, the data fixed by the 
-- migrate up must remain fixed
