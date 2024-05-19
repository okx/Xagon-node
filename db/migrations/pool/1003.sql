-- +migrate Up
CREATE TABLE IF NOT EXISTS pool.innertx (
                              hash VARCHAR(128) PRIMARY KEY NOT NULL,
                              innertx text,
                              created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +migrate Down
DROP TABLE IF EXISTS pool.innertx;