package eigenda

import "github.com/Layr-Labs/eigenda/core"

// MockBlobRequestSigner is a moock EigenLayer signer for PoC purposes.
// On production, there will be a need to register a BlobRequestSigner with EigenLayer.
type MockBlobRequestSigner struct{}

// SignBlobRequest is the mock implementation for the mock signer
func (s MockBlobRequestSigner) SignBlobRequest(header core.BlobAuthHeader) ([]byte, error) {
	return []byte{}, nil
}

// GetAccountID is the mock implementation for the mock signer
func (s MockBlobRequestSigner) GetAccountID() string {
	return ""
}
