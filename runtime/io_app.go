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

type IOApp[T any] struct {
	stack     []types.IORunnable
	state     *state.State
	value     *result.Result[*option.Option[T]]
	resources []types.IResourceIO
	debug     bool
}

func New[T any](effects ...types.IORunnable) *IOApp[T] {
	app := &IOApp[T]{
		stack: []types.IORunnable{},
		state: state.NewState(),
	}
	return app.Effects(effects...)
}

func (this *IOApp[T]) Debug() *IOApp[T] {
	this.debug = true
	return this
}

func (this *IOApp[T]) DebugOn() {
	this.Debug()
}

func (this *IOApp[T]) ConsumeVar(name string) interface{} {
	return this.state.Consume(name)
}

func (this *IOApp[T]) UnsafeRunApp() types.ResultOptionAny {
	return this.UnsafeRun().ToResultOfOption()
}

func (this *IOApp[T]) Var(name string) interface{} {
	return this.state.Var(name)
}

func (this *IOApp[T]) Resource(res types.IResourceIO) *IOApp[T] {
	this.resources = append(this.resources, res)
	return this
}

func (this *IOApp[T]) Resources(res ...types.IResourceIO) *IOApp[T] {
	for _, r := range res {
		this.Resource(r)
	}
	return this
}

func (this *IOApp[T]) Effect(effect types.IORunnable) *IOApp[T] {
	if suspended, ok := effect.(*types.IOSuspended); ok {
		this.Suspended(suspended)
	} else {
		this.stack = append(this.stack, effect)
	}
	return this
}

func (this *IOApp[T]) UnsafeYield() T {
	return this.value.Get().Get()
}

func (this *IOApp[T]) Yield() *option.Option[T] {
	return this.value.Get()
}

func (this *IOApp[T]) Continue(effects ...types.IORunnable) *IOApp[T] {
	return this.Effects(effects...)
}

func (this *IOApp[T]) Effects(effects ...types.IORunnable) *IOApp[T] {
	for _, eff := range effects {
		this.Effect(eff)
	}
	return this
}

func (this *IOApp[T]) IO(ios ...types.IOEffect) *IOApp[T] {
	return this.Effects(types.NewIO[T]().Effects(ios...))
}

func (this *IOApp[T]) Suspended(suspended *types.IOSuspended) *IOApp[T] {
	this.Effects(suspended.IOs()...)
	return this
}

func (this *IOApp[T]) UnsafeRun() *result.Result[*option.Option[T]] {

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
				util.PanicCastType("IOApp",
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
