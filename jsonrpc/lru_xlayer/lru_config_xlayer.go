package lru_xlayer

var (
	gConfigure Config
)

// Config is the configuration for the LRU cache
type Config struct {
	Enable bool
	Size   int
}

// SetConfig sets the configuration for the LRU cache
func SetConfig(c Config) {
	gConfigure = c
}

// GetConfig returns the configuration for the LRU cache
func GetConfig() Config {
	return gConfigure
}
