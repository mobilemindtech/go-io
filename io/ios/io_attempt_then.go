package ios

import (
	"fmt"
	"github.com/mobilemindtech/go-io/option"
	"github.com/mobilemindtech/go-io/result"
	"github.com/mobilemindtech/go-io/runtime"
	"github.com/mobilemindtech/go-io/state"
	"github.com/mobilemindtech/go-io/types"
	"github.com/mobilemindtech/go-io/util"
	"log"
	"reflect"
)

type IOAttemptThen[A any] struct {
	value      *result.Result[*option.Option[A]]
	prevEffect types.IOEffect

	fnPipe      func(A) *result.Result[A]
	fnPipeState func(A, *state.State) *result.Result[A]

	fnPipeOption      func(A) *result.Result[*option.Option[A]]
	fnPipeOptionState func(A, *state.State) *result.Result[*option.Option[A]]

	fnPipeIO      func(A) *types.IO[A]
	fnPipeIOState func(A, *state.State) *types.IO[A]

	state     *state.State
	debug     bool
	debugInfo *types.IODebugInfo
}

func NewAttemptThen[A any](f func(A) *result.Result[A]) *IOAttemptThen[A] {
	return &IOAttemptThen[A]{fnPipe: f}
}

func NewAttemptThenWithState[A any](f func(A, *state.State) *result.Result[A]) *IOAttemptThen[A] {
	return &IOAttemptThen[A]{fnPipeState: f}
}

func NewAttemptThenOption[A any](f func(A) *result.Result[*option.Option[A]]) *IOAttemptThen[A] {
	return &IOAttemptThen[A]{fnPipeOption: f}
}

func NewAttemptThenOptionWithState[A any](f func(A, *state.State) *result.Result[*option.Option[A]]) *IOAttemptThen[A] {
	return &IOAttemptThen[A]{fnPipeOptionState: f}
}

func NewAttemptThenIO[A any](f func(A) *types.IO[A]) *IOAttemptThen[A] {
	return &IOAttemptThen[A]{fnPipeIO: f}
}

func NewAttemptThenIOWithState[A any](f func(A, *state.State) *types.IO[A]) *IOAttemptThen[A] {
	return &IOAttemptThen[A]{fnPipeIOState: f}
}

func (this *IOAttemptThen[A]) Lift() *types.IO[A] {
	return types.NewIO[A]().Effects(this)
}

func (this *IOAttemptThen[A]) SetDebug(b bool) {
	this.debug = b
}

func (this *IOAttemptThen[A]) SetDebugInfo(info *types.IODebugInfo) {
	this.debugInfo = info
}

func (this *IOAttemptThen[A]) GetDebugInfo() *types.IODebugInfo {
	return this.debugInfo
}

func (this *IOAttemptThen[A]) TypeIn() reflect.Type {
	return reflect.TypeFor[A]()
}

func (this *IOAttemptThen[A]) TypeOut() reflect.Type {
	return reflect.TypeFor[A]()
}

func (this *IOAttemptThen[A]) String() string {
	return fmt.Sprintf("AttemptThen(fn=%v, value=%v)", this.getFuncName(), this.value.String())
}

func (this *IOAttemptThen[A]) SetState(st *state.State) {
	this.state = st
}

func (this *IOAttemptThen[A]) SetPrevEffect(prev types.IOEffect) {
	this.prevEffect = prev
}

func (this *IOAttemptThen[A]) GetPrevEffect() *option.Option[types.IOEffect] {
	return option.Of(this.prevEffect)
}

func (this *IOAttemptThen[A]) GetResult() types.ResultOptionAny {
	return this.value.ToResultOfOption()
}

func (this *IOAttemptThen[A]) UnsafeRun() types.IOEffect {

	var currEff interface{} = this
	prevEff := this.GetPrevEffect()
	this.value = result.OfValue(option.None[A]())
	isEmpty := false
	execute := true
	hasPrev := prevEff.NonEmpty()

	if hasPrev {
		prev := prevEff.Get()
		if prev.GetResult().IsError() {
			this.value = result.OfError[*option.Option[A]](
				prevEff.Get().GetResult().Failure())
			execute = false
		} else {
			isEmpty = prev.GetResult().Get().IsEmpty()
		}
	}

	defer func() {
		if r := recover(); r != nil {
			this.value = RecoverIO[A](this, this.debug, this.debugInfo, r)
		}
	}()

	if execute && !isEmpty {

		r := prevEff.Get().GetResult()
		val := r.Get().GetValue()

		if effValue, ok := val.(A); ok {

			if this.fnPipe != nil || this.fnPipeState != nil {

				var pr *result.Result[A]

				if this.fnPipe != nil {
					pr = this.fnPipe(effValue)
				} else {
					pr = this.fnPipeState(effValue, this.state)
				}

				pr.IfError(func(err error) {
					this.value = result.OfError[*option.Option[A]](err)
				}).IfOk(func(a A) {
					this.value = result.OfValue(option.Some(a))
				})

			} else if this.fnPipeOption != nil {
				this.value = this.fnPipeOption(effValue)
			} else if this.fnPipeOptionState != nil {
				this.value = this.fnPipeOptionState(effValue, this.state)
			} else {
				var runnableIO types.IORunnable

				if this.fnPipeIO != nil {
					runnableIO = this.fnPipeIO(effValue)
				} else {
					runnableIO = this.fnPipeIOState(effValue, this.state)
				}

				this.value = runtime.NewWithState[A](this.state, runnableIO).UnsafeRun()
			}

		} else {
			util.PanicCastType("IOAttempt",
				reflect.TypeOf(val), reflect.TypeFor[A]())
		}

	}
	if this.debug {
		log.Printf("%v\n", this.String())
	}

	return currEff.(types.IOEffect)
}

func (this *IOAttemptThen[A]) getFuncName() string {
	if this.fnPipe != nil {
		return "fnPipe"
	}
	if this.fnPipeState != nil {
		return "fnPipeState"
	}
	if this.fnPipeOption != nil {
		return "fnPipeOption"
	}
	if this.fnPipeOptionState != nil {
		return "fnPipeOptionState"
	}
	return "-"
}
