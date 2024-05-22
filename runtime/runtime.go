package runtime

import (
	"errors"
	"fmt"
	"github.com/mobilemindtec/go-io/option"
	"github.com/mobilemindtec/go-io/result"
	"github.com/mobilemindtec/go-io/types"
	"github.com/mobilemindtec/go-io/util"
)

type IRuntime interface {
	ConsumeVar(name string) interface{}
	Var(name string) interface{}
}

type Runtime[T any] struct {
	stack       []types.RuntimeIO
	allocations map[string]interface{}
	value       *result.Result[T]
	resources   []types.IResourceIO
}

func NewRuntime[T any]() *Runtime[T] {
	return &Runtime[T]{
		stack:       []types.RuntimeIO{},
		allocations: map[string]interface{}{},
	}
}

func (this *Runtime[T]) ToInterface() IRuntime {
	var rt interface{} = this
	return rt.(IRuntime)
}

func (this *Runtime[T]) ConsumeVar(name string) interface{} {
	val, ok := this.allocations[name]
	if ok {
		delete(this.allocations, name)
		return val
	}
	return nil
}

func (this *Runtime[T]) Var(name string) interface{} {
	val, ok := this.allocations[name]
	if ok {
		return val
	}
	return nil
}

func (this *Runtime[T]) Resource(res types.IResourceIO) *Runtime[T] {
	this.resources = append(this.resources, res)
	return this
}

func (this *Runtime[T]) Resources(res ...types.IResourceIO) *Runtime[T] {
	for _, r := range res {
		this.Resource(r)
	}
	return this
}

func (this *Runtime[T]) Effect(effect types.RuntimeIO) *Runtime[T] {
	this.stack = append(this.stack, effect)
	return this
}

func (this *Runtime[T]) Yield() T {
	return this.value.Get()
}

func (this *Runtime[T]) Effects(effects ...types.RuntimeIO) *Runtime[T] {
	for _, eff := range effects {
		this.Effect(eff)
	}
	return this
}

func (this *Runtime[T]) UnsafeRun() *result.Result[T] {

	var resultIO *result.Result[any]

	for _, r := range this.resources {
		res := r.Open()
		varName := r.GetVarName()

		res.IfError(func(err error) {
			panic(fmt.Sprintf(
				"fail on open resource %v: %v", r.GetVarName(), err))
		})

		res.ToOption().
			IfEmpty(func() {
				panic(fmt.Sprintf(
					"fail on open resource %v: resource not found", r.GetVarName()))
			})

		this.allocations[varName] = res.Get()
	}

	for _, io := range this.stack {
		varName := io.GetVarName()
		resultIO = io.UnsafeRunIO()

		if len(varName) == 0 {
			varName = fmt.Sprintf("__var__%v", len(this.allocations))
		}

		if resultIO.OptionNonEmpty() {
			this.allocations[varName] = resultIO.Get()
		} else {
			break
		}
	}

	if resultIO.ToOption().NonEmpty() {
		this.value = result.OfValue(resultIO.Get().(T))
	} else {
		this.value = result.OfError[T](resultIO.Error())
	}

	for _, r := range this.resources {
		r.Close().IfError(func(err error) {
			panic(fmt.Sprintf(
				"fail on close resource %v: %v", r.GetVarName(), err))
		})
	}

	return this.value
}

func ConsumeVar[T any](rt IRuntime, name string) (T, error) {
	opt := SafeConsumeVar[T](rt, name)
	if opt.NonEmpty() {
		return opt.Get(), nil
	}
	var t T
	return t, errors.New(fmt.Sprintf("variable %v not found on runtime", name))
}

func SafeConsumeVar[T any](rt IRuntime, name string) *option.Option[T] {
	value := rt.ConsumeVar(name)
	if !util.IsNil(value) {
		if val, ok := value.(T); ok {
			return option.Of(val)
		}
	}
	return option.None[T]()
}

func Var[T any](rt IRuntime, name string) (T, error) {
	opt := SafeVar[T](rt, name)
	if opt.NonEmpty() {
		return opt.Get(), nil
	}
	var t T
	return t, errors.New(fmt.Sprintf("variable %v not found on runtime", name))
}

func SafeVar[T any](rt IRuntime, name string) *option.Option[T] {
	value := rt.Var(name)
	if !util.IsNil(value) {
		if val, ok := value.(T); ok {
			return option.Of(val)
		}
	}
	return option.None[T]()
}
