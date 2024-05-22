package types

import (
	"github.com/mobilemindtec/go-io/option"
	"github.com/mobilemindtec/go-io/result"
)

type IOPure[T any] struct {
	value      *option.Option[T]
	prevEffect IOEffect
	f          func() T
}

func NewPure[T any](value T) *IOPure[T] {
	return &IOPure[T]{value: option.Of(value)}
}

func NewPureF[T any](f func() T) *IOPure[T] {
	return &IOPure[T]{f: f}
}

func (this *IOPure[T]) SetPrevEffect(prev IOEffect) {
	this.prevEffect = prev
}

func (this *IOPure[T]) GetPrevEffect() *option.Option[IOEffect] {
	return option.Of(this.prevEffect)
}

func (this *IOPure[T]) GetResult() *result.Result[any] {
	return result.OfValue[any](this.value.OrNil())
}

func (this *IOPure[T]) UnsafeRun() IOEffect {
	var eff interface{} = this

	if this.f != nil {
		this.value = option.Of(this.f())
	}

	return eff.(IOEffect)
}
