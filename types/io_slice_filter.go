package types

import (
	"github.com/mobilemindtec/go-io/option"
	"github.com/mobilemindtec/go-io/result"
)

type IOSliceFilter[A any] struct {
	valueA     *result.Result[[]A]
	prevEffect IOEffect
	f          func(A) bool
}

func NewSliceFilter[A any](f func(A) bool) *IOSliceFilter[A] {
	return &IOSliceFilter[A]{f: f}
}

func (this *IOSliceFilter[A]) SetPrevEffect(prev IOEffect) {
	this.prevEffect = prev
}

func (this *IOSliceFilter[A]) GetPrevEffect() *option.Option[IOEffect] {
	return option.Of(this.prevEffect)
}

func (this *IOSliceFilter[A]) GetResult() *result.Result[any] {
	return this.valueA.ToResultOfAny()
}

func (this *IOSliceFilter[A]) UnsafeRun() IOEffect {
	var currEff interface{} = this
	prevEff := this.GetPrevEffect()
	if prevEff.NonEmpty() {
		r := prevEff.Get().GetResult()
		if r.IsError() {
			this.valueA = result.OfError[[]A](r.Error())
		} else if r.OptionNonEmpty() {
			effValue := prevEff.Get().GetResult().ToOption().OrNil().([]A)
			var list []A
			for _, it := range effValue {
				if this.f(it) {
					list = append(list, it)
				}
			}
			if len(list) > 0 {
				this.valueA = result.OfValue(list)
			} else {
				this.valueA = result.OfNil[[]A]()
			}
		}
	}
	return currEff.(IOEffect)
}
