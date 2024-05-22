package types

import (
	"github.com/mobilemindtec/go-io/option"
	"github.com/mobilemindtec/go-io/result"
)

type IOAttempt[A any] struct {
	valueA     *result.Result[A]
	prevEffect IOEffect
	fr         func() *result.Result[A]
	ft         func() (A, error)
}

func NewAttempt[A any](f func() *result.Result[A]) *IOAttempt[A] {
	return &IOAttempt[A]{fr: f}
}

func NewAttemptTry[A any](f func() (A, error)) *IOAttempt[A] {
	return &IOAttempt[A]{ft: f}
}

func (this *IOAttempt[A]) SetPrevEffect(prev IOEffect) {
	this.prevEffect = prev
}

func (this *IOAttempt[A]) GetPrevEffect() *option.Option[IOEffect] {
	return option.Of(this.prevEffect)
}

func (this *IOAttempt[A]) GetResult() *result.Result[any] {
	return this.valueA.ToResultOfAny()
}

func (this *IOAttempt[A]) UnsafeRun() IOEffect {
	var currEff interface{} = this
	if this.fr != nil {
		this.valueA = this.fr()
	} else {
		this.valueA = result.Try(this.ft)
	}
	return currEff.(IOEffect)
}
