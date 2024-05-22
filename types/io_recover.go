package types

import (
	"github.com/mobilemindtec/go-io/option"
	"github.com/mobilemindtec/go-io/result"
)

type IORecover[A any] struct {
	valueA     *result.Result[A]
	prevEffect IOEffect
	f          func(error) A
}

func NewRecover[A any](f func(error) A) *IORecover[A] {
	return &IORecover[A]{f: f}
}

func (this *IORecover[A]) SetPrevEffect(prev IOEffect) {
	this.prevEffect = prev
}

func (this *IORecover[A]) GetPrevEffect() *option.Option[IOEffect] {
	return option.Of(this.prevEffect)
}

func (this *IORecover[A]) GetResult() *result.Result[any] {
	return this.valueA.ToResultOfAny()
}

func (this *IORecover[A]) UnsafeRun() IOEffect {
	var currEff interface{} = this
	prevEff := this.GetPrevEffect()
	if prevEff.NonEmpty() {
		r := prevEff.Get().GetResult()
		if r.IsError() {
			this.valueA = result.OfValue(this.f(r.Error()))
		} else if r.OptionNonEmpty() {
			this.valueA = result.OfValue(prevEff.Get().(A))
		}
	}
	return currEff.(IOEffect)
}
