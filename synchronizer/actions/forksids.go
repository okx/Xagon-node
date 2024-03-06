package actions

// ForkIdType is the type of the forkId
type ForkIdType uint64

const (
	// WildcardForkId It match for all forkIds
	WildcardForkId ForkIdType = 0
	// ForkIDIncaberry is the forkId for incaberry
	ForkIDIncaberry = ForkIdType(6) // nolint:gomnd
	// ForkIDEtrog is the forkId for etrog
	ForkIDEtrog = ForkIdType(7) //nolint:gomnd
	// ForkIDElderberry is the forkId for Elderberry
	ForkIDElderberry = ForkIdType(8) //nolint:gomnd
)

var (

	// ForksIdAll support all forkIds
	ForksIdAll = []ForkIdType{WildcardForkId}

	// ForksIdOnlyElderberry support only elderberry forkId
	ForksIdOnlyElderberry = []ForkIdType{ForkIDElderberry}

	// ForksIdOnlyEtrog support only etrog forkId
	ForksIdOnlyEtrog = []ForkIdType{ForkIDEtrog}

	// ForksIdToIncaberry support all forkIds till incaberry
	ForksIdToIncaberry = []ForkIdType{1, 2, 3, 4, 5, ForkIDIncaberry}

	// ForksIdEtrogAndIncaberry support only etrog and incaberry forkIds
	ForksIdEtrogAndIncaberry = []ForkIdType{ForkIDEtrog, ForkIDElderberry}
)
