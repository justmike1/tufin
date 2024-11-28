package config

// CmdOptions represents the available command options as an enum
type CmdOptions int

const (
	CLUSTER CmdOptions = iota
	DEPLOY
	STATUS
)

func (c CmdOptions) String() string {
	switch c {
	case CLUSTER:
		return "cluster"
	case DEPLOY:
		return "deploy"
	case STATUS:
		return "status"
	default:
		return "unknown"
	}
}

func ParseCmdOption(option string) CmdOptions {
	switch option {
	case "cluster":
		return CLUSTER
	case "deploy":
		return DEPLOY
	case "status":
		return STATUS
	default:
		return -1 // or some other sentinel value
	}
}
