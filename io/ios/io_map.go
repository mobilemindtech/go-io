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

type IOMap[A any, B any] struct {
	value      *result.Result[*option.Option[B]]
	prevEffect types.IOEffect
	f          func(A) B
	debug      bool
	debugInfo  *types.IODebugInfo
}

func NewMap[A any, B any](f func(A) B) *IOMap[A, B] {
	return &IOMap[A, B]{f: f}
}

func (this *IOMap[A, B]) Lift() *types.IO[B] {
	return types.NewIO[B]().Effects(this)
}

func (this *IOMap[A, B]) TypeIn() reflect.Type {
	return reflect.TypeFor[A]()
}

func (this *IOMap[A, B]) TypeOut() reflect.Type {
	return reflect.TypeFor[B]()
}

func (this *IOMap[A, B]) SetDebug(b bool) {
	this.debug = b
}

func (this *IOMap[A, B]) SetDebugInfo(info *types.IODebugInfo) {
	this.debugInfo = info
}

func (this *IOMap[A, B]) GetDebugInfo() *types.IODebugInfo {
	return this.debugInfo
}

func (this *IOMap[A, B]) String() string {
	return fmt.Sprintf("Map(%v)", this.value.String())
}

func (this *IOMap[A, B]) SetPrevEffect(prev types.IOEffect) {
	this.prevEffect = prev
}

func (this *IOMap[A, B]) GetPrevEffect() *option.Option[types.IOEffect] {
	return option.Of(this.prevEffect)
}

func (this *IOMap[A, B]) GetResult() types.ResultOptionAny {
	return this.value.ToResultOfOption()
}

func (this *IOMap[A, B]) UnsafeRun() types.IOEffect {
	var currEff interface{} = this
	prevEff := this.GetPrevEffect()
	this.value = result.OfValue(option.None[B]())

	if prevEff.NonEmpty() {
		r := prevEff.Get().GetResult()

		if r.IsError() {
			this.value = result.OfError[*option.Option[B]](r.GetError())
		} else if r.Get().NonEmpty() {
			val := r.Get().GetValue()
			if effValue, ok := val.(A); ok {
				this.value = result.OfValue(option.Some(this.f(effValue)))
			} else {
				util.PanicCastType("IOMap",
					reflect.TypeOf(val), reflect.TypeFor[B]())
			}
		}
	}

	if this.debug {
		log.Printf("%v\n", this.String())
	}

	return currEff.(types.IOEffect)
}
