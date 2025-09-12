package pipeline

import (
	"errors"
	"fmt"
	"github.com/mobilemindtech/go-io/option"
	"github.com/mobilemindtech/go-io/result"
	"github.com/mobilemindtech/go-io/state"
	"github.com/mobilemindtech/go-io/types"
	"github.com/mobilemindtech/go-io/types/unit"
	"github.com/mobilemindtech/go-io/util"
	"log"
	"reflect"
	"runtime/debug"
)

type IPipeline interface {
	GetComputations() []*Computation
	UnsafeRunPipeline() types.ResultOptionAny
}

type Computation struct {
	varName  string
	action   interface{}
	funcInfo *util.FuncInfo
}

// Pipeline should return:
// - any
// - (any, error)
// - (string, any) -> var name, type
// - *result.Result[any]
// - *option.Option[any]
// - *result.Result[*option.Option[any]]
// Any return of error or None stop pipeline.
// The last computation should be return a same type of Pipeline[T] generic type.
type Pipeline[T any] struct {
	state             *state.State
	computations      []*Computation
	computationResult *result.Result[*option.Option[T]]
	debug             bool
}

// New create new Pipeline
func New[T any]() *Pipeline[T] {
	return &Pipeline[T]{
		state:             state.NewState(),
		computations:      []*Computation{},
		computationResult: result.OfValue(option.None[T]()),
	}
}

func (this *Pipeline[T]) GetComputations() []*Computation {
	return this.computations
}

func (this *Pipeline[T]) UnsafeRunPipeline() types.ResultOptionAny {
	return this.UnsafeRun().ToResultOfOption()
}

func (this *Pipeline[T]) getVarName() string {
	return fmt.Sprintf("__var__%v", this.state.Count())
}

func (this *Pipeline[T]) addComputation(f interface{}) *Pipeline[T] {
	varName := this.getVarName()
	funcInfo := util.NewFuncInfo(f)
	this.computations = append(this.computations, &Computation{action: f, funcInfo: funcInfo, varName: varName})
	return this
}

// Suspension combine suspended Pipeline with current pipeline, add on end computations
func (this *Pipeline[T]) Suspension(pipe IPipeline) *Pipeline[T] {
	for _, cpu := range pipe.GetComputations() {
		this.computations = append(this.computations, cpu)
	}
	return this
}

// UnsafeYield Pipeline unsafe result
func (this *Pipeline[T]) UnsafeYield() T {
	return this.computationResult.Get().Get()
}

// Yield Pipeline result
func (this *Pipeline[T]) Yield() *option.Option[T] {
	return this.computationResult.Get()
}

// Next Pipeline computation
func (this *Pipeline[T]) Next(f interface{}) *Pipeline[T] {
	return this.addComputation(f)
}

// UnsafeRun Run Pipeline
func (this *Pipeline[T]) UnsafeRun() (value *result.Result[*option.Option[T]]) {

	currStackPointer := -1

	defer func() {
		if r := recover(); r != nil {
			msg := fmt.Sprintf(
				"Pipeline error on StackPointer %v. Message %v. StackTrace: %v",
				currStackPointer, r, string(debug.Stack()))

			fmt.Println(fmt.Sprintf("<Pipeline recover> %v", msg))
			value = result.OfError[*option.Option[T]](errors.New(msg))
			this.computationResult = value
		}
	}()

	var lastResult interface{}

	for i, step := range this.computations {

		stateCopy := this.state.Copy()
		currStackPointer = i
		nextFnInfo := step.funcInfo
		currStateSize := this.state.Count()
		var fnParams []reflect.Value
		var fnResults []interface{}
		var fnResultTypes []reflect.Type
		var isErrorFunc bool

		handleResult := func(res []reflect.Value) {
			if len(res) > 2 {
				panic(fmt.Sprintf("return type count should be < 3, but is %v", len(res)))
			}
			for i := 0; i < len(res); i++ {
				fnResults = append(fnResults, res[i].Interface())
				fnResultTypes = append(fnResultTypes, res[i].Type())
			}

			if len(fnResultTypes) == 2 {
				firstType := fnResultTypes[0]
				secondType := fnResultTypes[1]

				isErrorFunc = secondType.Implements(reflect.TypeFor[error]())

				if !isErrorFunc {
					panic(fmt.Sprintf("func should be return (any, error), but return (%v, %v)",
						firstType.String(), secondType.String()))
				}
			} else if len(fnResultTypes) > 2 {
				panic(fmt.Sprintf("return type count should be < 3, but is %v",
					len(fnResultTypes)))
			}
		}

		if currStateSize < nextFnInfo.ArgsCount {
			panic(fmt.Sprintf(
				"expected %v args, but have %v state results", nextFnInfo.ArgsCount, currStateSize))
		}

		if this.debug {
			log.Print("step %v, args %v", i, nextFnInfo.ArgsCount)
		}

		for j := 0; j < nextFnInfo.ArgsCount; j++ {
			_, val := state.LookupVar(stateCopy, nextFnInfo.ArgType(j), true)
			fnParams = append(fnParams, val)
		}

		handleResult(nextFnInfo.Call(fnParams))

		if len(fnResults) == 0 {
			// ignore result
			continue
		}

		var fnResult interface{}
		varName := this.getVarName()
		var errrorResult error

		switch len(fnResults) {
		case 1:
			fnResult = fnResults[0]
			break
		default:

			if util.IsNotNil(fnResults[1]) {
				errrorResult = fnResults[1].(error)
			}
			fnResult = fnResults[0]

		}

		if isErrorFunc && util.IsNotNil(errrorResult) {
			value = result.OfError[*option.Option[T]](errrorResult)
			this.computationResult = value
			return
		}

		if rs, ok := fnResult.(result.IResult); ok {
			if rs.HasError() {
				value = result.OfError[*option.Option[T]](rs.GetError())
				this.computationResult = value
				return
			} else {
				lastResult = rs.GetValue()
				if util.IsNotNil(lastResult) {

					if opt, ok := lastResult.(option.IOption); ok {
						if opt.IsEmpty() {
							value = result.OfValue[*option.Option[T]](option.None[T]())
							this.computationResult = value
							return
						} else {
							this.state.SetVar(varName, opt.GetValue())
						}
					} else {
						this.state.SetVar(varName, lastResult)
					}
				} else {
					value = result.OfValue[*option.Option[T]](option.None[T]())
					this.computationResult = value
					return
				}
			}
		} else if opt, ok := fnResult.(option.IOption); ok {
			if opt.IsEmpty() {
				value = result.OfValue[*option.Option[T]](option.None[T]())
				this.computationResult = value
				return
			} else {
				lastResult = opt.GetValue()
				this.state.SetVar(varName, lastResult)
			}
		} else {
			lastResult = fnResult
			this.state.SetVar(varName, fnResult)
		}
	}

	if reflect.TypeFor[T]() == reflect.TypeFor[*unit.Unit]() {
		var unit interface{} = unit.OfUnit()
		value = result.OfValue(option.Some(unit.(T)))
		this.computationResult = value

	} else {
		r := result.Cast[T](lastResult)
		value = result.OfValue(option.Some(r.Get()))
		this.computationResult = value
	}
	return
}
