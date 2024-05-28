package runtime

import (
	"fmt"
	"github.com/mobilemindtec/go-io/option"
	"github.com/mobilemindtec/go-io/result"
	"github.com/mobilemindtec/go-io/state"
	"github.com/mobilemindtec/go-io/types"
	"github.com/mobilemindtec/go-io/util"
	"reflect"
)

type Runtime[T any] struct {
	stack     []types.RuntimeIO
	state     *state.State
	value     *result.Result[*option.Option[T]]
	resources []types.IResourceIO
	debug     bool
}

func New[T any]() *Runtime[T] {
	return &Runtime[T]{
		stack: []types.RuntimeIO{},
		state: state.NewState(),
	}
}

func (this *Runtime[T]) Debug() *Runtime[T] {
	this.debug = true
	return this
}

func (this *Runtime[T]) DebugOn() {
	this.Debug()
}

func (this *Runtime[T]) ConsumeVar(name string) interface{} {
	return this.state.Consume(name)
}

func (this *Runtime[T]) UnsafeRunRuntime() types.ResultOptionAny {
	return this.UnsafeRun().ToResultOfOption()
}

func (this *Runtime[T]) Var(name string) interface{} {
	return this.state.Var(name)
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

func (this *Runtime[T]) UnsafeYield() T {
	return this.value.Get().Get()
}

func (this *Runtime[T]) Yield() *option.Option[T] {
	return this.value.Get()
}

func (this *Runtime[T]) Effects(effects ...types.RuntimeIO) *Runtime[T] {
	for _, eff := range effects {
		this.Effect(eff)
	}
	return this
}

func (this *Runtime[T]) IO(ios ...types.IOEffect) *Runtime[T] {
	this.Effects(types.NewIO[T]().Effects(ios...))
	return this
}

func (this *Runtime[T]) UnsafeRun() *result.Result[*option.Option[T]] {

	var resultIO types.ResultOptionAny
	this.value = result.OfValue(option.None[T]())

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

		this.state.SetVar(varName, res.Get())
	}

	for _, io := range this.stack {

		io.SetState(this.state)
		io.SetDebug(this.debug)

		varName := io.GetVarName()
		resultIO = io.UnsafeRunIO()

		if len(varName) == 0 {
			varName = fmt.Sprintf("__var__%v", this.state.Count())
		}

		if this.debug {
			fmt.Println("var = ", varName, ", IO result = ", resultIO.String())
		}

		if resultIO.IsOk() && resultIO.Get().NonEmpty() {
			this.state.SetVar(varName, resultIO.Get().Get())
		} else {
			break
		}
	}

	if resultIO.IsError() {
		this.value = result.OfError[*option.Option[T]](resultIO.Failure())
	} else {

		r := resultIO.Get()

		if r.NonEmpty() {
			if effValue, ok := r.GetValue().(T); ok {
				this.value = result.OfValue(option.Some(effValue))
			} else {
				util.PanicCastType("Runtime",
					reflect.TypeOf(r.GetValue()), reflect.TypeFor[T]())

			}
		}
	}

	for _, r := range this.resources {
		r.Close().
			IfError(func(err error) {
				panic(fmt.Sprintf(
					"fail on close resource %v: %v", r.GetVarName(), err))
			})
	}

	return this.value
}
