package types

import (
	"fmt"
	"github.com/mobilemindtec/go-io/option"
	"github.com/mobilemindtec/go-io/result"
	"github.com/mobilemindtec/go-io/util"
	"log"
	"reflect"
)

type IOCatchAll[A any] struct {
	value          *result.Result[*option.Option[A]]
	prevEffect     IOEffect
	fnResult       func(error) *result.Result[A]
	fnResultOption func(error) *result.Result[*option.Option[A]]
	fnOption       func(error) *option.Option[A]
	debug          bool
}

func NewCatchAll[A any](f func(error) *result.Result[*option.Option[A]]) *IOCatchAll[A] {
	return &IOCatchAll[A]{fnResultOption: f}
}

func NewCatchAllOfResult[A any](f func(error) *result.Result[A]) *IOCatchAll[A] {
	return &IOCatchAll[A]{fnResult: f}
}

func NewCatchAllOfOption[A any](f func(error) *option.Option[A]) *IOCatchAll[A] {
	return &IOCatchAll[A]{fnOption: f}
}

func (this *IOCatchAll[T]) TypeIn() reflect.Type {
	return reflect.TypeFor[T]()
}
func (this *IOCatchAll[T]) TypeOut() reflect.Type {
	return reflect.TypeFor[T]()
}

func (this *IOCatchAll[A]) SetDebug(b bool) {
	this.debug = b
}

func (this *IOCatchAll[A]) String() string {
	return fmt.Sprintf("CatchAll(%v)", this.value.String())
}

func (this *IOCatchAll[A]) SetPrevEffect(prev IOEffect) {
	this.prevEffect = prev
}

func (this *IOCatchAll[A]) GetPrevEffect() *option.Option[IOEffect] {
	return option.Of(this.prevEffect)
}

func (this *IOCatchAll[A]) GetResult() ResultOptionAny {
	return this.value.ToResultOfOption()
}

func (this *IOCatchAll[A]) UnsafeRun() IOEffect {
	var currEff interface{} = this
	prevEff := this.GetPrevEffect()
	this.value = result.OfValue(option.None[A]())

	if prevEff.NonEmpty() {
		r := prevEff.Get().GetResult()
		if r.IsError() {
			if this.fnResult != nil {
				res := this.fnResult(r.Failure())
				if res.HasError() {
					this.value = result.OfError[*option.Option[A]](res.Failure())
				} else {
					this.value = result.OfValue(option.Some(res.Get()))
				}
			} else if this.fnResultOption != nil {
				this.value = this.fnResultOption(r.Failure())
			} else {
				this.value = result.OfValue(this.fnOption(r.Failure()))
			}
		} else if r.Get().NonEmpty() {
			val := r.Get().GetValue()
			if effValue, ok := val.(A); ok {
				this.value = result.OfValue(option.Some(effValue))
			} else {
				util.PanicCastType("IOCatchAll",
					reflect.TypeOf(val), reflect.TypeFor[A]())

			}
		}
	}

	if this.debug {
		log.Printf("%v\n", this.String())

	}

	return currEff.(IOEffect)
}
