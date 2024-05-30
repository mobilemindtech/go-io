package types

import (
	"fmt"
	"github.com/mobilemindtec/go-io/option"
	"github.com/mobilemindtec/go-io/result"
	"github.com/mobilemindtec/go-io/state"
	"github.com/mobilemindtec/go-io/util"
	"log"
	"reflect"
)

type IOExecIfEmpty[A any] struct {
	value      *result.Result[*option.Option[A]]
	prevEffect IOEffect
	f          func()
	debug      bool
	state      *state.State
}

func NewExecIfEmpty[A any](f func()) *IOExecIfEmpty[A] {
	return &IOExecIfEmpty[A]{f: f}
}

func (this *IOExecIfEmpty[A]) SetState(st *state.State) {
	this.state = st
}

func (this *IOExecIfEmpty[A]) SetDebug(b bool) {
	this.debug = b
}

func (this *IOExecIfEmpty[A]) String() string {
	return fmt.Sprintf("ExecIfEmpty(%v)", this.value.String())
}

func (this *IOExecIfEmpty[A]) SetPrevEffect(prev IOEffect) {
	this.prevEffect = prev
}

func (this *IOExecIfEmpty[A]) GetPrevEffect() *option.Option[IOEffect] {
	return option.Of(this.prevEffect)
}

func (this *IOExecIfEmpty[A]) GetResult() ResultOptionAny {
	return this.value.ToResultOfOption()
}

func (this *IOExecIfEmpty[A]) UnsafeRun() IOEffect {
	var currEff interface{} = this
	prevEff := this.GetPrevEffect()
	this.value = result.OfValue(option.None[A]())

	if prevEff.NonEmpty() {
		r := prevEff.Get().GetResult()

		if r.IsError() {
			this.value = result.OfError[*option.Option[A]](r.Failure())
		} else if r.Get().NonEmpty() {
			val := r.Get().GetValue()
			if effValue, ok := val.(A); ok {
				this.value = result.OfValue(option.Some(effValue))
			} else {
				util.PanicCastType("IOExecIfEmpty",
					reflect.TypeOf(val), reflect.TypeFor[A]())
			}
		} else {
			this.f()
		}
	}

	if this.debug {
		log.Printf("%v\n", this.String())
	}

	return currEff.(IOEffect)
}
