package ios

import (
	"fmt"
	"github.com/mobilemindtech/go-io/option"
	"github.com/mobilemindtech/go-io/result"
	"github.com/mobilemindtech/go-io/types"
	"github.com/mobilemindtech/go-io/types/unit"
	"log"
	"reflect"
)

type IOPure[T any] struct {
	value      *option.Option[T]
	prevEffect types.IOEffect
	f          func() T
	debug      bool
	debugInfo  *types.IODebugInfo
}

func NewPureValue[T any](value T) *IOPure[T] {
	return &IOPure[T]{value: option.Of(value)}
}

func NewPure[T any](f func() T) *IOPure[T] {
	return &IOPure[T]{f: f}
}

func (this *IOPure[T]) TypeIn() reflect.Type {
	return reflect.TypeFor[*unit.Unit]()
}

func (this *IOPure[T]) TypeOut() reflect.Type {
	return reflect.TypeFor[T]()
}

func (this *IOPure[T]) Lift() *types.IO[T] {
	return types.NewIO[T]().Effects(this)
}

func (this *IOPure[T]) SetDebug(b bool) {
	this.debug = b
}

func (this *IOPure[A]) SetDebugInfo(info *types.IODebugInfo) {
	this.debugInfo = info
}

func (this *IOPure[A]) GetDebugInfo() *types.IODebugInfo {
	return this.debugInfo
}

func (this *IOPure[T]) String() string {
	return fmt.Sprintf("Pure(%v)", this.value.String())
}

func (this *IOPure[T]) SetPrevEffect(prev types.IOEffect) {
	this.prevEffect = prev
}

func (this *IOPure[T]) GetPrevEffect() *option.Option[types.IOEffect] {
	return option.Of(this.prevEffect)
}

func (this *IOPure[T]) GetResult() types.ResultOptionAny {

	if this.value.Empty() {
		return result.OfValue(option.None[any]())
	} else {
		return result.OfValue(option.Some[any](this.value.Get()))
	}
}

func (this *IOPure[T]) UnsafeRun() types.IOEffect {
	var eff interface{} = this

	if this.f != nil {
		this.value = option.Of(this.f())
	}

	if this.debug {
		log.Printf("%v\n", this.String())
	}

	return eff.(types.IOEffect)
}
