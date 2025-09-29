package config

type BuildOptimizationType int

const (
	BUILD_OPT_RELEASE BuildOptimizationType = iota
	BUILD_OPT_DEBUG
)

func (bt BuildOptimizationType) String() string {
	switch bt {
	case BUILD_OPT_RELEASE:
		return "release"
	case BUILD_OPT_DEBUG:
		return "debug"
	}
	return "unknown"
}
