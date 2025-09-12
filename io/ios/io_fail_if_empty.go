package ios

import (
	"fmt"
	"github.com/mobilemindtech/go-io/option"
	"github.com/mobilemindtech/go-io/result"
	"github.com/mobilemindtech/go-io/types"
	"github.com/mobilemindtech/go-io/types/unit"
	"github.com/mobilemindtech/go-io/util"
	"log"
	"reflect"
)

type IOFailIfEmpty[A any] struct {
	value      *result.Result[*option.Option[A]]
	prevEffect types.IOEffect
	f          func() error
	unit       bool
	debug      bool
	debugInfo  *types.IODebugInfo
}

func NewFailIfEmpty[A any](f func() error) *IOFailIfEmpty[A] {
	return &IOFailIfEmpty[A]{f: f}
}

func NewFailIfEmptyUnit(f func() error) *IOFailIfEmpty[*unit.Unit] {
	return &IOFailIfEmpty[*unit.Unit]{f: f, unit: true}
}

func (this *IOFailIfEmpty[A]) Lift() *types.IO[A] {
	return types.NewIO[A]().Effects(this)
}

func (this *IOFailIfEmpty[T]) TypeIn() reflect.Type {
	return reflect.TypeFor[T]()
}

func (this *IOFailIfEmpty[T]) TypeOut() reflect.Type {
	return reflect.TypeFor[T]()
}

func (this *IOFailIfEmpty[A]) SetDebug(b bool) {
	this.debug = b
}

func (this *IOFailIfEmpty[T]) SetDebugInfo(info *types.IODebugInfo) {
	this.debugInfo = info
}

func (this *IOFailIfEmpty[T]) GetDebugInfo() *types.IODebugInfo {
	return this.debugInfo
}

func (this *IOFailIfEmpty[A]) String() string {
	return fmt.Sprintf("FailIfEmpty(%v)", this.value.String())
}

func (this *IOFailIfEmpty[A]) SetPrevEffect(prev types.IOEffect) {
	this.prevEffect = prev
}

func (this *IOFailIfEmpty[A]) GetPrevEffect() *option.Option[types.IOEffect] {
	return option.Of(this.prevEffect)
}

func (this *IOFailIfEmpty[A]) GetResult() types.ResultOptionAny {
	return this.value.ToResultOfOption()
}

func (this *IOFailIfEmpty[A]) UnsafeRun() types.IOEffect {
	var currEff interface{} = this
	prevEff := this.GetPrevEffect()
	this.value = result.OfValue(option.None[A]())

	if prevEff.NonEmpty() {
		r := prevEff.Get().GetResult()

		if r.IsError() {
			this.value = result.OfError[*option.Option[A]](r.Failure())
		} else if r.Get().NonEmpty() {
			if !this.unit {
				val := r.Get().GetValue()
				if effValue, ok := val.(A); ok {
					this.value = result.OfValue(option.Some(effValue))
				} else {
					util.PanicCastType("IOFailIfEmpty",
						reflect.TypeOf(val), reflect.TypeFor[A]())
				}
			} else {
				var effValue interface{} = unit.OfUnit()
				this.value = result.OfValue(option.Some(effValue.(A)))
			}
		} else {
			this.value = result.OfError[*option.Option[A]](this.f())
		}
	}

	if this.debug {
		log.Printf("%v\n", this.String())
	}

	return currEff.(types.IOEffect)
}
