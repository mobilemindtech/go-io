package ios

import (
	"fmt"
	"github.com/mobilemindtec/go-io/option"
	"github.com/mobilemindtec/go-io/result"
	"github.com/mobilemindtec/go-io/state"
	"github.com/mobilemindtec/go-io/types"
	"log"
	"reflect"
)

type IOLoadVar[A any] struct {
	value      *result.Result[*option.Option[A]]
	prevEffect types.IOEffect
	debug      bool
	debugInfo  *types.IODebugInfo
	state      *state.State
}

func NewLoadVar[A any]() *IOLoadVar[A] {
	return &IOLoadVar[A]{}
}

func (this *IOLoadVar[A]) String() string {
	return fmt.Sprintf("LoadVar(%v)", this.value.String())
}

func (this *IOLoadVar[T]) TypeIn() reflect.Type {
	return reflect.TypeFor[*types.Unit]()
}

func (this *IOLoadVar[T]) TypeOut() reflect.Type {
	return reflect.TypeFor[T]()
}

func (this *IOLoadVar[A]) SetState(st *state.State) {
	this.state = st
}

func (this *IOLoadVar[A]) SetDebug(b bool) {
	this.debug = b
}

func (this *IOLoadVar[T]) SetDebugInfo(info *types.IODebugInfo) {
	this.debugInfo = info
}

func (this *IOLoadVar[T]) GetDebugInfo() *types.IODebugInfo {
	return this.debugInfo
}

func (this *IOLoadVar[A]) SetPrevEffect(prev types.IOEffect) {
	this.prevEffect = prev
}

func (this *IOLoadVar[A]) Lift() *types.IO[A] {
	return types.NewIO[A]().Effects(this)
}

func (this *IOLoadVar[A]) GetPrevEffect() *option.Option[types.IOEffect] {
	return option.Of(this.prevEffect)
}

func (this *IOLoadVar[A]) GetResult() types.ResultOptionAny {
	return this.value.ToResultOfOption()
}

func (this *IOLoadVar[A]) UnsafeRun() types.IOEffect {
	var currEff interface{} = this
	prevEff := this.GetPrevEffect()
	this.value = result.OfValue(option.None[A]())

	execute := true
	if prevEff.NonEmpty() {
		r := prevEff.Get().GetResult()
		if r.IsError() {
			this.value = result.OfError[*option.Option[A]](r.Failure())
			execute = false
		}
	}

	if execute {
		this.value = result.OfValue(option.Some(state.Var[A](this.state)))
	}

	if this.debug {
		log.Printf("%v\n", this.String())
	}

	return currEff.(types.IOEffect)
}
