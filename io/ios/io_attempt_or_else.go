package ios

import (
	"fmt"
	"github.com/mobilemindtec/go-io/option"
	"github.com/mobilemindtec/go-io/result"
	"github.com/mobilemindtec/go-io/runtime"
	"github.com/mobilemindtec/go-io/state"
	"github.com/mobilemindtec/go-io/types"
	"log"
	"reflect"
)

type IOAttemptOrElse[A any] struct {
	value      *result.Result[*option.Option[A]]
	prevEffect types.IOEffect

	fnState func(*state.State) *types.IO[A]
	fn      func() *types.IO[A]

	state     *state.State
	debug     bool
	debugInfo *types.IODebugInfo
}

func NewAttemptOrElseWithState[A any](f func(*state.State) *types.IO[A]) *IOAttemptOrElse[A] {
	return &IOAttemptOrElse[A]{fnState: f}
}

func NewAttemptOrElse[A any](f func() *types.IO[A]) *IOAttemptOrElse[A] {
	return &IOAttemptOrElse[A]{fn: f}
}

func (this *IOAttemptOrElse[A]) Lift() *types.IO[A] {
	return types.NewIO[A]().Effects(this)
}

func (this *IOAttemptOrElse[A]) SetDebug(b bool) {
	this.debug = b
}

func (this *IOAttemptOrElse[A]) SetDebugInfo(info *types.IODebugInfo) {
	this.debugInfo = info
}

func (this *IOAttemptOrElse[A]) GetDebugInfo() *types.IODebugInfo {
	return this.debugInfo
}

func (this *IOAttemptOrElse[A]) TypeIn() reflect.Type {
	return reflect.TypeFor[*types.Unit]()
}

func (this *IOAttemptOrElse[A]) TypeOut() reflect.Type {
	return reflect.TypeFor[A]()
}

func (this *IOAttemptOrElse[A]) String() string {
	return fmt.Sprintf("AttemptOrElse(fn=%v, value=%v)", this.getFuncName(), this.value.String())
}

func (this *IOAttemptOrElse[A]) SetState(st *state.State) {
	this.state = st
}

func (this *IOAttemptOrElse[A]) SetPrevEffect(prev types.IOEffect) {
	this.prevEffect = prev
}

func (this *IOAttemptOrElse[A]) GetPrevEffect() *option.Option[types.IOEffect] {
	return option.Of(this.prevEffect)
}

func (this *IOAttemptOrElse[A]) GetResult() types.ResultOptionAny {
	return this.value.ToResultOfOption()
}

func (this *IOAttemptOrElse[A]) UnsafeRun() types.IOEffect {

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

			var runnableIO types.IORunnable
			if this.fn != nil {
				runnableIO = this.fn()
			} else {
				runnableIO = this.fnState(this.state)
			}

			if this.debug {
				runnableIO.SetDebug(this.debug)
			}

			this.value = runtime.NewWithState[A](this.state, runnableIO).UnsafeRun()

		} else {
			this.value = TryGetLastIOResult[A](this, prevEff)
		}
	}

	if this.debug {
		log.Printf("%v\n", this.String())
	}

	return currEff.(types.IOEffect)
}

func (this *IOAttemptOrElse[A]) getFuncName() string {
	if this.fnState != nil {
		return "fnState"
	}
	if this.fn != nil {
		return "fn"
	}
	return "-"
}
