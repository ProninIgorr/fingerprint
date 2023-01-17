package errs

type Severity uint32

const (
	SeverityWarning Severity = iota
	SeverityError
	SeverityCritical
)

var AllSeverities = [...]Severity{
	SeverityWarning,
	SeverityError,
	SeverityCritical,
}

func (s Severity) String() string {
	switch s {
	case SeverityWarning:
		return "wrn"
	case SeverityError:
		return "err"
	case SeverityCritical:
		return "cri"
	default:
		return "unknown error severity"
	}
}
