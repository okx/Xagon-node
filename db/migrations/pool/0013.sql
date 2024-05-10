-- +migrate Up
ALTER TABLE pool.transaction
    ADD COLUMN reserved_zkcounters jsonb DEFAULT '{}'::jsonb;

UPDATE pool."transaction" set reserved_zkcounters  = jsonb_set(reserved_zkcounters , '{GasUsed}', cast(cumulative_gas_used  as text)::jsonb, true);
UPDATE pool."transaction" set reserved_zkcounters  = jsonb_set(reserved_zkcounters , '{KeccakHashes}', cast(used_keccak_hashes as text)::jsonb, true);
UPDATE pool."transaction" set reserved_zkcounters  = jsonb_set(reserved_zkcounters , '{PoseidonHashes}', cast(used_poseidon_hashes  as text)::jsonb, true);
UPDATE pool."transaction" set reserved_zkcounters  = jsonb_set(reserved_zkcounters , '{PoseidonPaddings}', cast(used_poseidon_paddings as text)::jsonb, true);
UPDATE pool."transaction" set reserved_zkcounters  = jsonb_set(reserved_zkcounters , '{MemAligns}', cast(used_mem_aligns as text)::jsonb, true);
UPDATE pool."transaction" set reserved_zkcounters  = jsonb_set(reserved_zkcounters , '{Arithmetics}', cast(used_arithmetics as text)::jsonb, true);
UPDATE pool."transaction" set reserved_zkcounters  = jsonb_set(reserved_zkcounters , '{Binaries}', cast(used_binaries as text)::jsonb, true);
UPDATE pool."transaction" set reserved_zkcounters  = jsonb_set(reserved_zkcounters , '{Steps}', cast(used_steps  as text)::jsonb, true);
UPDATE pool."transaction" set reserved_zkcounters  = jsonb_set(reserved_zkcounters , '{Sha256Hashes_V2}', cast(used_sha256_hashes as text)::jsonb, true);

-- +migrate Down
ALTER TABLE pool.transaction
    DROP COLUMN reserved_zkcounters;
