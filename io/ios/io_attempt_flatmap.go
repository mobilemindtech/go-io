package ios

import (
	"fmt"
	"github.com/mobilemindtec/go-io/option"
	"github.com/mobilemindtec/go-io/result"
	"github.com/mobilemindtec/go-io/runtime"
	"github.com/mobilemindtec/go-io/state"
	"github.com/mobilemindtec/go-io/types"
	"github.com/mobilemindtec/go-io/util"
	"log"
	"reflect"
)

type IOAttemptFlatMap[A, B any] struct {
	value      *result.Result[*option.Option[B]]
	prevEffect types.IOEffect

	f func(A, *state.State) *types.IO[B]

	state     *state.State
	debug     bool
	debugInfo *types.IODebugInfo
}

func NewAttemptFlatMap[A, B any](f func(A, *state.State) *types.IO[B]) *IOAttemptFlatMap[A, B] {
	return &IOAttemptFlatMap[A, B]{f: f}
}

func (this *IOAttemptFlatMap[A, B]) Lift() *types.IO[B] {
	return types.NewIO[B]().Effects(this)
}

func (this *IOAttemptFlatMap[A, B]) SetDebug(b bool) {
	this.debug = b
}

func (this *IOAttemptFlatMap[A, B]) SetDebugInfo(info *types.IODebugInfo) {
	this.debugInfo = info
}

func (this *IOAttemptFlatMap[A, B]) GetDebugInfo() *types.IODebugInfo {
	return this.debugInfo
}

func (this *IOAttemptFlatMap[A, B]) TypeIn() reflect.Type {
	return reflect.TypeFor[A]()
}

func (this *IOAttemptFlatMap[A, B]) TypeOut() reflect.Type {
	return reflect.TypeFor[B]()
}

func (this *IOAttemptFlatMap[A, B]) String() string {
	return fmt.Sprintf("AttemptFlatMap(fn=%v, value=%v)", "f", this.value.String())
}

func (this *IOAttemptFlatMap[A, B]) SetState(st *state.State) {
	this.state = st
}

func (this *IOAttemptFlatMap[A, B]) SetPrevEffect(prev types.IOEffect) {
	this.prevEffect = prev
}

func (this *IOAttemptFlatMap[A, B]) GetPrevEffect() *option.Option[types.IOEffect] {
	return option.Of(this.prevEffect)
}

func (this *IOAttemptFlatMap[A, B]) GetResult() types.ResultOptionAny {
	return this.value.ToResultOfOption()
}

func (this *IOAttemptFlatMap[A, B]) UnsafeRun() types.IOEffect {

	var currEff interface{} = this
	prevEff := this.GetPrevEffect()
	this.value = result.OfValue(option.None[B]())
	isEmpty := false
	execute := true
	hasPrev := prevEff.NonEmpty()

	if hasPrev {
		prev := prevEff.Get()
		if prev.GetResult().IsError() {
			this.value = result.OfError[*option.Option[B]](
				prevEff.Get().GetResult().Failure())
			execute = false
		} else {
			isEmpty = prev.GetResult().Get().IsEmpty()
		}
	}

	defer func() {
		if r := recover(); r != nil {
			this.value = RecoverIO[B](this, this.debug, this.debugInfo, r)
		}
	}()

	if execute && !isEmpty {

		r := prevEff.Get().GetResult()
		val := r.Get().GetValue()

		if effValue, ok := val.(A); ok {
			runnableIO := this.f(effValue, this.state)
			this.value = runtime.NewWithState[B](this.state, runnableIO).UnsafeRun()
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
