package ios

import (
	"fmt"
	"github.com/mobilemindtech/go-io/option"
	"github.com/mobilemindtech/go-io/result"
	"github.com/mobilemindtech/go-io/state"
	"github.com/mobilemindtech/go-io/types"
	"github.com/mobilemindtech/go-io/types/unit"
	"log"
	"reflect"
)

type IOAttempt[A any] struct {
	value      *result.Result[*option.Option[A]]
	prevEffect types.IOEffect

	fnResultOption func() *result.Result[*option.Option[A]]
	fnResult       func() *result.Result[A]
	fnOption       func() *option.Option[A]
	fnError        func() (A, error)

	fnUint      func()
	fnStateUint func(*state.State)

	fnStateResultOption func(*state.State) *result.Result[*option.Option[A]]
	fnStateResult       func(*state.State) *result.Result[A]
	fnStateOption       func(*state.State) *option.Option[A]
	fnStateError        func(*state.State) (A, error)
	fnValueState        func(*state.State) A

	state     *state.State
	debug     bool
	debugInfo *types.IODebugInfo
}

func NewAttempt[A any](f func() *result.Result[A]) *IOAttempt[A] {
	return &IOAttempt[A]{fnResult: f}
}

func NewAttemptOfOption[A any](f func() *option.Option[A]) *IOAttempt[A] {
	return &IOAttempt[A]{fnOption: f}
}

func NewAttemptOfResultOption[A any](f func() *result.Result[*option.Option[A]]) *IOAttempt[A] {
	return &IOAttempt[A]{fnResultOption: f}
}

func NewAttemptState[A any](f func(*state.State) *result.Result[A]) *IOAttempt[A] {
	return &IOAttempt[A]{fnStateResult: f}
}

func NewAttemptStateOfOption[A any](f func(*state.State) *option.Option[A]) *IOAttempt[A] {
	return &IOAttempt[A]{fnStateOption: f}
}

func NewAttemptStateOfResultOption[A any](f func(*state.State) *result.Result[*option.Option[A]]) *IOAttempt[A] {
	return &IOAttempt[A]{fnStateResultOption: f}
}

func NewAttemptOfUnit(f func()) *IOAttempt[*unit.Unit] {
	return &IOAttempt[*unit.Unit]{fnUint: f}
}

func NewAttemptStateOfUnit(f func(*state.State)) *IOAttempt[*unit.Unit] {
	return &IOAttempt[*unit.Unit]{fnStateUint: f}
}

func NewAttemptOfError[A any](f func() (A, error)) *IOAttempt[A] {
	return &IOAttempt[A]{fnError: f}
}

func NewAttemptStateOfError[A any](f func(*state.State) (A, error)) *IOAttempt[A] {
	return &IOAttempt[A]{fnStateError: f}
}

func NewAttemptValueState[A any](f func(*state.State) A) *IOAttempt[A] {
	return &IOAttempt[A]{fnValueState: f}
}

func (this *IOAttempt[T]) Lift() *types.IO[T] {
	return types.NewIO[T]().Effects(this)
}

func (this *IOAttempt[A]) SetDebug(b bool) {
	this.debug = b
}

func (this *IOAttempt[T]) SetDebugInfo(info *types.IODebugInfo) {
	this.debugInfo = info
}

func (this *IOAttempt[T]) GetDebugInfo() *types.IODebugInfo {
	return this.debugInfo
}

func (this *IOAttempt[T]) TypeIn() reflect.Type {
	return reflect.TypeFor[*unit.Unit]()
}
func (this *IOAttempt[T]) TypeOut() reflect.Type {
	return reflect.TypeFor[T]()
}

func (this *IOAttempt[A]) String() string {
	return fmt.Sprintf("Attempt(fn=%v, value=%v)", this.getFuncName(), this.value.String())
}

func (this *IOAttempt[A]) SetState(st *state.State) {
	this.state = st
}

func (this *IOAttempt[A]) SetPrevEffect(prev types.IOEffect) {
	this.prevEffect = prev
}

func (this *IOAttempt[A]) GetPrevEffect() *option.Option[types.IOEffect] {
	return option.Of(this.prevEffect)
}

func (this *IOAttempt[A]) GetResult() types.ResultOptionAny {
	return this.value.ToResultOfOption()
}

func (this *IOAttempt[A]) UnsafeRun() types.IOEffect {

	var currEff interface{} = this
	prevEff := this.GetPrevEffect()
	this.value = result.OfValue(option.None[A]())
	execute := true
	hasPrev := prevEff.NonEmpty()

	if hasPrev {
		prev := prevEff.Get()
		if prev.GetResult().IsError() {
			this.value = result.OfError[*option.Option[A]](
				prevEff.Get().GetResult().Failure())
			execute = false
		} else {
			execute = prev.GetResult().Get().IsSome()
		}
	}

	defer func() {
		if r := recover(); r != nil {
			this.value = RecoverIO[A](this, this.debug, this.debugInfo, r)
		}
	}()

	if execute {
		if this.fnResultOption != nil {
			this.value = this.fnResultOption()
		} else if this.fnOption != nil {
			this.value = result.OfValue(this.fnOption())
		} else if this.fnResult != nil {
			r := this.fnResult()
			if r.HasError() {
				this.value = result.OfError[*option.Option[A]](r.GetError())
			} else {
				this.value = result.OfValue(option.Some(r.Get()))
			}
		} else if this.fnError != nil {
			this.value = result.TryOption(this.fnError)
		} else if this.fnStateResultOption != nil {
			this.value = this.fnStateResultOption(this.state)
		} else if this.fnStateOption != nil {
			this.value = result.OfValue(this.fnStateOption(this.state))
		} else if this.fnStateResult != nil {
			r := this.fnStateResult(this.state)
			if r.HasError() {
				this.value = result.OfError[*option.Option[A]](r.GetError())
			} else {
				this.value = result.OfValue(option.Some(r.Get()))
			}
		} else if this.fnStateError != nil {
			this.value = result.TryOption(
				func() (A, error) {
					return this.fnStateError(this.state)
				})
		} else if this.fnValueState != nil {
			this.value = result.OfValue(option.Of(this.fnValueState(this.state)))
		} else if this.fnUint != nil {
			this.fnUint()
			var unit interface{} = unit.OfUnit()
			this.value = result.OfValue(option.Some(unit.(A)))
		} else if this.fnStateUint != nil {
			this.fnStateUint(this.state)
			var unit interface{} = unit.OfUnit()
			this.value = result.OfValue(option.Some(unit.(A)))
		}
	}

	if this.debug {
		log.Printf("%v\n", this.String())
	}

	return currEff.(types.IOEffect)
}

func (this *IOAttempt[A]) getFuncName() string {
	if this.fnResultOption != nil {
		return "fnResultOption"
	}
	if this.fnResult != nil {
		return "fnResult"
	}
	if this.fnOption != nil {
		return "fnOption"
	}
	if this.fnError != nil {
		return "fnError"
	}
	if this.fnUint != nil {
		return "fnUint"
	}
	if this.fnStateUint != nil {
		return "fnStateUint"
	}
	if this.fnStateResultOption != nil {
		return "fnStateResultOption"
	}
	if this.fnStateResult != nil {
		return "fnStateResult"
	}
	if this.fnStateOption != nil {
		return "fnStateOption"
	}
	if this.fnStateError != nil {
		return "fnStateError"
	}
	if this.fnValueState != nil {
		return "fnValueState"
	}
	return "-"
}
