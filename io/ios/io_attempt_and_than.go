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

type IOAttemptAndThan[A any] struct {
	value      *result.Result[*option.Option[A]]
	prevEffect types.IOEffect

	fnState func(*state.State) types.IORunnable
	fn      func() types.IORunnable

	state     *state.State
	debug     bool
	debugInfo *types.IODebugInfo
}

func NewAttemptAndThanWithState[A any](f func(*state.State) types.IORunnable) *IOAttemptAndThan[A] {
	return &IOAttemptAndThan[A]{fnState: f}
}

func NewAttemptAndThan[A any](f func() types.IORunnable) *IOAttemptAndThan[A] {
	return &IOAttemptAndThan[A]{fn: f}
}

func NewAttemptRunIOWithState[A any](f func(*state.State) types.IORunnable) *IOAttemptAndThan[A] {
	return &IOAttemptAndThan[A]{fnState: f}
}

func NewAttemptRunIO[A any](f func() types.IORunnable) *IOAttemptAndThan[A] {
	return &IOAttemptAndThan[A]{fn: f}
}

func (this *IOAttemptAndThan[A]) Lift() *types.IO[A] {
	return types.NewIO[A]().Effects(this)
}

func (this *IOAttemptAndThan[A]) SetDebug(b bool) {
	this.debug = b
}

func (this *IOAttemptAndThan[A]) SetDebugInfo(info *types.IODebugInfo) {
	this.debugInfo = info
}

func (this *IOAttemptAndThan[A]) GetDebugInfo() *types.IODebugInfo {
	return this.debugInfo
}

func (this *IOAttemptAndThan[A]) TypeIn() reflect.Type {
	return reflect.TypeFor[*types.Unit]()
}

func (this *IOAttemptAndThan[A]) TypeOut() reflect.Type {
	return reflect.TypeFor[A]()
}

func (this *IOAttemptAndThan[A]) String() string {
	return fmt.Sprintf("AttemptAndThan(fn=%v, value=%v)", this.getFuncName(), this.value.String())
}

func (this *IOAttemptAndThan[A]) SetState(st *state.State) {
	this.state = st
}

func (this *IOAttemptAndThan[A]) SetPrevEffect(prev types.IOEffect) {
	this.prevEffect = prev
}

func (this *IOAttemptAndThan[A]) GetPrevEffect() *option.Option[types.IOEffect] {
	return option.Of(this.prevEffect)
}

func (this *IOAttemptAndThan[A]) GetResult() types.ResultOptionAny {
	return this.value.ToResultOfOption()
}

func (this *IOAttemptAndThan[A]) UnsafeRun() types.IOEffect {

	var currEff interface{} = this
	prevEff := this.GetPrevEffect()
	this.value = result.OfValue(option.None[A]())
	hasPrev := prevEff.NonEmpty()
	execute := true

	if hasPrev {
		prev := prevEff.Get()
		if prev.GetResult().IsError() {
			this.value = result.OfError[*option.Option[A]](
				prevEff.Get().GetResult().Failure())
			execute = false
		}
	}

	if execute {
		defer func() {
			if r := recover(); r != nil {
				this.value = RecoverIO[A](this, this.debug, this.debugInfo, r)
			}
		}()

		var runnableIO types.IORunnable
		if this.fn != nil {
			runnableIO = this.fn()
		} else {
			runnableIO = this.fnState(this.state)
		}

		if this.debug {
			runnableIO.SetDebug(this.debug)
		}

		this.value = runtime.New[A](runnableIO).UnsafeRun()
	}

	if this.debug {
		log.Printf("%v\n", this.String())
	}

	return currEff.(types.IOEffect)
}

func (this *IOAttemptAndThan[A]) getFuncName() string {
	if this.fnState != nil {
		return "fnState"
	}
	if this.fn != nil {
		return "fn"
	}
	return "-"
}
