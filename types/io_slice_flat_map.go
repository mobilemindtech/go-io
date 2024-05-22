package types

import (
	"github.com/mobilemindtec/go-io/option"
	"github.com/mobilemindtec/go-io/result"
)

type IOSliceFlatMap[A any, B any] struct {
	//valueA     *result.Result[[]A]
	valueB     *result.Result[[]B]
	prevEffect IOEffect
	f          func(A) *IO[B]
}

func NewSliceFlatMap[A any, B any](f func(A) *IO[B]) *IOSliceFlatMap[A, B] {
	return &IOSliceFlatMap[A, B]{f: f}
}

func (this *IOSliceFlatMap[A, B]) SetPrevEffect(prev IOEffect) {
	this.prevEffect = prev
}

func (this *IOSliceFlatMap[A, B]) GetPrevEffect() *option.Option[IOEffect] {
	return option.Of(this.prevEffect)
}

func (this *IOSliceFlatMap[A, B]) GetResult() *result.Result[any] {
	return this.valueB.ToResultOfAny()
}

func (this *IOSliceFlatMap[A, B]) UnsafeRun() IOEffect {
	var currEff interface{} = this
	prevEff := this.GetPrevEffect()
	if prevEff.NonEmpty() {
		r := prevEff.Get().GetResult()
		if r.IsError() {
			this.valueB = result.OfError[[]B](r.Error())
		} else if r.OptionNonEmpty() {
			effValue := prevEff.Get().GetResult().ToOption().OrNil().([]A)
			var list []B
			for _, it := range effValue {
				resultIO := this.f(it).UnsafeRun()
				if resultIO.IsError() || resultIO.OptionEmpty() {
					this.valueB = result.OfError[[]B](resultIO.Error())
					break
				}
				list = append(list, resultIO.Get())
			}
			this.valueB = result.OfValue(list)
		}
	}
	return currEff.(IOEffect)
}
