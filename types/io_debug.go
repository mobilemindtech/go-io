package types

import (
	"fmt"
	"github.com/mobilemindtec/go-io/option"
	"github.com/mobilemindtec/go-io/result"
)

type IODebug[A any] struct {
	valueA     *result.Result[A]
	prevEffect IOEffect
	label      string
}

func NewDebug[A any](label string) *IODebug[A] {
	return &IODebug[A]{label: label}
}

func (this *IODebug[A]) SetPrevEffect(prev IOEffect) {
	this.prevEffect = prev
}

func (this *IODebug[A]) GetPrevEffect() *option.Option[IOEffect] {
	return option.Of(this.prevEffect)
}

func (this *IODebug[A]) GetResult() *result.Result[any] {
	return this.valueA.ToResultOfAny()
}

func (this *IODebug[A]) UnsafeRun() IOEffect {
	var currEff interface{} = this
	prevEff := this.GetPrevEffect()
	if prevEff.NonEmpty() {
		r := prevEff.Get().GetResult()
		fmt.Println(fmt.Sprintf("<DEBUG>: %v - %v", this.label, prevEff.Get()))
		if r.IsError() {
			this.valueA = result.OfError[A](r.Error())
		} else if r.OptionNonEmpty() {
			effValue := prevEff.Get().GetResult().ToOption().OrNil().(A)
			this.valueA = result.OfValue(effValue)
		}
	} else {
		fmt.Println(fmt.Sprintf("<DEBUG>: %v - IO(empty)", this.label))
	}
	return currEff.(IOEffect)
}
