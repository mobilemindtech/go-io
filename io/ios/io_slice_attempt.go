package ios

import (
	"fmt"
	"github.com/mobilemindtech/go-io/option"
	"github.com/mobilemindtech/go-io/result"
	"github.com/mobilemindtech/go-io/state"
	"github.com/mobilemindtech/go-io/types"
	"github.com/mobilemindtech/go-io/util"
	"log"
	"reflect"
	"runtime/debug"
)

type IOSliceAttempt[A any] struct {
	value      *result.Result[*option.Option[[]A]]
	prevEffect types.IOEffect
	f          func([]A) *result.Result[[]A]
	fState     func([]A, *state.State) *result.Result[[]A]
	debug      bool
	debugInfo  *types.IODebugInfo
	state      *state.State
}

func NewSliceAttempt[A any](f func([]A) *result.Result[[]A]) *IOSliceAttempt[A] {
	return &IOSliceAttempt[A]{f: f}
}

func NewSliceAttemptWithState[A any](f func([]A, *state.State) *result.Result[[]A]) *IOSliceAttempt[A] {
	return &IOSliceAttempt[A]{fState: f}
}

func (this *IOSliceAttempt[A]) Lift() *types.IO[A] {
	return types.NewIO[A]().Effects(this)
}

func (this *IOSliceAttempt[A]) TypeIn() reflect.Type {
	return reflect.TypeFor[[]A]()
}

func (this *IOSliceAttempt[A]) TypeOut() reflect.Type {
	return reflect.TypeFor[[]A]()
}

func (this *IOSliceAttempt[A]) SetDebug(b bool) {
	this.debug = b
}

func (this *IOSliceAttempt[A]) SetState(st *state.State) {
	this.state = st
}

func (this *IOSliceAttempt[A]) SetDebugInfo(info *types.IODebugInfo) {
	this.debugInfo = info
}

func (this *IOSliceAttempt[A]) GetDebugInfo() *types.IODebugInfo {
	return this.debugInfo
}

func (this *IOSliceAttempt[A]) String() string {
	return fmt.Sprintf("SliceAttempt(%v)", this.value.String())
}

func (this *IOSliceAttempt[A]) SetPrevEffect(prev types.IOEffect) {
	this.prevEffect = prev
}

func (this *IOSliceAttempt[A]) GetPrevEffect() *option.Option[types.IOEffect] {
	return option.Of(this.prevEffect)
}

func (this *IOSliceAttempt[A]) GetResult() types.ResultOptionAny {
	return this.value.ToResultOfOption()
}

func (this *IOSliceAttempt[A]) UnsafeRun() types.IOEffect {
	var currEff interface{} = this
	prevEff := this.GetPrevEffect()
	this.value = result.OfValue(option.None[[]A]())

	if prevEff.NonEmpty() {
		r := prevEff.Get().GetResult()
		if r.IsError() {
			this.value = result.OfError[*option.Option[[]A]](r.Failure())
		} else if !r.Get().IsEmpty() {

			defer func() {
				if r := recover(); r != nil {

					err := types.NewIOError(fmt.Sprintf("%v", r), debug.Stack())

					if this.debug {
						if this.debugInfo != nil {
							log.Printf("[DEBUG IOSliceAttempt]=>> added in: %v:%v", this.debugInfo.Filename, this.debugInfo.Line)
						}
						log.Printf("[DEBUG IOSliceAttempt]=>> Error: %v\n", err.Error())
						log.Printf("[DEBUG IOSliceAttempt]=>> StackTrace: %v\n", err.StackTrace)
					}

					this.value = result.OfError[*option.Option[[]A]](err)
				}
			}()

			val := r.Get().GetValue()
			if effValue, ok := val.([]A); ok {

				var res *result.Result[[]A]

				if this.f != nil {
					res = this.f(effValue)
				} else {
					res = this.fState(effValue, this.state)
				}

				if res.IsError() {
					this.value = result.OfError[*option.Option[[]A]](res.Failure())
				} else {
					if len(res.Get()) > 0 {
						this.value = result.OfValue(option.Some(res.Get()))
					}
				}

			} else {
				util.PanicCastType("IOSliceAttempt",
					reflect.TypeOf(val), reflect.TypeFor[[]A]())
			}

		}
	}

	if this.debug {
		log.Printf("%v\n", this.String())
	}

	return currEff.(types.IOEffect)
}
