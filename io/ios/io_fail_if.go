package ios

import (
	"fmt"
	"github.com/mobilemindtech/go-io/option"
	"github.com/mobilemindtech/go-io/result"
	"github.com/mobilemindtech/go-io/types"
	"github.com/mobilemindtech/go-io/util"
	"log"
	"reflect"
)

type IOFailIf[A any] struct {
	value      *result.Result[*option.Option[A]]
	prevEffect types.IOEffect
	f          func(A) error
	debug      bool
	debugInfo  *types.IODebugInfo
}

func NewFailIf[A any](f func(A) error) *IOFailIf[A] {
	return &IOFailIf[A]{f: f}
}

func (this *IOFailIf[A]) Lift() *types.IO[A] {
	return types.NewIO[A]().Effects(this)
}

func (this *IOFailIf[T]) TypeIn() reflect.Type {
	return reflect.TypeFor[T]()
}

func (this *IOFailIf[T]) TypeOut() reflect.Type {
	return reflect.TypeFor[T]()
}

func (this *IOFailIf[A]) SetDebug(b bool) {
	this.debug = b
}

func (this *IOFailIf[T]) SetDebugInfo(info *types.IODebugInfo) {
	this.debugInfo = info
}

func (this *IOFailIf[T]) GetDebugInfo() *types.IODebugInfo {
	return this.debugInfo
}

func (this *IOFailIf[A]) String() string {
	return fmt.Sprintf("FailIf(%v)", this.value.String())
}

func (this *IOFailIf[A]) SetPrevEffect(prev types.IOEffect) {
	this.prevEffect = prev
}

func (this *IOFailIf[A]) GetPrevEffect() *option.Option[types.IOEffect] {
	return option.Of(this.prevEffect)
}

func (this *IOFailIf[A]) GetResult() types.ResultOptionAny {
	return this.value.ToResultOfOption()
}

func (this *IOFailIf[A]) UnsafeRun() types.IOEffect {
	var currEff interface{} = this
	prevEff := this.GetPrevEffect()
	this.value = result.OfValue(option.None[A]())

	if prevEff.NonEmpty() {
		r := prevEff.Get().GetResult()

		if r.IsError() {
			this.value = result.OfError[*option.Option[A]](r.Failure())
		} else if r.Get().NonEmpty() {
			val := r.Get().GetValue()
			if effValue, ok := val.(A); ok {

				err := this.f(effValue)

				if err != nil {
					this.value = result.OfError[*option.Option[A]](err)
				} else {
					this.value = result.OfValue(option.Some(effValue))
				}

			} else {
				util.PanicCastType("FailIf",
					reflect.TypeOf(val), reflect.TypeFor[A]())
			}
		}
	}

	if this.debug {
		log.Printf("%v\n", this.String())
	}

	return currEff.(types.IOEffect)
}
