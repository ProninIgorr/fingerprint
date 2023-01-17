package logging

import (
	cou "github.com/ProninIgorr/fingerprint/internal/contexts"
	"github.com/sirupsen/logrus"
)

type ContextFormatter struct {
	BaseFormatter logrus.Formatter
}

func (f *ContextFormatter) Format(e *logrus.Entry) ([]byte, error) {
	if ctx := e.Context; nil != ctx {
		if ops := cou.GetContextOperations(ctx).String(); ops != "" {
			e.Data["ops"] = ops
		}
	}
	return f.BaseFormatter.Format(e)
}
