package logging

import (
	"context"
	"fmt"
	cou "github.com/ProninIgorr/fingerprint/internal/contexts"
	"github.com/sirupsen/logrus"
	"time"
)

func LogMsg(args ...interface{}) *logrus.Entry {
	entry := logrus.WithFields(logrus.Fields{
		"cts":  time.Now().Format(DefaultTimeFormat),
		"rec":  "msg",
		"type": "string",
	})
	if len(args) == 1 {
		switch arg := args[0].(type) {
		case cou.Operation, cou.Operations:
			return entry.WithField("ops", fmt.Sprintf("%s", arg))
		case context.Context:
			return entry.WithContext(arg)
		}
	}
	return entry
}
