package t

// DottyConfig represents the `dotty.toml` file
type DottyConfig struct {
	ConfigDir     string `toml:"configDir"`
	SystemDirSrc  string `toml:"systemDirSrc"`
	SystemDirDest string `toml:"systemDirDest"`
	UserDirSrc    string `toml:"userDirSrc"`
	UserDirDest   string `toml:"userDirDest"`
	LocalDirSrc   string `toml:"localDirSrc"`
}
