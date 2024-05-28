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

type IOSliceFlatMap[A any, B any] struct {
	value      *result.Result[*option.Option[[]B]]
	prevEffect IOEffect
	f          func(A) *IO[B]
	debug      bool
	state      *state.State
}

func NewSliceFlatMap[A any, B any](f func(A) *IO[B]) *IOSliceFlatMap[A, B] {
	return &IOSliceFlatMap[A, B]{f: f}
}

func (this *IOSliceFlatMap[A, B]) SetState(st *state.State) {
	this.state = st
}

func (this *IOSliceFlatMap[A, B]) SetDebug(b bool) {
	this.debug = b
}

func (this *IOSliceFlatMap[A, B]) String() string {
	return fmt.Sprintf("SliceFlatMap(%v)", this.value.String())
}

func (this *IOSliceFlatMap[A, B]) SetPrevEffect(prev IOEffect) {
	this.prevEffect = prev
}

func (this *IOSliceFlatMap[A, B]) GetPrevEffect() *option.Option[IOEffect] {
	return option.Of(this.prevEffect)
}

func (this *IOSliceFlatMap[A, B]) GetResult() ResultOptionAny {
	return this.value.ToResultOfOption()
}

func (this *IOSliceFlatMap[A, B]) UnsafeRun() IOEffect {
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

					ioEffect := this.f(it)
					ioEffect.SetState(this.state.Copy())
					ioEffect.SetDebug(this.debug)
					resultIO := ioEffect.UnsafeRun()

					if resultIO.IsError() {
						this.value = result.OfError[*option.Option[[]B]](resultIO.Failure())
					} else if resultIO.Get().NonEmpty() {
						list = append(list, resultIO.Get().Get())
					}

				}
				if len(list) > 0 {
					this.value = result.OfValue(option.Some(list))
				}

			} else {
				util.PanicCastType("IOSliceFlatMap",
					reflect.TypeOf(val), reflect.TypeFor[B]())

			}
		}
	}

	if this.debug {
		log.Printf("%v\n", this.String())
	}

	return currEff.(IOEffect)
}
