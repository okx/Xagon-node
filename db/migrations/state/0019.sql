-- +migrate Up
CREATE TABLE state.blob_inner
(
    blob_inner_num BIGINT PRIMARY KEY,
    data           BYTEA,
    block_num      BIGINT NOT NULL REFERENCES state.block (block_num) ON DELETE CASCADE    
);

ALTER TABLE state.virtual_batch
    ADD COLUMN IF NOT EXISTS blob_inner_num BIGINT, -- REFERENCES state.blob_inner (blob_inner_num),
    ADD COLUMN IF NOT EXISTS prev_l1_it_root VARCHAR,
    ADD COLUMN IF NOT EXISTS prev_l1_it_index BIGINT;

ALTER TABLE IF EXISTS state.proof RENAME TO batch_proof;

ALTER TABLE state.batch_proof 
    ADD COLUMN IF NOT EXISTS blob_inner_num  BIGINT; -- NOT NULL REFERENCES state.blob_inner (blob_inner_num) ON DELETE CASCADE;

CREATE TABLE state.blob_inner_proof
(
    blob_inner_num      BIGINT NOT NULL, -- REFERENCES state.blob_inner (blob_inner_num) ON DELETE CASCADE,
    proof_id            VARCHAR,
    proof               VARCHAR,
    input_prover        VARCHAR,
    prover              VARCHAR,
    prover_id           VARCHAR,  
    created_at          TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    generating_since    TIMESTAMP WITH TIME ZONE,
    PRIMARY KEY (blob_inner_num)    
);

CREATE TABLE state.blob_outer_proof
(
    blob_outer_num          BIGINT NOT NULL, -- REFERENCES state.blob_inner (blob_inner_num) ON DELETE CASCADE,
    blob_outer_num_final    BIGINT NOT NULL, -- REFERENCES state.blob_inner (blob_inner_num) ON DELETE CASCADE,
    proof_id                VARCHAR,
    proof                   VARCHAR,
    input_prover            VARCHAR,
    prover                  VARCHAR,
    prover_id               VARCHAR,  
    created_at              TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at              TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    generating_since        TIMESTAMP WITH TIME ZONE,
    PRIMARY KEY (blob_outer_num, blob_outer_num_final)    
);

-- +migrate Down
ALTER TABLE state.virtual_batch
    DROP COLUMN IF EXISTS blob_inner_num,
    DROP COLUMN IF EXISTS prev_l1_it_root,
    DROP COLUMN IF EXISTS prev_l1_it_index;

DROP TABLE state.blob_outer_proof;

DROP TABLE state.blob_inner_proof;

DROP TABLE state.batch_proof;

DROP TABLE state.blob_inner;

CREATE TABLE state.proof
(
    batch_num           BIGINT NOT NULL REFERENCES state.batch (batch_num) ON DELETE CASCADE,
    batch_num_final     BIGINT NOT NULL REFERENCES state.batch (batch_num) ON DELETE CASCADE,
    proof_id            VARCHAR,
    proof               VARCHAR,
    input_prover        VARCHAR,
    prover              VARCHAR,
    prover_id           VARCHAR,  
    created_at          TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    generating_since    TIMESTAMP WITH TIME ZONE,
    PRIMARY KEY (batch_num, batch_num_final)
);

ALTER TABLE state.virtual_batch
    DROP COLUMN IF EXISTS blob_inner_num,
    DROP COLUMN IF EXISTS prev_l1_it_root,
    DROP COLUMN IF EXISTS prev_l1_it_index;
    