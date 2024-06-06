package ios

import (
	"fmt"
	"github.com/mobilemindtec/go-io/option"
	"github.com/mobilemindtec/go-io/result"
	"github.com/mobilemindtec/go-io/types"
	"log"
	"reflect"
	"runtime/debug"
)

type IOSliceAttemptIfEmpty[A any] struct {
	value      *result.Result[*option.Option[[]A]]
	prevEffect types.IOEffect
	f          func() *result.Result[[]A]
	debug      bool
	debugInfo  *types.IODebugInfo
}

func NewSliceAttemptIfEmpty[A any](f func() *result.Result[[]A]) *IOSliceAttemptIfEmpty[A] {
	return &IOSliceAttemptIfEmpty[A]{f: f}
}

func (this *IOSliceAttemptIfEmpty[A]) TypeIn() reflect.Type {
	return reflect.TypeFor[[]A]()
}

func (this *IOSliceAttemptIfEmpty[A]) TypeOut() reflect.Type {
	return reflect.TypeFor[[]A]()
}

func (this *IOSliceAttemptIfEmpty[A]) SetDebug(b bool) {
	this.debug = b
}

func (this *IOSliceAttemptIfEmpty[A]) SetDebugInfo(info *types.IODebugInfo) {
	this.debugInfo = info
}

func (this *IOSliceAttemptIfEmpty[A]) GetDebugInfo() *types.IODebugInfo {
	return this.debugInfo
}

func (this *IOSliceAttemptIfEmpty[A]) String() string {
	return fmt.Sprintf("SliceAttemptIfEmpty(%v)", this.value.String())
}

func (this *IOSliceAttemptIfEmpty[A]) SetPrevEffect(prev types.IOEffect) {
	this.prevEffect = prev
}

func (this *IOSliceAttemptIfEmpty[A]) GetPrevEffect() *option.Option[types.IOEffect] {
	return option.Of(this.prevEffect)
}

func (this *IOSliceAttemptIfEmpty[A]) GetResult() types.ResultOptionAny {
	return this.value.ToResultOfOption()
}

func (this *IOSliceAttemptIfEmpty[A]) UnsafeRun() types.IOEffect {
	var currEff interface{} = this
	prevEff := this.GetPrevEffect()
	this.value = result.OfValue(option.None[[]A]())

	if prevEff.NonEmpty() {
		r := prevEff.Get().GetResult()
		if r.IsError() {
			this.value = result.OfError[*option.Option[[]A]](r.Failure())
		} else if r.Get().IsEmpty() {

			defer func() {
				if r := recover(); r != nil {

					err := types.NewIOError(fmt.Sprintf("%v", r), debug.Stack())

					if this.debug {
						if this.debugInfo != nil {
							log.Printf("[DEBUG IOSliceAttemptIfEmpty]=>> added in: %v:%v", this.debugInfo.Filename, this.debugInfo.Line)
						}
						log.Printf("[DEBUG IOSliceAttemptIfEmpty]=>> Error: %v\n", err.Error())
						log.Printf("[DEBUG IOSliceAttemptIfEmpty]=>> StackTrace: %v\n", err.StackTrace)
					}

					this.value = result.OfError[*option.Option[[]A]](err)
				}
			}()

			res := this.f()

			if res.IsError() {
				this.value = result.OfError[*option.Option[[]A]](res.Failure())
			} else {
				if len(res.Get()) > 0 {
					this.value = result.OfValue(option.Some(res.Get()))
				}
			}
		}
	}

	if this.debug {
		log.Printf("%v\n", this.String())
	}

	return currEff.(types.IOEffect)
}
