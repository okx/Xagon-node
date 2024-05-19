-- +migrate Up
CREATE TABLE IF NOT EXISTS pool.readytx(
                             id SERIAL PRIMARY KEY NOT NULL,
                             count INT,
                             updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
INSERT INTO pool.readytx(id, count) VALUES(1, 0)
    ON CONFLICT(id) do UPDATE
    SET count = 0;

-- +migrate Down
DROP TABLE IF EXISTS pool.readytx;
