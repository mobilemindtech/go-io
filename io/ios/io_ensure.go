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

type IOEnsure[A any] struct {
	value      *result.Result[*option.Option[A]]
	prevEffect types.IOEffect
	f          func()
	debug      bool
	debugInfo  *types.IODebugInfo
}

func NewEnsure[A any](f func()) *IOEnsure[A] {
	return &IOEnsure[A]{f: f}
}

func (this *IOEnsure[A]) Lift() *types.IO[A] {
	return types.NewIO[A]().Effects(this)
}

func (this *IOEnsure[A]) TypeIn() reflect.Type {
	return reflect.TypeFor[A]()
}

func (this *IOEnsure[A]) TypeOut() reflect.Type {
	return reflect.TypeFor[A]()
}

func (this *IOEnsure[A]) SetDebug(b bool) {
	this.debug = b
}

func (this *IOEnsure[A]) SetDebugInfo(info *types.IODebugInfo) {
	this.debugInfo = info
}

func (this *IOEnsure[A]) GetDebugInfo() *types.IODebugInfo {
	return this.debugInfo
}

func (this *IOEnsure[A]) String() string {
	return fmt.Sprintf("Ensure(%v)", this.value.String())
}

func (this *IOEnsure[A]) SetPrevEffect(prev types.IOEffect) {
	this.prevEffect = prev
}

func (this *IOEnsure[A]) GetPrevEffect() *option.Option[types.IOEffect] {
	return option.Of(this.prevEffect)
}

func (this *IOEnsure[A]) GetResult() types.ResultOptionAny {
	return this.value.ToResultOfOption()
}

func (this *IOEnsure[A]) UnsafeRun() types.IOEffect {
	var currEff interface{} = this
	prevEff := this.GetPrevEffect()
	this.value = result.OfValue(option.None[A]())

	this.f()

	if prevEff.NonEmpty() {
		r := prevEff.Get().GetResult()
		if r.IsError() {
			this.value = result.OfError[*option.Option[A]](r.Failure())
		} else if r.Get().NonEmpty() {
			val := r.Get().GetValue()
			if effValue, ok := val.(A); ok {
				this.value = result.OfValue(option.Some(effValue))
			} else {
				util.PanicCastType("IODebug",
					reflect.TypeOf(val), reflect.TypeFor[A]())
			}
		}
	}

	if this.debug {
		log.Printf("%v\n", this.String())
	}

	return currEff.(types.IOEffect)
}
