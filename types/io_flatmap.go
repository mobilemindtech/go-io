package types

import (
	"github.com/mobilemindtec/go-io/option"
	"github.com/mobilemindtec/go-io/result"
)

type IOFlatMap[A any, B any] struct {
	//valueA     *option.Option[A]
	valueB     *result.Result[B]
	prevEffect IOEffect
	f          func(A) *IO[B]
}

func NewFlatMap[A any, B any](f func(A) *IO[B]) *IOFlatMap[A, B] {
	return &IOFlatMap[A, B]{f: f}
}

func (this *IOFlatMap[A, B]) SetPrevEffect(prev IOEffect) {
	this.prevEffect = prev
}

func (this *IOFlatMap[A, B]) GetPrevEffect() *option.Option[IOEffect] {
	return option.Of(this.prevEffect)
}

func (this *IOFlatMap[A, B]) GetResult() *result.Result[any] {
	return this.valueB.ToResultOfAny()
}

func (this *IOFlatMap[A, B]) UnsafeRun() IOEffect {
	var currEff interface{} = this
	prevEff := this.GetPrevEffect()
	if prevEff.NonEmpty() {
		effResult := prevEff.Get().GetResult()
		if effResult.IsError() {
			this.valueB = result.OfError[B](effResult.Error())
		} else if effResult.OptionNonEmpty() {
			effValue := prevEff.Get().GetResult().ToOption().OrNil().(A)
			this.valueB = this.f(effValue).UnsafeRun()
		}
	}
	return currEff.(IOEffect)
}
