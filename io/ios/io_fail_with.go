package ios

import (
	"fmt"
	"github.com/mobilemindtec/go-io/option"
	"github.com/mobilemindtec/go-io/result"
	"github.com/mobilemindtec/go-io/types"
	"log"
	"reflect"
)

type IOFailWith[A any] struct {
	value      *result.Result[*option.Option[A]]
	prevEffect types.IOEffect
	f          func() error
	debug      bool
	debugInfo  *types.IODebugInfo
}

func NewFailWith[A any](f func() error) *IOFailWith[A] {
	return &IOFailWith[A]{f: f}
}

func (this *IOFailWith[A]) Lift() *types.IO[A] {
	return types.NewIO[A]().Effects(this)
}

func (this *IOFailWith[T]) TypeIn() reflect.Type {
	return reflect.TypeFor[T]()
}

func (this *IOFailWith[T]) TypeOut() reflect.Type {
	return reflect.TypeFor[T]()
}

func (this *IOFailWith[A]) SetDebug(b bool) {
	this.debug = b
}

func (this *IOFailWith[T]) SetDebugInfo(info *types.IODebugInfo) {
	this.debugInfo = info
}

func (this *IOFailWith[T]) GetDebugInfo() *types.IODebugInfo {
	return this.debugInfo
}

func (this *IOFailWith[A]) String() string {
	return fmt.Sprintf("FailWith(%v)", this.value.String())
}

func (this *IOFailWith[A]) SetPrevEffect(prev types.IOEffect) {
	this.prevEffect = prev
}

func (this *IOFailWith[A]) GetPrevEffect() *option.Option[types.IOEffect] {
	return option.Of(this.prevEffect)
}

func (this *IOFailWith[A]) GetResult() types.ResultOptionAny {
	return this.value.ToResultOfOption()
}

func (this *IOFailWith[A]) UnsafeRun() types.IOEffect {
	var currEff interface{} = this
	prevEff := this.GetPrevEffect()
	this.value = result.OfValue(option.None[A]())

	if prevEff.NonEmpty() {
		r := prevEff.Get().GetResult()

		if r.IsError() {
			this.value = result.OfError[*option.Option[A]](r.Failure())
		} else {
			this.value = result.OfError[*option.Option[A]](this.f())
		}
	}

	if this.debug {
		log.Printf("%v\n", this.String())
	}

	return currEff.(types.IOEffect)
}
