package ios

import (
	"fmt"
	"github.com/mobilemindtec/go-io/option"
	"github.com/mobilemindtec/go-io/result"
	"github.com/mobilemindtec/go-io/types"
	"github.com/mobilemindtec/go-io/util"
	"log"
	"reflect"
)

type IORecover[A any] struct {
	value         *result.Result[*option.Option[A]]
	prevEffect    types.IOEffect
	fpure         func(error) A
	fresult       func(error) *result.Result[A]
	foption       func(error) *option.Option[A]
	fresultOption func(error) *result.Result[*option.Option[A]]
	debug         bool
	debugInfo     *types.IODebugInfo
}

func NewRecoverPure[A any](f func(error) A) *IORecover[A] {
	return &IORecover[A]{fpure: f}
}

func NewRecover[A any](f func(error) *result.Result[A]) *IORecover[A] {
	return &IORecover[A]{fresult: f}
}

func NewRecoverOption[A any](f func(error) *option.Option[A]) *IORecover[A] {
	return &IORecover[A]{foption: f}
}

func NewRecoverResultOption[A any](f func(error) *result.Result[*option.Option[A]]) *IORecover[A] {
	return &IORecover[A]{fresultOption: f}
}

func (this *IORecover[A]) Lift() *types.IO[A] {
	return types.NewIO[A]().Effects(this)
}

func (this *IORecover[A]) TypeIn() reflect.Type {
	return reflect.TypeFor[A]()
}

func (this *IORecover[A]) TypeOut() reflect.Type {
	return reflect.TypeFor[A]()
}

func (this *IORecover[A]) SetDebug(b bool) {
	this.debug = b
}

func (this *IORecover[A]) SetDebugInfo(info *types.IODebugInfo) {
	this.debugInfo = info
}

func (this *IORecover[A]) GetDebugInfo() *types.IODebugInfo {
	return this.debugInfo
}

func (this *IORecover[A]) String() string {
	return fmt.Sprintf("Recover(%v)", this.value.String())
}

func (this *IORecover[A]) SetPrevEffect(prev types.IOEffect) {
	this.prevEffect = prev
}

func (this *IORecover[A]) GetPrevEffect() *option.Option[types.IOEffect] {
	return option.Of(this.prevEffect)
}

func (this *IORecover[A]) GetResult() types.ResultOptionAny {
	return this.value.ToResultOfOption()
}

func (this *IORecover[A]) UnsafeRun() types.IOEffect {
	var currEff interface{} = this
	prevEff := this.GetPrevEffect()
	this.value = result.OfValue(option.None[A]())

	if prevEff.NonEmpty() {
		r := prevEff.Get().GetResult()
		if r.IsError() {

			if this.fpure != nil {
				this.value = result.OfValue(option.Of(this.fpure(r.Failure())))
			} else if this.fresultOption != nil {
				this.value = this.fresultOption(r.Failure())
			} else if this.foption != nil {
				this.value = result.OfValue(this.foption(r.Failure()))
			} else if this.fresult != nil {
				this.fresult(r.Failure()).
					IfOk(func(a A) {
						this.value = result.OfValue(option.Of(a))
					}).IfError(func(err error) {
					this.value = result.OfError[*option.Option[A]](err)
				})
			}
		} else if r.Get().NonEmpty() {
			val := r.Get().GetValue()
			if effValue, ok := val.(A); ok {

				this.value = result.OfValue(option.Some(effValue))
			} else {
				util.PanicCastType("IORecover",
					reflect.TypeOf(val), reflect.TypeFor[A]())

			}
		}
	}

	if this.debug {
		log.Printf("%v\n", this.String())

	}

	return currEff.(types.IOEffect)
}
