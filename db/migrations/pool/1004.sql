-- +migrate Up
CREATE TABLE IF NOT EXISTS pool.free_gas (
    addr varchar NOT NULL PRIMARY KEY
);

-- +migrate Down
DROP TABLE IF EXISTS pool.free_gas;