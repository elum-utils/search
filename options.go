package search

type Config struct {
	Interests []string
	LocalFile string
	Reset     bool
}

var defaultValue = Config{
	Interests: []string{},
	LocalFile: "file::memory:",
	Reset:     false,
}

func configDefault(config ...Config) Config {

	if len(config) < 1 {
		return defaultValue
	}

	cfg := config[0]

	if cfg.LocalFile == "" {
		cfg.LocalFile = defaultValue.LocalFile
	}

	return cfg

}
