package routerscan

type StOptions struct {
	EnableDebug         bool
	DebugVerbosity      int
	UserAgent           string
	UseCustomPage       bool
	CustomPage          string
	DualAuthCheck       bool
	PairsBasic          string
	PairsDigest         string
	ProxyType           int
	ProxyIp             string
	ProxyPort           int
	CredentialsUsername string
	CredentialsPassword string
	PairsForm           string
	FilterRules         string
	ProxyUseAuth        bool
	ProxyUser           string
	ProxyPass           string
}
