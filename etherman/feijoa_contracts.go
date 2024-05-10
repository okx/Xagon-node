package etherman

import (
	"github.com/0xPolygonHermez/zkevm-node/etherman/smartcontracts/feijoapolygonzkevm"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

// FeijoaContracts represents the contracts of the Feijoa upgrade
type FeijoaContracts struct {
	FeijoaZKEVMAddress common.Address
	FeijoaZKEVM        *feijoapolygonzkevm.Feijoapolygonzkevm
	//FeijoaRollupManager         *feijoapolygonrollupmanager.Feijoapolygonrollupmanager
	//FeijoaGlobalExitRootManager *feijoapolygonzkevmglobalexitroot.Feijoapolygonzkevmglobalexitroot
}

// NewFeijoaContracts creates a new FeijoaContracts
func NewFeijoaContracts(ethClient bind.ContractBackend, l1Config L1Config) (*FeijoaContracts, error) {
	FeijoaZKEVMAddress := l1Config.ZkEVMAddr
	FeijoaZKEVM, err := feijoapolygonzkevm.NewFeijoapolygonzkevm(FeijoaZKEVMAddress, ethClient)
	if err != nil {
		log.Errorf("error creating FeijoaZKEVM client (addr: %s). Error: %w", FeijoaZKEVMAddress.String(), err)
		return nil, err
	}

	return &FeijoaContracts{
		FeijoaZKEVMAddress: FeijoaZKEVMAddress,
		FeijoaZKEVM:        FeijoaZKEVM,
	}, nil
}

// GetAddresses returns the addresses of the contracts
func (f *FeijoaContracts) GetAddresses() []common.Address {
	return []common.Address{
		f.FeijoaZKEVMAddress,
	}
}
