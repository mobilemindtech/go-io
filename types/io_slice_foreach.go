package types

import (
	"github.com/mobilemindtec/go-io/option"
	"github.com/mobilemindtec/go-io/result"
)

type IOSliceForeach[A any] struct {
	valueA     *result.Result[[]A]
	prevEffect IOEffect
	f          func(A)
}

func NewSliceForeach[A any](f func(A)) *IOSliceForeach[A] {
	return &IOSliceForeach[A]{f: f}
}

func (this *IOSliceForeach[A]) SetPrevEffect(prev IOEffect) {
	this.prevEffect = prev
}

func (this *IOSliceForeach[A]) GetPrevEffect() *option.Option[IOEffect] {
	return option.Of(this.prevEffect)
}

func (this *IOSliceForeach[A]) GetResult() *result.Result[any] {
	return this.valueA.ToResultOfAny()
}

func (this *IOSliceForeach[A]) UnsafeRun() IOEffect {
	var currEff interface{} = this
	prevEff := this.GetPrevEffect()
	if prevEff.NonEmpty() {
		r := prevEff.Get().GetResult()
		if r.IsError() {
			this.valueA = result.OfError[[]A](r.Error())
		} else if r.OptionNonEmpty() {
			effValue := prevEff.Get().GetResult().ToOption().OrNil().([]A)
			this.valueA = result.OfValue(effValue)
			for _, it := range effValue {
				this.f(it)
			}

		}
	}
	return currEff.(IOEffect)
}
