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

type IOExec[A any] struct {
	value      *result.Result[*option.Option[A]]
	prevEffect IOEffect
	f          func(A)
	fstate     func(A, *state.State)
	debug      bool
	state      *state.State
}

func NewExec[A any](f func(A)) *IOExec[A] {
	return &IOExec[A]{f: f}
}

func NewExecState[A any](f func(A, *state.State)) *IOExec[A] {
	return &IOExec[A]{fstate: f}
}

func (this *IOExec[A]) SetState(st *state.State) {
	this.state = st
}

func (this *IOExec[A]) SetDebug(b bool) {
	this.debug = b
}

func (this *IOExec[A]) String() string {
	return fmt.Sprintf("Exec(%v)", this.value.String())
}

func (this *IOExec[A]) SetPrevEffect(prev IOEffect) {
	this.prevEffect = prev
}

func (this *IOExec[A]) GetPrevEffect() *option.Option[IOEffect] {
	return option.Of(this.prevEffect)
}

func (this *IOExec[A]) GetResult() ResultOptionAny {
	return this.value.ToResultOfOption()
}

func (this *IOExec[A]) UnsafeRun() IOEffect {
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

				if this.f != nil {
					this.f(effValue)
				} else {
					this.fstate(effValue, this.state)
				}
			} else {
				util.PanicCastType("IOExec",
					reflect.TypeOf(val), reflect.TypeFor[A]())

			}
		}
	}

	if this.debug {
		log.Printf("%v\n", this.String())
	}

	return currEff.(IOEffect)
}
