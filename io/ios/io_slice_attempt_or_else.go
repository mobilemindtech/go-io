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

type IOSliceAttemptOrElse[A any] struct {
	value      *result.Result[*option.Option[[]A]]
	prevEffect types.IOEffect
	f          func() *result.Result[[]A]
	fstate     func(*state.State) *result.Result[[]A]
	debug      bool
	debugInfo  *types.IODebugInfo
	state      *state.State
}

func NewSliceAttemptOrElse[A any](f func() *result.Result[[]A]) *IOSliceAttemptOrElse[A] {
	return &IOSliceAttemptOrElse[A]{f: f}
}

func NewSliceAttemptOrElseWithState[A any](f func(*state.State) *result.Result[[]A]) *IOSliceAttemptOrElse[A] {
	return &IOSliceAttemptOrElse[A]{fstate: f}
}

func (this *IOSliceAttemptOrElse[A]) Lift() *types.IO[A] {
	return types.NewIO[A]().Effects(this)
}

func (this *IOSliceAttemptOrElse[A]) TypeIn() reflect.Type {
	return reflect.TypeFor[[]A]()
}

func (this *IOSliceAttemptOrElse[A]) TypeOut() reflect.Type {
	return reflect.TypeFor[[]A]()
}

func (this *IOSliceAttemptOrElse[A]) SetDebug(b bool) {
	this.debug = b
}

func (this *IOSliceAttemptOrElse[A]) SetState(st *state.State) {
	this.state = st
}

func (this *IOSliceAttemptOrElse[A]) SetDebugInfo(info *types.IODebugInfo) {
	this.debugInfo = info
}

func (this *IOSliceAttemptOrElse[A]) GetDebugInfo() *types.IODebugInfo {
	return this.debugInfo
}

func (this *IOSliceAttemptOrElse[A]) String() string {
	return fmt.Sprintf("SliceAttemptOrElse(%v)", this.value.String())
}

func (this *IOSliceAttemptOrElse[A]) SetPrevEffect(prev types.IOEffect) {
	this.prevEffect = prev
}

func (this *IOSliceAttemptOrElse[A]) GetPrevEffect() *option.Option[types.IOEffect] {
	return option.Of(this.prevEffect)
}

func (this *IOSliceAttemptOrElse[A]) GetResult() types.ResultOptionAny {
	return this.value.ToResultOfOption()
}

func (this *IOSliceAttemptOrElse[A]) UnsafeRun() types.IOEffect {
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
							log.Printf("[DEBUG IOSliceAttemptOrElse]=>> added in: %v:%v", this.debugInfo.Filename, this.debugInfo.Line)
						}
						log.Printf("[DEBUG IOSliceAttemptOrElse]=>> Error: %v\n", err.Error())
						log.Printf("[DEBUG IOSliceAttemptOrElse]=>> StackTrace: %v\n", err.StackTrace)
					}

					this.value = result.OfError[*option.Option[[]A]](err)
				}
			}()

			var res *result.Result[[]A]

			if this.f != nil {
				res = this.f()
			} else {
				res = this.fstate(this.state)
			}

			if res.IsError() {
				this.value = result.OfError[*option.Option[[]A]](res.Failure())
			} else {
				if len(res.Get()) > 0 {
					this.value = result.OfValue(option.Some(res.Get()))
				}
			}
		} else {
			val := r.Get().GetValue()
			if effValue, ok := val.([]A); ok {
				this.value = result.OfValue(option.Some(effValue))
			} else {
				util.PanicCastType("IOSliceAttemptOrElse",
					reflect.TypeOf(val), reflect.TypeFor[[]A]())
			}
		}
	}

	if this.debug {
		log.Printf("%v\n", this.String())
	}

	return currEff.(types.IOEffect)
}
