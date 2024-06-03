package types

import (
	"fmt"
	"github.com/mobilemindtec/go-io/option"
	"github.com/mobilemindtec/go-io/result"
	"log"
	"reflect"
)

type IOPure[T any] struct {
	value      *option.Option[T]
	prevEffect IOEffect
	f          func() T
	debug      bool
}

func NewPureValue[T any](value T) *IOPure[T] {
	return &IOPure[T]{value: option.Of(value)}
}

func NewPure[T any](f func() T) *IOPure[T] {
	return &IOPure[T]{f: f}
}

func (this *IOPure[T]) TypeIn() reflect.Type {
	return reflect.TypeFor[T]()
}

func (this *IOPure[T]) TypeOut() reflect.Type {
	return reflect.TypeFor[T]()
}

func (this *IOPure[T]) Lift() *IO[T] {
	return NewIO[T]().Effects(this)
}

func (this *IOPure[T]) SetDebug(b bool) {
	this.debug = b
}

func (this *IOPure[T]) String() string {
	return fmt.Sprintf("Pure(%v)", this.value.String())
}

func (this *IOPure[T]) SetPrevEffect(prev IOEffect) {
	this.prevEffect = prev
}

func (this *IOPure[T]) GetPrevEffect() *option.Option[IOEffect] {
	return option.Of(this.prevEffect)
}

func (this *IOPure[T]) GetResult() ResultOptionAny {

	if this.value.Empty() {
		return result.OfValue(option.None[any]())
	} else {
		return result.OfValue(option.Some[any](this.value.Get()))
	}
}

func (this *IOPure[T]) UnsafeRun() IOEffect {
	var eff interface{} = this

	if this.f != nil {
		this.value = option.Of(this.f())
	}

	if this.debug {
		log.Printf("%v\n", this.String())
	}

	return eff.(IOEffect)
}
