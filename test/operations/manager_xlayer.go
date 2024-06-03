package operations

const (
	DefaultL1DataCommitteeContract        = "0x6Ae5b0863dBF3477335c0102DBF432aFf04ceb22"
	DefaultL1AdminAddress                 = "0x8f8E2d6cF621f30e9a11309D6A56A876281Fd534"
	DefaultL1AdminPrivateKey              = "0x815405dddb0e2a99b12af775fd2929e526704e1d1aea6a0b4e74dc33e2f7fcd2"
	MaxBatchesForL1                uint64 = 10
)

// StartDACDB starts the data availability node DB
func (m *Manager) StartDACDB() error {
	return StartComponent("dac-db", func() (bool, error) { return true, nil })
}

// StopDACDB stops the data availability node DB
func (m *Manager) StopDACDB() error {
	return StopComponent("dac-db")
}

// StartPermissionlessNodeForcedToSYncThroughDAC starts a permissionless node that is froced to sync through the DAC
func (m *Manager) StartPermissionlessNodeForcedToSYncThroughDAC() error {
	return StartComponent("permissionless-dac", func() (bool, error) { return true, nil })
}

// StopPermissionlessNodeForcedToSYncThroughDAC stops the permissionless node that is froced to sync through the DAC
func (m *Manager) StopPermissionlessNodeForcedToSYncThroughDAC() error {
	return StopComponent("permissionless-dac")
}
