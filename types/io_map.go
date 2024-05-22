package types

import (
	"github.com/mobilemindtec/go-io/option"
	"github.com/mobilemindtec/go-io/result"
)

type IOMap[A any, B any] struct {
	//valueA     *option.Option[A]
	valueB     *result.Result[B]
	prevEffect IOEffect
	f          func(A) B
}

func NewMap[A any, B any](f func(A) B) *IOMap[A, B] {
	return &IOMap[A, B]{f: f}
}

func (this *IOMap[A, B]) SetPrevEffect(prev IOEffect) {
	this.prevEffect = prev
}

func (this *IOMap[A, B]) GetPrevEffect() *option.Option[IOEffect] {
	return option.Of(this.prevEffect)
}

func (this *IOMap[A, B]) GetResult() *result.Result[any] {
	return this.valueB.ToResultOfAny()
}

func (this *IOMap[A, B]) UnsafeRun() IOEffect {
	var currEff interface{} = this
	prevEff := this.GetPrevEffect()
	if prevEff.NonEmpty() {
		r := prevEff.Get().GetResult()
		if r.IsError() {
			this.valueB = result.OfError[B](r.GetError())
		} else if r.OptionNonEmpty() {
			effValue := prevEff.Get().GetResult().ToOption().OrNil().(A)
			this.valueB = result.OfValue(this.f(effValue))
		}
	}
	return currEff.(IOEffect)
}
