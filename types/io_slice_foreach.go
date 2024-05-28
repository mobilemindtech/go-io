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

type IOSliceForeach[A any] struct {
	value      *result.Result[*option.Option[[]A]]
	prevEffect IOEffect
	f          func(A)
	debug      bool
	state      *state.State
}

func NewSliceForeach[A any](f func(A)) *IOSliceForeach[A] {
	return &IOSliceForeach[A]{f: f}
}

func (this *IOSliceForeach[A]) SetState(st *state.State) {
	this.state = st
}

func (this *IOSliceForeach[A]) SetDebug(b bool) {
	this.debug = b
}

func (this *IOSliceForeach[A]) String() string {
	return fmt.Sprintf("SliceForeach(%v)", this.value.String())
}

func (this *IOSliceForeach[A]) SetPrevEffect(prev IOEffect) {
	this.prevEffect = prev
}

func (this *IOSliceForeach[A]) GetPrevEffect() *option.Option[IOEffect] {
	return option.Of(this.prevEffect)
}

func (this *IOSliceForeach[A]) GetResult() ResultOptionAny {
	return this.value.ToResultOfOption()
}

func (this *IOSliceForeach[A]) UnsafeRun() IOEffect {
	var currEff interface{} = this
	prevEff := this.GetPrevEffect()
	this.value = result.OfValue(option.None[[]A]())

	if prevEff.NonEmpty() {
		r := prevEff.Get().GetResult()
		if r.IsError() {
			this.value = result.OfError[*option.Option[[]A]](r.Failure())
		} else if r.Get().NonEmpty() {

			val := r.Get().GetValue()

			if effValue, ok := val.([]A); ok {
				this.value = result.OfValue(option.Some(effValue))

				for _, it := range effValue {
					this.f(it)
				}

			} else {
				util.PanicCastType("IOSliceForeach",
					reflect.TypeOf(val), reflect.TypeFor[A]())

			}

		}
	}

	if this.debug {
		log.Printf("%v\n", this.String())
	}

	return currEff.(IOEffect)
}
