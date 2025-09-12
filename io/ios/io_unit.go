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

type IOUnit struct {
	value      *result.Result[*option.Option[*unit.Unit]]
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
	return reflect.TypeFor[*unit.Unit]()
}

func (this *IOUnit) TypeOut() reflect.Type {
	return reflect.TypeFor[*unit.Unit]()
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

func (this *IOUnit) Lift() *types.IO[*unit.Unit] {
	return types.NewIO[*unit.Unit]().Effects(this)
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
	this.value = result.OfValue(option.None[*unit.Unit]())

	execute := true

	if prevEff.NonEmpty() {
		r := prevEff.Get().GetResult()
		if r.IsError() {
			this.value = result.OfError[*option.Option[*unit.Unit]](r.Failure())
			execute = false
		}
	}

	if execute {
		this.value = result.OfValue(option.Some(unit.OfUnit()))
	}

	if this.debug {
		log.Printf("%v\n", this.String())
	}

	return currEff.(types.IOEffect)
}
