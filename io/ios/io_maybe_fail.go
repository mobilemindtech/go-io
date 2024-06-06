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

type IOMaybeFail[A any] struct {
	value      *result.Result[*option.Option[A]]
	prevEffect types.IOEffect
	fe         func(A) error
	fr         func(A) *result.Result[A]
	debug      bool
	debugInfo  *types.IODebugInfo
}

func NewMaybeFail[A any](f func(A) *result.Result[A]) *IOMaybeFail[A] {
	return &IOMaybeFail[A]{fr: f}
}
func NewMaybeFailError[A any](f func(A) error) *IOMaybeFail[A] {
	return &IOMaybeFail[A]{fe: f}
}

func (this *IOMaybeFail[A]) TypeIn() reflect.Type {
	return reflect.TypeFor[A]()
}

func (this *IOMaybeFail[A]) TypeOut() reflect.Type {
	return reflect.TypeFor[A]()
}

func (this *IOMaybeFail[A]) SetDebug(b bool) {
	this.debug = b
}

func (this *IOMaybeFail[A]) SetDebugInfo(info *types.IODebugInfo) {
	this.debugInfo = info
}

func (this *IOMaybeFail[A]) GetDebugInfo() *types.IODebugInfo {
	return this.debugInfo
}
func (this *IOMaybeFail[A]) String() string {
	return fmt.Sprintf("MaybeFail(%v)", this.value.String())
}

func (this *IOMaybeFail[A]) SetPrevEffect(prev types.IOEffect) {
	this.prevEffect = prev
}

func (this *IOMaybeFail[A]) GetPrevEffect() *option.Option[types.IOEffect] {
	return option.Of(this.prevEffect)
}

func (this *IOMaybeFail[A]) GetResult() types.ResultOptionAny {
	return this.value.ToResultOfOption()
}

func (this *IOMaybeFail[A]) UnsafeRun() types.IOEffect {
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

				var res *result.Result[A]

				if this.fr != nil {
					res = this.fr(effValue)
				} else {
					res = result.Make(effValue, this.fe(effValue))
				}

				if res.HasError() {
					this.value = result.OfError[*option.Option[A]](res.Failure())
				} else {
					this.value = result.OfValue(option.Some(effValue))
				}

			} else {
				util.PanicCastType("IOMaybeFail",
					reflect.TypeOf(val), reflect.TypeFor[A]())

			}
		}
	}

	if this.debug {
		log.Printf("%v\n", this.String())
	}

	return currEff.(types.IOEffect)
}
