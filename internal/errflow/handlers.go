package errflow

import (
	"sync"

	"github.com/ProninIgorr/fingerprint/internal/errs"
	"github.com/ProninIgorr/fingerprint/internal/logging"
)

//func CriticalErrorHandlerBuilder(cancel context.CancelFunc, kinds []errs.Kind) FuncErrorHandler {
//	mKinds := make(map[errs.Kind]struct{}, len(kinds))
//	for _, kind := range kinds {
//		mKinds[kind] = struct{}{}
//	}
//	return func(cerr <-chan errs.Error, wg *sync.WaitGroup) {
//		defer wg.Done()
//		for err := range cerr {
//			if err != nil {
//				logging.LogError(err)
//				if _, exists := mKinds[err.Kind()]; exists {
//					cancel()
//					break
//				}
//			}
//		}
//	}
//}

func LoggingErrorHandler(cerr <-chan errs.Error, wg *sync.WaitGroup) {
	defer wg.Done()
	for err := range cerr {
		logging.LogError(err)
	}
}
