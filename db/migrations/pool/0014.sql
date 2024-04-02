-- +migrate Up
CREATE INDEX IF NOT EXISTS idx_transaction_gas_price ON pool.transaction (gas_price);

-- +migrate Down
DROP INDEX IF EXISTS pool.idx_transaction_gas_price;
