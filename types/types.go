package types

import (
	"github.com/mobilemindtec/go-io/option"
	"github.com/mobilemindtec/go-io/result"
	"github.com/mobilemindtec/go-io/state"
	"reflect"
)

type ResultOptionAny = *result.Result[*option.Option[any]]

type IODebugInfo struct {
	Line     int
	Filename string
}

type IOEffect interface {
	GetPrevEffect() *option.Option[IOEffect]
	SetPrevEffect(IOEffect)
	GetResult() ResultOptionAny
	UnsafeRun() IOEffect
	SetDebug(bool)
	String() string
	TypeIn() reflect.Type
	TypeOut() reflect.Type
	SetDebugInfo(*IODebugInfo)
	GetDebugInfo() *IODebugInfo
}

type IOStateful interface {
	SetState(*state.State)
}

type IOLift[T any] interface {
	Lift() *IO[T]
}

type IORunnable interface {
	UnsafeRunIO() ResultOptionAny
	GetVarName() string
	SetDebug(bool)
	SetState(*state.State)
	CheckTypesFlow()
	IOType() reflect.Type
	GetLastEffect() IOEffect
	SetPrevEffect(IOEffect)
}

type IOApp interface {
	ConsumeVar(name string) interface{}
	Var(name string) interface{}
	UnsafeRunApp() ResultOptionAny
	DebugOn()
}

type IOError struct {
	Message    string
	StackTrace string
}

func NewIOError(message string, stacktrace []byte) *IOError {
	return &IOError{Message: message, StackTrace: string(stacktrace)}
}

func (this IOError) Error() string {
	return this.Message
}
