package types

import (
	"github.com/mobilemindtec/go-io/option"
	"github.com/mobilemindtec/go-io/result"
)

type IOFailWith[A any] struct {
	valueA     *result.Result[A]
	prevEffect IOEffect
	fe         func(A) error
	fr         func(A) *result.Result[A]
}

func NewFailWith[A any](f func(A) *result.Result[A]) *IOFailWith[A] {
	return &IOFailWith[A]{fr: f}
}
func NewFailWithError[A any](f func(A) error) *IOFailWith[A] {
	return &IOFailWith[A]{fe: f}
}

func (this *IOFailWith[A]) SetPrevEffect(prev IOEffect) {
	this.prevEffect = prev
}

func (this *IOFailWith[A]) GetPrevEffect() *option.Option[IOEffect] {
	return option.Of(this.prevEffect)
}

func (this *IOFailWith[A]) GetResult() *result.Result[any] {
	return this.valueA.ToResultOfAny()
}

func (this *IOFailWith[A]) UnsafeRun() IOEffect {
	var currEff interface{} = this
	prevEff := this.GetPrevEffect()
	if prevEff.NonEmpty() {
		r := prevEff.Get().GetResult()
		if r.IsError() {
			this.valueA = result.OfError[A](r.GetError())
		} else if r.OptionNonEmpty() {
			effValue := r.ToOption().OrNil().(A)
			if this.fr != nil {
				this.valueA = this.fr(effValue)
			} else {
				this.valueA = result.Make(effValue, this.fe(effValue))
			}
		}
	}
	return currEff.(IOEffect)
}
