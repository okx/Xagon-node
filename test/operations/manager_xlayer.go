package operations

const (
	DefaultL1DataCommitteeContract        = "0x6Ae5b0863dBF3477335c0102DBF432aFf04ceb22"
	DefaultL1AdminAddress                 = "0x2ecf31ece36ccac2d3222a303b1409233ecbb225"
	DefaultL1AdminPrivateKey              = "0xde3ca643a52f5543e84ba984c4419ff40dbabd0e483c31c1d09fee8168d68e38"
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
