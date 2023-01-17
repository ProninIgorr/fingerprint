package logging

import (
	"github.com/ProninIgorr/fingerprint/internal/errs"
	"github.com/sirupsen/logrus"
)

func LogError(args ...interface{}) {
	err := errs.E(args...)
	lvl := Severity2LogLevel[err.Severity()]
	logrus.WithFields(errorFields(err)).Log(lvl, err)
	//if lvl == logrus.FatalLevel{
	//	logrus.Exit(1)
	//}
}

func errorFields(err errs.Error) map[string]interface{} {
	fields := make(map[string]interface{}, 8)
	fields["rec"] = "error"
	fields["type"] = "errs.Error" //fmt.Sprintf("%T", err) = *errs.errorData
	fields["severity"] = err.Severity().String()
	if err.Kind() != errs.KindOther {
		fields["kind"] = err.Kind().String()
	}
	if !err.OperationPath().Empty() {
		fields["ops"] = err.OperationPath().String()
	}
	if len(err.StackTrace()) > 0 {
		fields["frames"] = err.StackTrace()
	}
	fields["cts"] = err.TimeStamp().Format(DefaultTimeFormat)
	return fields
}

var Severity2LogLevel = map[errs.Severity]logrus.Level{
	errs.SeverityWarning:  logrus.WarnLevel,
	errs.SeverityError:    logrus.ErrorLevel,
	errs.SeverityCritical: logrus.FatalLevel,
}

func GetSeveritiesFilter4CurrentLogLevel() (result []errs.Severity) {
	for severity, logLevel := range Severity2LogLevel {
		if logrus.IsLevelEnabled(logLevel) {
			result = append(result, severity)
		}
	}
	return
}
