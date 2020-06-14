package routerscan

type StValueBool uint
type StValueInt uint
type StValueString uint
type StValuePointer uint

const (
	StEnableDebug          StValueBool    = 0
	StDebugVerbosity       StValueInt     = 1
	StWriteLogCallback     StValuePointer = 2
	StSetTableDataCallback StValuePointer = 3
	StUserAgent            StValueString  = 4
	StUseCustomPage        StValueBool    = 5
	StCustomPage           StValueString  = 6
	StDualAuthCheck        StValueBool    = 7
	StPairsBasic           StValueString  = 8
	StPairsDigest          StValueString  = 9
	StProxyType            StValueInt     = 10
	StProxyIp              StValueString  = 11
	StProxyPort            StValueInt     = 12
	StUseCredentials       StValueBool    = 13
	StCredentialsUsername  StValueString  = 14
	StCredentialsPassword  StValueString  = 15
	StPairsForm            StValueString  = 16
	StFilterRules          StValueString  = 17
	StProxyUseAuth         StValueBool    = 18
	StProxyUser            StValueString  = 19
	StProxyPass            StValueString  = 20
)
