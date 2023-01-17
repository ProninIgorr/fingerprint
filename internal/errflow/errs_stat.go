package errflow

import "github.com/ProninIgorr/fingerprint/internal/errs"

type ErrStatKey struct {
	Severity   errs.Severity
	Kind       errs.Kind
	Operations string
}

func (e ErrStatKey) Less(other interface{}) bool {
	ee := other.(ErrStatKey)
	if e.Severity == ee.Severity {
		if e.Operations == ee.Operations {
			return e.Kind < ee.Kind
		}
		return e.Operations < ee.Operations
	}
	return e.Severity > ee.Severity
}
