package types

import (
	"github.com/mobilemindtec/go-io/option"
	"github.com/mobilemindtec/go-io/result"
)

type IOEffect interface {
	GetPrevEffect() *option.Option[IOEffect]
	SetPrevEffect(IOEffect)
	GetResult() *result.Result[any]
	UnsafeRun() IOEffect
}

type RuntimeIO interface {
	UnsafeRunIO() *result.Result[any]
	GetVarName() string
}
