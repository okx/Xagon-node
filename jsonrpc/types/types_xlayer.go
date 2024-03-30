package types

// BatchFilter is a list of batch numbers to retrieve
type BatchFilter struct {
	Numbers []BatchNumber `json:"numbers"`
}

// BatchData is an abbreviated structure that only contains the number and L2 batch data
type BatchData struct {
	Number      ArgUint64 `json:"number"`
	BatchL2Data ArgBytes  `json:"batchL2Data,omitempty"`
	Empty       bool      `json:"empty"`
}

// BatchDataResult is a list of BatchData for a BatchFilter
type BatchDataResult struct {
	Data []*BatchData `json:"data"`
}

// Contains checks if a string is contained in a slice of strings
func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
