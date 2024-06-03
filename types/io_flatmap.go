package types

import (
	"fmt"
	"github.com/mobilemindtec/go-io/option"
	"github.com/mobilemindtec/go-io/result"
	"github.com/mobilemindtec/go-io/state"
	"github.com/mobilemindtec/go-io/util"
	"reflect"
)

type IOFlatMap[A any, B any] struct {
	value      *result.Result[*option.Option[B]]
	prevEffect IOEffect
	f          func(A) *IO[B]
	debug      bool
	state      *state.State
}

func NewFlatMap[A any, B any](f func(A) *IO[B]) *IOFlatMap[A, B] {
	return &IOFlatMap[A, B]{f: f}
}

func (this *IOFlatMap[A, B]) SetState(st *state.State) {
	this.state = st
}

func (this *IOFlatMap[A, B]) TypeIn() reflect.Type {
	return reflect.TypeFor[A]()
}

func (this *IOFlatMap[A, B]) TypeOut() reflect.Type {
	return reflect.TypeFor[B]()
}

func (this *IOFlatMap[A, B]) SetDebug(b bool) {
	this.debug = b
}

func (this *IOFlatMap[A, B]) String() string {
	return fmt.Sprintf("FlatMap(%v)", this.value.String())
}

func (this *IOFlatMap[A, B]) SetPrevEffect(prev IOEffect) {
	this.prevEffect = prev
}

func (this *IOFlatMap[A, B]) GetPrevEffect() *option.Option[IOEffect] {
	return option.Of(this.prevEffect)
}

func (this *IOFlatMap[A, B]) GetResult() ResultOptionAny {
	return this.value.ToResultOfOption()
}

func (this *IOFlatMap[A, B]) UnsafeRun() IOEffect {
	var currEff interface{} = this
	prevEff := this.GetPrevEffect()
	this.value = result.OfValue(option.None[B]())

	if prevEff.NonEmpty() {
		r := prevEff.Get().GetResult()

		if r.IsError() {
			this.value = result.OfError[*option.Option[B]](r.Failure())
		} else if r.Get().NonEmpty() {
			val := r.Get().GetValue()
			if effValue, ok := val.(A); ok {
				ioEffect := this.f(effValue)
				ioEffect.SetState(this.state.Copy())
				ioEffect.SetDebug(this.debug)
				this.value = ioEffect.UnsafeRun()
			} else {
				util.PanicCastType("IOFlatMap",
					reflect.TypeOf(val), reflect.TypeFor[B]())
			}
		}
	}

	if this.debug {
		fmt.Println(this)
	}

	return currEff.(IOEffect)
}
