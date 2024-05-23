package pipeline

import (
	"errors"
	"fmt"
	"github.com/mobilemindtec/go-io/option"
	"github.com/mobilemindtec/go-io/result"
	"github.com/mobilemindtec/go-io/util"
	"reflect"

	"runtime/debug"
)

type IPipeline interface {
	GetComputations() []*Computation
	UnsafeRunPipeline() *result.Result[any]
}

type Computation struct {
	varName  string
	action   interface{}
	funcInfo *util.FuncInfo
}

func NewStateItem(name string, value interface{}) *StateItem {
	return &StateItem{name: name, value: value, typ: reflect.TypeOf(value)}
}

type Pipeline[T any] struct {
	state        *State
	computations []*Computation
	//value        *result.Result[T]
	debug bool
}

func New[T any]() *Pipeline[T] {
	return &Pipeline[T]{
		state:        NewState(),
		computations: []*Computation{},
	}
}

func (this Pipeline[T]) GetComputations() []*Computation {
	return this.computations
}

func (this Pipeline[T]) UnsafeRunPipeline() *result.Result[any] {
	return this.UnsafeRun().ToResultOfAny()
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

func (this Pipeline[T]) Next(f interface{}) Pipeline[T] {
	this.computations = append(this.computations, &Computation{action: f})
	return this
}

func (this Pipeline[T]) findInCtxbyType(state *State, argType reflect.Type, stackPointer string) reflect.Value {

	var item reflect.Value
	var key string
	var found bool
	tuples := state.ToTuples()

	for i := len(tuples); i >= 0; i-- {
		tp := tuples[i]
		rtype := reflect.TypeOf(tp.val)
		if rtype == argType {
			item = reflect.ValueOf(tp.val)
			key = tp.key
			found = true
			break
		}
	}

	if found {
		// consume value and delete from key
		state.Delete(key)
		return item
	}

	if argType == reflect.TypeOf(state) {
		return reflect.ValueOf(state)
	}

	panic(fmt.Sprintf("StackPointer %v: arg type %v not found in state results", stackPointer, argType))
}

func (this Pipeline[T]) UnsafeRun() (value *result.Result[T]) {

	currStackPointer := -1

	defer func() {
		if r := recover(); r != nil {
			msg := fmt.Sprintf(
				"Pipeline error on StackPointer %v. Message %v. StackTrace: %v",
				currStackPointer, r, string(debug.Stack()))

			fmt.Println(fmt.Sprintf("<Pipeline recover> %v", msg))
			value = result.OfError[T](errors.New(msg))
		}
	}()

	var lastResult interface{}

	for i, step := range this.computations {

		stateCopy := this.state.ToCopy()
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

		//r = step.action.(ActionR2)(this.results[0], this.results[1])
		for i = 0; i < nextFnInfo.ArgsCount; i++ {
			fnParams = append(fnParams,
				this.findInCtxbyType(stateCopy, nextFnInfo.ArgType(i), fmt.Sprintf("%v", i)))
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
				value = result.OfError[T](rs.GetError())
				return
			} else {
				lastResult = rs.GetValue()
				if util.IsNotNil(lastResult) {
					this.state.SetVar(varName, lastResult)
				} else {
					value = result.OfNil[T]()
					return
				}
			}
		} else if opt, ok := fnResult.(option.IOption); ok {
			if opt.IsEmpty() {
				value = result.OfNil[T]()
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

	value = result.Cast[T](lastResult)
	return
}
