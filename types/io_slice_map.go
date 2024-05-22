package types

import (
	"github.com/mobilemindtec/go-io/option"
	"github.com/mobilemindtec/go-io/result"
)

type IOSliceMap[A any, B any] struct {
	//valueA     *result.Result[[]A]
	valueB     *result.Result[[]B]
	prevEffect IOEffect
	f          func(A) B
}

func NewSliceMap[A any, B any](f func(A) B) *IOSliceMap[A, B] {
	return &IOSliceMap[A, B]{f: f}
}

func (this *IOSliceMap[A, B]) SetPrevEffect(prev IOEffect) {
	this.prevEffect = prev
}

func (this *IOSliceMap[A, B]) GetPrevEffect() *option.Option[IOEffect] {
	return option.Of(this.prevEffect)
}

func (this *IOSliceMap[A, B]) GetResult() *result.Result[any] {
	return this.valueB.ToResultOfAny()
}

func (this *IOSliceMap[A, B]) UnsafeRun() IOEffect {
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
				list = append(list, this.f(it))
			}
			this.valueB = result.OfValue(list)
		}
	}
	return currEff.(IOEffect)
}
