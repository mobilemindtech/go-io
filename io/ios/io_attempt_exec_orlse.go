package ios

import (
	"fmt"
	"github.com/mobilemindtec/go-io/option"
	"github.com/mobilemindtec/go-io/result"
	"github.com/mobilemindtec/go-io/state"
	"github.com/mobilemindtec/go-io/types"
	"log"
	"reflect"
)

type IOAttemptExecOrElse[A any] struct {
	value      *result.Result[*option.Option[A]]
	prevEffect types.IOEffect

	fnExecOrElse      func()
	fnExecOrElseState func(*state.State)

	state     *state.State
	debug     bool
	debugInfo *types.IODebugInfo
}

func NewAttemptExecOrElse[A any](f func()) *IOAttemptExecOrElse[A] {
	return &IOAttemptExecOrElse[A]{fnExecOrElse: f}
}

func NewAttemptExecOrElseWithState[A any](f func(*state.State)) *IOAttemptExecOrElse[A] {
	return &IOAttemptExecOrElse[A]{fnExecOrElseState: f}
}

func (this *IOAttemptExecOrElse[A]) Lift() *types.IO[A] {
	return types.NewIO[A]().Effects(this)
}

func (this *IOAttemptExecOrElse[A]) SetDebug(b bool) {
	this.debug = b
}

func (this *IOAttemptExecOrElse[A]) SetDebugInfo(info *types.IODebugInfo) {
	this.debugInfo = info
}

func (this *IOAttemptExecOrElse[A]) GetDebugInfo() *types.IODebugInfo {
	return this.debugInfo
}

func (this *IOAttemptExecOrElse[A]) TypeIn() reflect.Type {
	return reflect.TypeFor[*types.Unit]()
}

func (this *IOAttemptExecOrElse[A]) TypeOut() reflect.Type {
	return reflect.TypeFor[A]()
}

func (this *IOAttemptExecOrElse[A]) String() string {
	return fmt.Sprintf("AttemptExecOrElse(fn=%v, value=%v)", this.getFuncName(), this.value.String())
}

func (this *IOAttemptExecOrElse[A]) SetState(st *state.State) {
	this.state = st
}

func (this *IOAttemptExecOrElse[A]) SetPrevEffect(prev types.IOEffect) {
	this.prevEffect = prev
}

func (this *IOAttemptExecOrElse[A]) GetPrevEffect() *option.Option[types.IOEffect] {
	return option.Of(this.prevEffect)
}

func (this *IOAttemptExecOrElse[A]) GetResult() types.ResultOptionAny {
	return this.value.ToResultOfOption()
}

func (this *IOAttemptExecOrElse[A]) UnsafeRun() types.IOEffect {

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

	if execute {

		if isEmpty { // not error

			if this.fnExecOrElse != nil {
				this.fnExecOrElse()
			} else if this.fnExecOrElse != nil {
				this.fnExecOrElseState(this.state)
			}

		} else {
			this.value = TryGetLastIOResult[A](this, prevEff)
		}
	}
	if this.debug {
		log.Printf("%v\n", this.String())
	}

	return currEff.(types.IOEffect)
}

func (this *IOAttemptExecOrElse[A]) getFuncName() string {
	if this.fnExecOrElse != nil {
		return "fnExecOrElse"
	}
	if this.fnExecOrElseState != nil {
		return "fnExecOrElseState"
	}
	return "-"
}
