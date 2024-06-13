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
	// ForkIDElderberry2 is the forkId for Elderberry2
	ForkIDElderberry2 = ForkIdType(9) //nolint:gomnd
	// ForkIDFeijoa is the forkId for Feijoa
	ForkIDFeijoa = ForkIdType(10) //nolint:gomnd
)

var (

	/// ************** ALL ***************///

	// ForksIdAll support all forkIds
	ForksIdAll = []ForkIdType{WildcardForkId}

	/// ************** SINGLE ***************///

	// ForksIdOnlyFeijoa support only etrog forkId
	ForksIdOnlyFeijoa = []ForkIdType{ForkIDFeijoa}

	// ForksIdOnlyElderberry support only elderberry forkId
	ForksIdOnlyElderberry = []ForkIdType{ForkIDElderberry, ForkIDElderberry2}

	// ForksIdOnlyEtrog support only etrog forkId
	ForksIdOnlyEtrog = []ForkIdType{ForkIDEtrog}

	/// ************** MULTIPLE ***************///

	// ForksIdToIncaberry support all forkIds till incaberry
	ForksIdToIncaberry = []ForkIdType{1, 2, 3, 4, 5, ForkIDIncaberry}

	// ForksIdToEtrog support all forkIds till etrog
	ForksIdToEtrog = append(ForksIdToIncaberry, ForksIdOnlyEtrog...)

	// ForksIdToElderberry support all forkIds till elderberry
	ForksIdToElderberry = append(ForksIdToEtrog, ForksIdOnlyElderberry...)

	// ForksIdToFeijoa support all forkIds till feijoa
	ForksIdToFeijoa = append(ForksIdToElderberry, ForksIdOnlyFeijoa...)
)
