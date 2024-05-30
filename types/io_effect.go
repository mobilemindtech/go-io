package types

import (
	"github.com/mobilemindtec/go-io/option"
	"github.com/mobilemindtec/go-io/result"
	"github.com/mobilemindtec/go-io/state"
)

type ResultOptionAny = *result.Result[*option.Option[any]]

type IOEffect interface {
	GetPrevEffect() *option.Option[IOEffect]
	SetPrevEffect(IOEffect)
	GetResult() ResultOptionAny
	UnsafeRun() IOEffect
	SetDebug(bool)
	SetState(*state.State)
	String() string
}

type IOLift[T any] interface {
	Lift() *IO[T]
}

type IORunnable interface {
	UnsafeRunIO() ResultOptionAny
	GetVarName() string
	SetDebug(bool)
	SetState(*state.State)
}

type IOApp interface {
	ConsumeVar(name string) interface{}
	Var(name string) interface{}
	UnsafeRunApp() ResultOptionAny
	DebugOn()
}
