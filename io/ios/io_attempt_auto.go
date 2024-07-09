package ios

import (
	"fmt"
	"github.com/mobilemindtec/go-io/either"
	"github.com/mobilemindtec/go-io/option"
	"github.com/mobilemindtec/go-io/result"
	"github.com/mobilemindtec/go-io/state"
	"github.com/mobilemindtec/go-io/types"
	"github.com/mobilemindtec/go-io/util"
	"log"
	"reflect"
)

type IOAttemptAuto[A any] struct {
	value      *result.Result[*option.Option[A]]
	prevEffect types.IOEffect

	fnAuto interface{}

	state     *state.State
	debug     bool
	debugInfo *types.IODebugInfo
}

// NewAttemptAuto f should be a func that return:
// - type A
// - type (A, error)
// - type *option.Option[A],
// - type *result.Result[A]
// - type *result.Result[*option.Option[A]
// - type *either.Either[?, ?]
// the func can br 0 or N args that will be sought in the state.State
func NewAttemptAuto[A any](f interface{}) *IOAttemptAuto[A] {
	return &IOAttemptAuto[A]{fnAuto: f}
}

func (this *IOAttemptAuto[A]) Lift() *types.IO[A] {
	return types.NewIO[A]().Effects(this)
}

func (this *IOAttemptAuto[A]) SetDebug(b bool) {
	this.debug = b
}

func (this *IOAttemptAuto[A]) SetDebugInfo(info *types.IODebugInfo) {
	this.debugInfo = info
}

func (this *IOAttemptAuto[A]) GetDebugInfo() *types.IODebugInfo {
	return this.debugInfo
}

func (this *IOAttemptAuto[A]) TypeIn() reflect.Type {
	return reflect.TypeFor[*types.Unit]()
}

func (this *IOAttemptAuto[A]) TypeOut() reflect.Type {
	return reflect.TypeFor[A]()
}

func (this *IOAttemptAuto[A]) String() string {
	return fmt.Sprintf("AttemptAuto(fn=%v, value=%v)", this.getFuncName(), this.value.String())
}

func (this *IOAttemptAuto[A]) SetState(st *state.State) {
	this.state = st
}

func (this *IOAttemptAuto[A]) SetPrevEffect(prev types.IOEffect) {
	this.prevEffect = prev
}

func (this *IOAttemptAuto[A]) GetPrevEffect() *option.Option[types.IOEffect] {
	return option.Of(this.prevEffect)
}

func (this *IOAttemptAuto[A]) GetResult() types.ResultOptionAny {
	return this.value.ToResultOfOption()
}

func (this *IOAttemptAuto[A]) UnsafeRun() types.IOEffect {

	var currEff interface{} = this
	prevEff := this.GetPrevEffect()
	this.value = result.OfValue(option.None[A]())
	execute := true
	hasPrev := prevEff.NonEmpty()

	if hasPrev {
		prev := prevEff.Get()
		if prev.GetResult().IsError() {
			this.value = result.OfError[*option.Option[A]](
				prevEff.Get().GetResult().Failure())
			execute = false
		} else {
			execute = prev.GetResult().Get().IsSome()
		}
	}

	defer func() {
		if r := recover(); r != nil {
			this.value = RecoverIO[A](this, this.debug, this.debugInfo, r)
		}
	}()

	if execute { // not error

		info := util.NewFuncInfo(this.fnAuto)
		var fnParams []reflect.Value
		stateCopy := this.state.Copy()
		for i := 0; i < info.ArgsCount; i++ {
			_, val := state.LookupVar(stateCopy, info.ArgType(i), true)
			fnParams = append(fnParams, val)
		}

		fnResults := info.Call(fnParams)

		switch len(fnResults) {
		case 0:
			var unit interface{} = types.OfUnit()
			this.value = result.OfValue(option.Some(unit.(A)))
		case 1: // should be a Result[A]

			rVal := fnResults[0]

			if util.CanNil(rVal.Kind()) && rVal.IsNil() {
				panic("func result can't be null")
			}

			retValue := rVal.Interface()

			if r, ok := retValue.(result.IResult); ok {

				if r.HasError() {
					this.value = result.OfError[*option.Option[A]](r.GetError())
				} else {
					if o, ok := r.GetValue().(option.IOption); ok {
						if o.IsEmpty() {
							this.value = result.OfValue(option.None[A]())
						} else {

							if val, ok := o.GetValue().(A); ok {
								this.value = result.OfValue(option.Some(val))
							} else {
								util.PanicCastType("IOAttempt",
									reflect.TypeOf(o.GetValue()), reflect.TypeFor[A]())
							}
						}
					} else if val, ok := r.GetValue().(A); ok {
						this.value = result.OfValue(option.Some(val))
					} else {
						panic("wrong result type")
					}
				}

			} else if opt, ok := retValue.(option.IOption); ok {

				if opt.IsEmpty() {
					this.value = result.OfValue(option.None[A]())
				} else {
					if val, ok := opt.GetValue().(A); ok {
						this.value = result.OfValue(option.Some(val))
					} else {
						util.PanicCastType("IOAttempt",
							reflect.TypeOf(opt.GetValue()), reflect.TypeFor[A]())
					}
				}

			} else if _, ok := retValue.(either.IEither); ok {
				if eia, ok := retValue.(A); ok {
					this.value = result.OfValue[*option.Option[A]](option.Some(eia))
				} else {
					util.PanicCastType("IOAttempt",
						reflect.TypeOf(retValue), reflect.TypeFor[A]())
				}
			} else {
				panic("wrong result type")
			}

		case 2: // should be a (A, error)

			rA := fnResults[0]
			rE := fnResults[1]
			canNil := util.CanNil(rA.Kind())

			if (canNil && rA.IsNil()) && rE.IsNil() {
				panic("func result can't be null")
			}

			retTypeA := rA.Interface()
			retTypeE := rE.Interface()

			var valA A
			var valE error
			var ok bool

			if !canNil || !rA.IsNil() {
				valA, ok = rA.Interface().(A)
				if !ok {
					panic(fmt.Sprintf("cat't cast %v to %b", retTypeA, reflect.TypeOf(valA)))
				}
			}

			if !rE.IsNil() {
				valE, ok = rE.Interface().(error)
				if !ok {
					panic(fmt.Sprintf("cat't cast %v to %b", retTypeE, reflect.TypeOf(valE)))
				}
			}

			this.value = result.TryOption(func() (A, error) {
				return valA, valE
			})

		default:
			panic("func should be result.Result[*option.Option[A]] or (A, error)")
		}

	}

	if this.debug {
		log.Printf("%v\n", this.String())
	}

	return currEff.(types.IOEffect)
}

func (this *IOAttemptAuto[A]) getFuncName() string {
	if this.fnAuto != nil {
		return "fnAuto"
	}
	return "-"
}
