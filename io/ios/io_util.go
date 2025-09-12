package ios

import (
	"fmt"
	"github.com/mobilemindtech/go-io/option"
	"github.com/mobilemindtech/go-io/result"
	"github.com/mobilemindtech/go-io/types"
	"github.com/mobilemindtech/go-io/util"
	"log"
	"reflect"
	"runtime/debug"
)

func TryGetLastIOResult[A any](refIO interface{}, prevEff *option.Option[types.IOEffect]) *result.Result[*option.Option[A]] {
	r := prevEff.Get().GetResult()
	val := r.Get().GetValue()
	if effValue, ok := val.(A); ok {
		return result.OfValue(option.Some(effValue))
	} else {

		ioName := reflect.TypeOf(refIO).Elem().Name()

		util.PanicCastType(ioName,
			reflect.TypeOf(val), reflect.TypeFor[A]())
	}
	return nil // never execute
}

func RecoverIO[A any](refIO interface{}, isDebug bool, debugInfo *types.IODebugInfo, err interface{}) *result.Result[*option.Option[A]] {

	errIO := types.NewIOError(fmt.Sprintf("%v", err), debug.Stack())

	if isDebug {

		ioName := reflect.TypeOf(refIO).Elem().Name()

		if debugInfo != nil {
			log.Printf("[DEBUG %v]=>> added in: %v:%v", ioName, debugInfo.Filename, debugInfo.Line)
		}
		log.Printf("[DEBUG %v]=>> Error: %v\n", ioName, errIO.Error())
		log.Printf("[DEBUG %v]=>> StackTrace: %v\n", ioName, errIO.StackTrace)
	}

	return result.OfError[*option.Option[A]](errIO)
}

func ResultToResultOption[T any](res *result.Result[T]) *result.Result[*option.Option[T]] {
	if res.IsError() {
		return result.OfError[*option.Option[T]](res.GetError())
	}
	return result.OfValue(option.Of(res.Get()))
}
