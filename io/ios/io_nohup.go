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

type IONohup[A any] struct {
	value      *result.Result[*option.Option[A]]
	prevEffect types.IOEffect
	debug      bool
	debugInfo  *types.IODebugInfo
}

func NewNohup[A any]() *IONohup[A] {
	return &IONohup[A]{}
}

func (this *IONohup[T]) String() string {
	return fmt.Sprintf("Nohup(%v)", this.value.String())
}

func (this *IONohup[T]) TypeIn() reflect.Type {
	return reflect.TypeFor[*unit.Unit]()
}

func (this *IONohup[T]) TypeOut() reflect.Type {
	return reflect.TypeFor[T]()
}

func (this *IONohup[T]) SetDebug(b bool) {
	this.debug = b
}

func (this *IONohup[T]) SetDebugInfo(info *types.IODebugInfo) {
	this.debugInfo = info
}

func (this *IONohup[T]) GetDebugInfo() *types.IODebugInfo {
	return this.debugInfo
}

func (this *IONohup[T]) SetPrevEffect(prev types.IOEffect) {
	this.prevEffect = prev
}

func (this *IONohup[T]) Lift() *types.IO[T] {
	return types.NewIO[T]().Effects(this)
}

func (this *IONohup[T]) GetPrevEffect() *option.Option[types.IOEffect] {
	return option.Of(this.prevEffect)
}

func (this *IONohup[T]) GetResult() types.ResultOptionAny {
	return this.value.ToResultOfOption()
}

func (this *IONohup[T]) UnsafeRun() types.IOEffect {
	var currEff interface{} = this
	prevEff := this.GetPrevEffect()
	this.value = result.OfValue(option.None[T]())
	execute := true

	if prevEff.NonEmpty() {
		r := prevEff.Get().GetResult()
		if r.IsError() {
			this.value = result.OfError[*option.Option[T]](r.Failure())
			execute = false
		}
	}

	if execute {
		this.value = result.OfValue(option.Some(util.NewOf[T]()))
	}

	if this.debug {
		log.Printf("%v\n", this.String())
	}

	return currEff.(types.IOEffect)
}
