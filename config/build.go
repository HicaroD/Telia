package config

type BuildType int

const (
	RELEASE BuildType = iota
	DEBUG
)

func (bt BuildType) String() string {
	switch bt {
	case RELEASE:
		return "release"
	case DEBUG:
		return "debug"
	}
	return "unknown"
}
