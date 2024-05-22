package types

import (
	"github.com/mobilemindtec/go-io/option"
	"github.com/mobilemindtec/go-io/result"
)

type IOFilter[A any] struct {
	valueA     *result.Result[A]
	prevEffect IOEffect
	f          func(A) bool
}

func NewFilter[A any](f func(A) bool) *IOFilter[A] {
	return &IOFilter[A]{f: f}
}

func (this *IOFilter[A]) SetPrevEffect(prev IOEffect) {
	this.prevEffect = prev
}

func (this *IOFilter[A]) GetPrevEffect() *option.Option[IOEffect] {
	return option.Of(this.prevEffect)
}

func (this *IOFilter[A]) GetResult() *result.Result[any] {
	return this.valueA.ToResultOfAny()
}

func (this *IOFilter[A]) UnsafeRun() IOEffect {
	var currEff interface{} = this
	prevEff := this.GetPrevEffect()
	if prevEff.NonEmpty() {
		r := prevEff.Get().GetResult()
		if r.IsError() {
			this.valueA = result.OfError[A](r.Error())
		} else if r.OptionNonEmpty() {
			effValue := prevEff.Get().GetResult().ToOption().OrNil().(A)
			if this.f(effValue) {
				this.valueA = result.OfValue(effValue)
			} else {
				this.valueA = result.OfNil[A]()
			}
		}
	}
	return currEff.(IOEffect)
}
