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

type IOSliceMap[A any, B any] struct {
	value      *result.Result[*option.Option[[]B]]
	prevEffect IOEffect
	f          func(A) B
	debug      bool
	state      *state.State
}

func NewSliceMap[A any, B any](f func(A) B) *IOSliceMap[A, B] {
	return &IOSliceMap[A, B]{f: f}
}

func (this *IOSliceMap[A, B]) SetState(st *state.State) {
	this.state = st
}

func (this *IOSliceMap[A, B]) SetDebug(b bool) {
	this.debug = b
}

func (this *IOSliceMap[A, B]) String() string {
	return fmt.Sprintf("SliceMap(%v)", this.value.String())
}

func (this *IOSliceMap[A, B]) SetPrevEffect(prev IOEffect) {
	this.prevEffect = prev
}

func (this *IOSliceMap[A, B]) GetPrevEffect() *option.Option[IOEffect] {
	return option.Of(this.prevEffect)
}

func (this *IOSliceMap[A, B]) GetResult() ResultOptionAny {
	return this.value.ToResultOfOption()
}

func (this *IOSliceMap[A, B]) UnsafeRun() IOEffect {
	var currEff interface{} = this
	prevEff := this.GetPrevEffect()
	this.value = result.OfValue(option.None[[]B]())

	if prevEff.NonEmpty() {
		r := prevEff.Get().GetResult()
		if r.IsError() {
			this.value = result.OfError[*option.Option[[]B]](r.Failure())
		} else if r.Get().NonEmpty() {
			val := r.Get().GetValue()
			if effValue, ok := val.([]A); ok {
				var list []B
				for _, it := range effValue {
					list = append(list, this.f(it))
				}
				this.value = result.OfValue(option.Some(list))
			} else {
				util.PanicCastType("IOSliceMap",
					reflect.TypeOf(val), reflect.TypeFor[B]())

			}
		}
	}

	if this.debug {
		log.Printf("%v\n", this.String())
	}

	return currEff.(IOEffect)
}
