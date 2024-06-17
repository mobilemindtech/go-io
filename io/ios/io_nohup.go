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

type IONohup[A any] struct {
	value      *result.Result[*option.Option[A]]
	prevEffect types.IOEffect
	debug      bool
	debugInfo  *types.IODebugInfo
}

func NewNohup[A any]() *IONohup[A] {
	return &IONohup[A]{}
}

func (this *IONohup[A]) String() string {
	return fmt.Sprintf("Nohup(%v)", this.value.String())
}

func (this *IONohup[T]) TypeIn() reflect.Type {
	return reflect.TypeFor[T]()
}

func (this *IONohup[T]) TypeOut() reflect.Type {
	return reflect.TypeFor[T]()
}

func (this *IONohup[A]) SetDebug(b bool) {
	this.debug = b
}

func (this *IONohup[T]) SetDebugInfo(info *types.IODebugInfo) {
	this.debugInfo = info
}

func (this *IONohup[T]) GetDebugInfo() *types.IODebugInfo {
	return this.debugInfo
}

func (this *IONohup[A]) SetPrevEffect(prev types.IOEffect) {
	this.prevEffect = prev
}

func (this *IONohup[A]) Lift() *types.IO[A] {
	return types.NewIO[A]().Effects(this)
}

func (this *IONohup[A]) GetPrevEffect() *option.Option[types.IOEffect] {
	return option.Of(this.prevEffect)
}

func (this *IONohup[A]) GetResult() types.ResultOptionAny {
	return this.value.ToResultOfOption()
}

func (this *IONohup[A]) UnsafeRun() types.IOEffect {
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
				this.value = result.OfValue(option.Some(effValue))
			} else {
				util.PanicCastType("IONohup",
					reflect.TypeOf(val), reflect.TypeFor[A]())

			}
		}
	}

	if this.debug {
		log.Printf("%v\n", this.String())
	}

	return currEff.(types.IOEffect)
}
