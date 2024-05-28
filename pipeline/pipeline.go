package pipeline

import (
	"errors"
	"fmt"
	"github.com/mobilemindtec/go-io/option"
	"github.com/mobilemindtec/go-io/result"
	"github.com/mobilemindtec/go-io/state"
	"github.com/mobilemindtec/go-io/types"
	"github.com/mobilemindtec/go-io/util"
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

type Pipeline[T any] struct {
	state             *state.State
	computations      []*Computation
	computationResult *result.Result[*option.Option[T]]
	debug             bool
}

func New[T any]() *Pipeline[T] {
	return &Pipeline[T]{
		state:             state.NewState(),
		computations:      []*Computation{},
		computationResult: result.OfValue(option.None[T]()),
	}
}

func (this Pipeline[T]) GetComputations() []*Computation {
	return this.computations
}

func (this Pipeline[T]) UnsafeRunPipeline() types.ResultOptionAny {
	return this.UnsafeRun().ToResultOfOption()
}

func (this Pipeline[T]) getVarName() string {
	return fmt.Sprintf("__var__%v", this.state.Count())
}

func (this Pipeline[T]) addComputation(f interface{}) Pipeline[T] {
	varName := this.getVarName()
	funcInfo := util.NewFuncInfo(f)
	this.computations = append(this.computations, &Computation{action: f, funcInfo: funcInfo, varName: varName})
	return this
}

func (this Pipeline[T]) Suspension(pipe IPipeline) Pipeline[T] {
	for _, cpu := range pipe.GetComputations() {
		this.computations = append(this.computations, cpu)
	}
	return this
}

func (this *Pipeline[T]) UnsafeYield() T {
	return this.computationResult.Get().Get()
}

func (this *Pipeline[T]) Yield() *option.Option[T] {
	return this.computationResult.Get()
}

func (this Pipeline[T]) Next(f interface{}) Pipeline[T] {
	return this.addComputation(f)
}

func (this Pipeline[T]) UnsafeRun() (value *result.Result[*option.Option[T]]) {

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

		handleResult := func(res []reflect.Value) {
			if len(res) > 2 {
				panic(fmt.Sprintf("return type count should be < 3, but is %v", len(res)))
			}
			for i := 0; i < len(res); i++ {
				fnResults = append(fnResults, res[i].Interface())
				fnResultTypes = append(fnResultTypes, res[i].Type())
			}

			if len(fnResultTypes) == 2 {
				firstType := fnResultTypes[1]
				if firstType.Kind() != reflect.String {
					panic(fmt.Sprintf("the first return type should be a string (var namr), but has %v",
						firstType.String()))
				}
			}
		}

		if currStateSize < nextFnInfo.ArgsCount {
			panic(fmt.Sprintf(
				"expected %v args, but have %v state results", nextFnInfo.ArgsCount, currStateSize))
		}

		for j := 0; j < nextFnInfo.ArgsCount; j++ {
			_, val := state.LookupVar(stateCopy, nextFnInfo.ArgType(j))
			fnParams = append(fnParams, val)
		}

		handleResult(nextFnInfo.Call(fnParams))

		if len(fnResults) == 0 {
			// ignore result
			continue
		}

		var fnResult interface{}
		varName := this.getVarName()

		switch len(fnResults) {
		case 1:
			fnResult = fnResults[0]
			break
		default:
			varName = fnResults[0].(string)
			fnResult = fnResults[1]
		}

		if rs, ok := fnResult.(result.IResult); ok {
			if rs.HasError() {
				value = result.OfError[*option.Option[T]](rs.GetError())
				this.computationResult = value
				return
			} else {
				lastResult = rs.GetValue()
				if util.IsNotNil(lastResult) {
					this.state.SetVar(varName, lastResult)
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

	r := result.Cast[T](lastResult)
	value = result.OfValue(option.Some(r.Get()))
	this.computationResult = value
	return
}
