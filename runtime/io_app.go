package runtime

import (
	"fmt"
	"github.com/mobilemindtec/go-io/option"
	"github.com/mobilemindtec/go-io/result"
	"github.com/mobilemindtec/go-io/state"
	"github.com/mobilemindtec/go-io/types"
	"github.com/mobilemindtec/go-io/util"
	"log"
	"reflect"
	"runtime/debug"
)

type IOApp[T any] struct {
	stack          []types.IORunnable
	state          *state.State
	value          *result.Result[*option.Option[T]]
	resources      []types.IResourceIO
	_debug         bool
	showStackTrace bool
	fnCatch        func(error) *result.Result[*option.Option[T]]
}

func New[T any](effects ...types.IORunnable) *IOApp[T] {
	app := &IOApp[T]{
		stack: []types.IORunnable{},
		state: state.NewState(),
	}
	return app.Effects(effects...)
}

func (this *IOApp[T]) Debug() *IOApp[T] {
	this._debug = true
	return this
}

func (this *IOApp[T]) ShowStackTrace() *IOApp[T] {
	this.showStackTrace = true
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
	if suspended, ok := effect.(types.IIOSuspended); ok {
		this.Suspended(suspended)
	} else {
		effect.CheckTypesFlow()
		this.stack = append(this.stack, effect)
	}
	return this
}

func (this *IOApp[T]) Catch(f func(error) *result.Result[*option.Option[T]]) *IOApp[T] {
	this.fnCatch = f
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

func (this *IOApp[T]) Suspended(suspended types.IIOSuspended) *IOApp[T] {
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

	var lastEffect types.IOEffect

	for _, io := range this.stack {

		io.SetState(this.state)

		if this._debug {
			io.SetDebug(this._debug)
		}

		io.SetPrevEffect(lastEffect)
		varName := io.GetVarName()
		resultIO = io.UnsafeRunIO()
		lastEffect = io.GetLastEffect()

		if len(varName) == 0 {
			varName = fmt.Sprintf("__var__%v", this.state.Count())
		}

		if this._debug {
			log.Printf("var = %v, IO result = %v", varName, resultIO.String())
		}

		if resultIO.IsOk() && resultIO.Get().NonEmpty() {
			this.state.SetVar(varName, resultIO.Get().Get())
		} else {
			//break
		}
	}

	if resultIO.IsError() {

		if this.showStackTrace {
			debug.PrintStack()
			this.showStackTrace = false // only first error
		}

		if this.fnCatch != nil {
			this.value = this.fnCatch(resultIO.Failure())
		} else {
			this.value = result.OfError[*option.Option[T]](resultIO.Failure())
		}

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
