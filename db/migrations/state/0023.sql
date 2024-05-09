-- +migrate Up

CREATE TABLE IF NOT EXISTS state.blob_sequence
(
    index                BIGINT PRIMARY KEY,
    coinbase             VARCHAR,
    final_acc_input_hash VARCHAR,
    first_blob_sequenced  BIGINT,
    last_blob_sequenced   BIGINT,
    created_at           TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    received_at          TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    block_num          BIGINT NOT NULL REFERENCES state.block (block_num) ON DELETE CASCADE 
);

comment on column state.blob_sequence.index is 'It is the id of this sequence, this value is internal and incremental';
comment on column state.blob_sequence.block_num is 'L1 Block where appear this sequence';
comment on column state.blob_sequence.first_blob_sequenced is 'first (included) blob_inner_num of this sequence (state.blob_inner.blob_inner_num)';
comment on column state.blob_sequence.first_blob_sequenced is 'last (included) blob_inner_num of this sequence (state.blob_inner.blob_inner_num)';
comment on column state.blob_sequence.received_at is 'time when it was received in node';
comment on column state.blob_sequence.created_at is 'time when was created on L1 (L1block tstamp)';

CREATE TABLE IF NOT EXISTS state.blob_inner_in 
(
    blob_inner_num      BIGINT PRIMARY KEY,
    blob_sequence_index BIGINT NOT NULL REFERENCES state.blob_sequence (index) ON DELETE CASCADE,
    blob_type         VARCHAR,
    max_sequence_timestamp TIMESTAMP WITH TIME ZONE,
    zk_gas_limit        BIGINT,
    l1_info_tree_leaf_index BIGINT,
    l1_info_tree_root VARCHAR,
    blob_data_hash VARCHAR,
    updated_at          TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    -- if blob_type== blob
    blob_type_index    BIGINT,
    blob_type_z       VARCHAR,
    blob_type_y      VARCHAR,
    blob_type_commitment VARCHAR,
    blob_type_proof VARCHAR
);

comment on column state.blob_inner_in.updated_at is 'the creation time is blob_sequence.created_at, this is the last time when was updated (tipically Now() )';
comment on column state.blob_inner_in.blob_type is 'call_data, blob or forced';
comment on column state.blob_inner_in.blob_data_hash is 'is the hash of the blobData';

CREATE TABLE IF NOT EXISTS state.incoming_batch
(
    batch_num                BIGINT PRIMARY KEY,
    blob_inner_num          BIGINT NOT NULL REFERENCES state.blob_inner_in (blob_inner_num) ON DELETE CASCADE,
    data           BYTEA,
    created_at           TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),   
    updated_at          TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
); 

-- +migrate Down
DROP TABLE IF EXISTS state.incoming_batch;
DROP TABLE IF EXISTS state.blob_inner_in;
DROP TABLE IF EXISTS state.blob_sequence;
