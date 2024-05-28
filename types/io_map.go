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

type IOMap[A any, B any] struct {
	value      *result.Result[*option.Option[B]]
	prevEffect IOEffect
	f          func(A) B
	debug      bool
	state      *state.State
}

func NewMap[A any, B any](f func(A) B) *IOMap[A, B] {
	return &IOMap[A, B]{f: f}
}

func (this *IOMap[A, B]) SetState(st *state.State) {
	this.state = st
}

func (this *IOMap[A, B]) SetDebug(b bool) {
	this.debug = b
}

func (this *IOMap[A, B]) String() string {
	return fmt.Sprintf("Map(%v)", this.value.String())
}

func (this *IOMap[A, B]) SetPrevEffect(prev IOEffect) {
	this.prevEffect = prev
}

func (this *IOMap[A, B]) GetPrevEffect() *option.Option[IOEffect] {
	return option.Of(this.prevEffect)
}

func (this *IOMap[A, B]) GetResult() ResultOptionAny {
	return this.value.ToResultOfOption()
}

func (this *IOMap[A, B]) UnsafeRun() IOEffect {
	var currEff interface{} = this
	prevEff := this.GetPrevEffect()
	this.value = result.OfValue(option.None[B]())

	if prevEff.NonEmpty() {
		r := prevEff.Get().GetResult()

		if r.IsError() {
			this.value = result.OfError[*option.Option[B]](r.GetError())
		} else if r.Get().NonEmpty() {
			val := r.Get().GetValue()
			if effValue, ok := val.(A); ok {
				this.value = result.OfValue(option.Some(this.f(effValue)))
			} else {
				util.PanicCastType("IOMap",
					reflect.TypeOf(val), reflect.TypeFor[B]())
			}
		}
	}

	if this.debug {
		log.Printf("%v\n", this.String())
	}

	return currEff.(IOEffect)
}
