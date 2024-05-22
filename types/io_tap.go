package types

import (
	"github.com/mobilemindtec/go-io/option"
	"github.com/mobilemindtec/go-io/result"
)

type IOTap[A any] struct {
	valueA     *result.Result[A]
	prevEffect IOEffect
	f          func(A) bool
}

func NewTap[A any](f func(A) bool) *IOTap[A] {
	return &IOTap[A]{f: f}
}

func (this *IOTap[A]) SetPrevEffect(prev IOEffect) {
	this.prevEffect = prev
}

func (this *IOTap[A]) GetPrevEffect() *option.Option[IOEffect] {
	return option.Of(this.prevEffect)
}

func (this *IOTap[A]) GetResult() *result.Result[any] {
	return this.valueA.ToResultOfAny()
}

func (this *IOTap[A]) UnsafeRun() IOEffect {
	var currEff interface{} = this
	prevEff := this.GetPrevEffect()
	if prevEff.NonEmpty() {
		r := prevEff.Get().GetResult()
		if r.IsError() {
			this.valueA = result.OfError[A](r.Error())
		} else if r.OptionNonEmpty() {
			effValue := prevEff.Get().GetResult().ToOption().OrNil().(A)
			this.valueA = result.OfValue(effValue)
			this.f(effValue)
		}
	}
	return currEff.(IOEffect)
}
