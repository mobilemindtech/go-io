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

type IOError[A any] struct {
	value      *result.Result[*option.Option[A]]
	prevEffect types.IOEffect
	debug      bool
	debugInfo  *types.IODebugInfo
}

func NewError[A any](err error) *IOError[A] {
	return &IOError[A]{
		value: result.OfErrorOption[A](err),
	}
}

func (this *IOError[A]) Lift() *types.IO[A] {
	return types.NewIO[A]().Effects(this)
}

func (this *IOError[A]) TypeIn() reflect.Type {
	return reflect.TypeFor[*unit.Unit]()
}

func (this *IOError[A]) TypeOut() reflect.Type {
	return reflect.TypeFor[*unit.Unit]()
}

func (this *IOError[A]) SetDebug(b bool) {
	this.debug = b
}

func (this *IOError[A]) SetDebugInfo(info *types.IODebugInfo) {
	this.debugInfo = info
}

func (this *IOError[A]) GetDebugInfo() *types.IODebugInfo {
	return this.debugInfo
}

func (this *IOError[A]) String() string {
	return fmt.Sprintf("Error(%v)", this.value.String())
}

func (this *IOError[A]) SetPrevEffect(prev types.IOEffect) {
	this.prevEffect = prev
}

func (this *IOError[A]) GetPrevEffect() *option.Option[types.IOEffect] {
	return option.Of(this.prevEffect)
}

func (this *IOError[A]) GetResult() types.ResultOptionAny {
	return this.value.ToResultOfOption()
}

func (this *IOError[A]) UnsafeRun() types.IOEffect {
	var currEff interface{} = this

	if this.debug {
		log.Printf("%v\n", this.String())
	}

	return currEff.(types.IOEffect)
}
