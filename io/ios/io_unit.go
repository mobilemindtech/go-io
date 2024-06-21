package ios

import (
	"fmt"
	"github.com/mobilemindtec/go-io/option"
	"github.com/mobilemindtec/go-io/result"
	"github.com/mobilemindtec/go-io/types"
	"log"
	"reflect"
)

type IOUnit struct {
	value      *result.Result[*option.Option[*types.Unit]]
	prevEffect types.IOEffect
	debug      bool
	debugInfo  *types.IODebugInfo
}

func NewUnit() *IOUnit {
	return &IOUnit{}
}

func (this *IOUnit) String() string {
	return fmt.Sprintf("Unit(%v)", this.value.String())
}

func (this *IOUnit) TypeIn() reflect.Type {
	return reflect.TypeFor[*types.Unit]()
}

func (this *IOUnit) TypeOut() reflect.Type {
	return reflect.TypeFor[*types.Unit]()
}

func (this *IOUnit) SetDebug(b bool) {
	this.debug = b
}

func (this *IOUnit) SetDebugInfo(info *types.IODebugInfo) {
	this.debugInfo = info
}

func (this *IOUnit) GetDebugInfo() *types.IODebugInfo {
	return this.debugInfo
}

func (this *IOUnit) SetPrevEffect(prev types.IOEffect) {
	this.prevEffect = prev
}

func (this *IOUnit) Lift() *types.IO[*types.Unit] {
	return types.NewIO[*types.Unit]().Effects(this)
}

func (this *IOUnit) GetPrevEffect() *option.Option[types.IOEffect] {
	return option.Of(this.prevEffect)
}

func (this *IOUnit) GetResult() types.ResultOptionAny {
	return this.value.ToResultOfOption()
}

func (this *IOUnit) UnsafeRun() types.IOEffect {
	var currEff interface{} = this
	prevEff := this.GetPrevEffect()
	this.value = result.OfValue(option.None[*types.Unit]())

	execute := true

	if prevEff.NonEmpty() {
		r := prevEff.Get().GetResult()
		if r.IsError() {
			this.value = result.OfError[*option.Option[*types.Unit]](r.Failure())
			execute = false
		}
	}

	if execute {
		this.value = result.OfValue(option.Some(types.OfUnit()))
	}

	if this.debug {
		log.Printf("%v\n", this.String())
	}

	return currEff.(types.IOEffect)
}
