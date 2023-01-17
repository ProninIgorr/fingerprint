package logging

import (
	"context"
	"fmt"
	"github.com/ProninIgorr/fingerprint/internal/contexts"
	"github.com/ProninIgorr/fingerprint/internal/errs"
	"github.com/ProninIgorr/fingerprint/internal/fh"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"os/user"
	"runtime/trace"
	"strings"
)

var logFile, traceFile *os.File

func Initialize(ctx context.Context, logFileName, level, format, traceFileName string, usr *user.User) errs.Error {
	ctx = contexts.BuildContext(ctx, contexts.AddContextOperation("log_init"), errs.SetDefaultErrsSeverity(errs.SeverityCritical))
	logrus.SetOutput(os.Stdout)
	logrus.RegisterExitHandler(Finalize)
	if logFileName != "" {
		var file *os.File
		var err error
		if logFileName, err := fh.SafeParentResolvePath(logFileName, usr, 0700); err != nil {
			return errs.E(ctx, errs.KindInvalidValue, fmt.Errorf("invalid log file name <%s>: %w", logFileName, err))
		}
		file, err = os.OpenFile(logFileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0664)
		if err != nil {
			return errs.E(ctx, errs.KindOSOpenFile, fmt.Errorf("open file <%s> for logging failed: %w", logFileName, err))
		} else {
			logrus.SetOutput(file)
			logFile = file
			fmt.Println("logging to ", file.Name())
		}
	} else {
		fmt.Println("logging to standard output")
	}

	fieldMap := logrus.FieldMap{
		logrus.FieldKeyTime:  "ts",
		logrus.FieldKeyLevel: "lvl",
		logrus.FieldKeyMsg:   "msg"}
	switch strings.ToUpper(format) {
	case "JSON":
		logrus.SetFormatter(
			&ContextFormatter{
				&logrus.JSONFormatter{
					FieldMap: fieldMap,
				},
			})
	case "TEXT":
		logrus.SetFormatter(
			&ContextFormatter{
				&logrus.TextFormatter{
					ForceQuote:       false,
					DisableTimestamp: false,
					FullTimestamp:    true,
					TimestampFormat:  DefaultTimeFormat,
					QuoteEmptyFields: true,
					FieldMap:         fieldMap,
				},
			})
	default:
		return errs.E(ctx, errs.KindInvalidValue, fmt.Errorf("invalid log format [%s]. supported formats: json, text", format))
	}

	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		return errs.E(ctx, errs.KindInvalidValue, fmt.Errorf("parsing log level from config failed: %w", err))
	}
	logrus.SetLevel(lvl)
	LogMsg(ctx).Infof("Logging initialized with level <%s>", level) // first record in log file

	logrus.SetReportCaller(lvl > logrus.InfoLevel)
	errs.WithFrames(lvl > logrus.InfoLevel)
	if lvl == logrus.TraceLevel {
		var err error
		if traceFileName == "" {
			traceFile = os.Stderr
		} else {
			if traceFileName, err := fh.SafeParentResolvePath(traceFileName, usr, 0700); err != nil {
				return errs.E(ctx, errs.KindInvalidValue, fmt.Errorf("invalid trace file name <%s>: %w", traceFileName, err))
			}
			if traceFile, err = os.OpenFile(traceFileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0664); err != nil {
				return errs.E(ctx, errs.KindOSOpenFile, fmt.Errorf("creating trace file failed: %w", err))
			}
		}
		if err := trace.Start(traceFile); err != nil {
			return errs.E(ctx, errs.KindInternal, fmt.Errorf("starting trace failed: %w", err))
		}
		if traceFile != nil {
			LogMsg(ctx).Tracef("Tracing started with output to %s", traceFile.Name())
		}
	}
	log.SetOutput(logrus.StandardLogger().Writer()) // to use with standard log pkg
	return nil
}

func Finalize() {
	op := contexts.Operation("log_finalize")
	if traceFile != nil {
		trace.Stop()
		LogMsg(op).Trace("trace stopped")
		if traceFile != os.Stderr {
			if err := traceFile.Close(); err != nil {
				LogMsg(op).Errorf("closing trace file failed: %v", err)
			} else {
				LogMsg(op).Trace("trace file closed")
			}
		}
	}
	if nil != logFile {
		LogMsg(op).Debug("Logging finalized.")
		if err := logFile.Sync(); err != nil {
			LogMsg(op).Errorf("sync log buffer with file [%s] - failed: %v", logFile.Name(), err)
		}
		if err := logFile.Close(); err != nil {
			LogMsg(op).Errorf("closing log file [%s] - failed: %v\n", logFile.Name(), err)
		}
	}
}
