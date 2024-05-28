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

type IOSliceFilter[A any] struct {
	value      *result.Result[*option.Option[[]A]]
	prevEffect IOEffect
	f          func(A) bool
	debug      bool
	state      *state.State
}

func NewSliceFilter[A any](f func(A) bool) *IOSliceFilter[A] {
	return &IOSliceFilter[A]{f: f}
}

func (this *IOSliceFilter[A]) SetState(st *state.State) {
	this.state = st
}

func (this *IOSliceFilter[A]) SetDebug(b bool) {
	this.debug = b
}

func (this *IOSliceFilter[A]) String() string {
	return fmt.Sprintf("SliceFilter(%v)", this.value.String())
}

func (this *IOSliceFilter[A]) SetPrevEffect(prev IOEffect) {
	this.prevEffect = prev
}

func (this *IOSliceFilter[A]) GetPrevEffect() *option.Option[IOEffect] {
	return option.Of(this.prevEffect)
}

func (this *IOSliceFilter[A]) GetResult() ResultOptionAny {
	return this.value.ToResultOfOption()
}

func (this *IOSliceFilter[A]) UnsafeRun() IOEffect {
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

				var list []A
				for _, it := range effValue {
					if this.f(it) {
						list = append(list, it)
					}
				}
				if len(list) > 0 {
					this.value = result.OfValue(option.Some(list))
				}

			} else {
				util.PanicCastType("IOSliceFilter",
					reflect.TypeOf(val), reflect.TypeFor[A]())

			}

		}
	}

	if this.debug {
		log.Printf("%v\n", this.String())
	}

	return currEff.(IOEffect)
}
