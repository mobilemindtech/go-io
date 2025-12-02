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
)

type IOAttemptExec[A any] struct {
	value      *result.Result[*option.Option[A]]
	prevEffect types.IOEffect

	fnExec      func(A)
	fnExecState func(A, *state.State)

	state     *state.State
	debug     bool
	debugInfo *types.IODebugInfo
}

func NewAttemptExec[A any](f func(A)) *IOAttemptExec[A] {
	return &IOAttemptExec[A]{fnExec: f}
}

func NewAttemptExecWithState[A any](f func(A, *state.State)) *IOAttemptExec[A] {
	return &IOAttemptExec[A]{fnExecState: f}
}

func (this *IOAttemptExec[A]) Lift() *types.IO[A] {
	return types.NewIO[A]().Effects(this)
}

func (this *IOAttemptExec[A]) SetDebug(b bool) {
	this.debug = b
}

func (this *IOAttemptExec[A]) SetDebugInfo(info *types.IODebugInfo) {
	this.debugInfo = info
}

func (this *IOAttemptExec[A]) GetDebugInfo() *types.IODebugInfo {
	return this.debugInfo
}

func (this *IOAttemptExec[A]) TypeIn() reflect.Type {
	return reflect.TypeFor[A]()
}

func (this *IOAttemptExec[A]) TypeOut() reflect.Type {
	return reflect.TypeFor[A]()
}

func (this *IOAttemptExec[A]) String() string {
	return fmt.Sprintf("AttemptExec(fn=%v, value=%v)", this.getFuncName(), this.value.String())
}

func (this *IOAttemptExec[A]) SetState(st *state.State) {
	this.state = st
}

func (this *IOAttemptExec[A]) SetPrevEffect(prev types.IOEffect) {
	this.prevEffect = prev
}

func (this *IOAttemptExec[A]) GetPrevEffect() *option.Option[types.IOEffect] {
	return option.Of(this.prevEffect)
}

func (this *IOAttemptExec[A]) GetResult() types.ResultOptionAny {
	return this.value.ToResultOfOption()
}

func (this *IOAttemptExec[A]) UnsafeRun() types.IOEffect {

	var currEff interface{} = this
	prevEff := this.GetPrevEffect()
	this.value = result.OfValue(option.None[A]())
	isEmpty := false
	execute := true
	hasPrev := prevEff.NonEmpty()

	log.Printf("==>>>> %v\n", this.String())

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

		if !isEmpty {
			r := prevEff.Get().GetResult()
			val := r.Get().GetValue()

			if effValue, ok := val.(A); ok {
				if this.fnExec != nil {
					this.fnExec(effValue)
				} else {
					this.fnExecState(effValue, this.state)
				}
			} else {
				util.PanicCastType("IOAttempt",
					reflect.TypeOf(val), reflect.TypeFor[A]())
			}
		}
	}
	if this.debug {
		log.Printf("%v\n", this.String())
	}

	return currEff.(types.IOEffect)
}

func (this *IOAttemptExec[A]) getFuncName() string {
	if this.fnExec != nil {
		return "fnExec"
	}
	if this.fnExecState != nil {
		return "fnExecState"
	}
	return "-"
}
