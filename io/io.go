package io

import (
	"github.com/mobilemindtec/go-io/either"
	"github.com/mobilemindtec/go-io/option"
	"github.com/mobilemindtec/go-io/pipeline"
	"github.com/mobilemindtec/go-io/result"
	"github.com/mobilemindtec/go-io/runtime"
	"github.com/mobilemindtec/go-io/state"
	"github.com/mobilemindtec/go-io/types"
)

func IOUnit(effs ...types.IOEffect) *types.IO[*types.Unit] {
	return IO[*types.Unit](effs...)
}
func IO[T any](effs ...types.IOEffect) *types.IO[T] {
	return types.NewIO[T]().Effects(effs...)
}

func Attempt[A any](f func() *result.Result[A]) *types.IOAttempt[A] {
	return types.NewAttempt(f)
}

func AttemptOfOption[A any](f func() *option.Option[A]) *types.IOAttempt[A] {
	return types.NewAttemptOfOption(f)
}

func AttemptOfResultOption[A any](f func() *result.Result[*option.Option[A]]) *types.IOAttempt[A] {
	return types.NewAttemptOfResultOption(f)
}

func AttemptState[A any](f func(*state.State) *result.Result[A]) *types.IOAttempt[A] {
	return types.NewAttemptState(f)
}

func AttemptStateOfOption[A any](f func(*state.State) *option.Option[A]) *types.IOAttempt[A] {
	return types.NewAttemptStateOfOption(f)
}

func AttemptStateOfResultOption[A any](f func(*state.State) *result.Result[*option.Option[A]]) *types.IOAttempt[A] {
	return types.NewAttemptStateOfResultOption(f)
}

func AttemptOfResultEither[E error, A any](f func() *result.Result[*either.Either[E, A]]) *types.IOAttempt[*either.Either[E, A]] {
	return types.NewAttemptOfResultEither(f)
}

func AttemptStateOfResultEither[E error, A any](f func(*state.State) *result.Result[*either.Either[E, A]]) *types.IOAttempt[*either.Either[E, A]] {
	return types.NewAttemptStateOfResultEither(f)
}

func AttemptOfEither[E error, A any](f func() *either.Either[E, A]) *types.IOAttempt[*either.Either[E, A]] {
	return types.NewAttemptOfEither(f)
}

func AttemptStateOfEither[E error, A any](f func(*state.State) *either.Either[E, A]) *types.IOAttempt[*either.Either[E, A]] {
	return types.NewAttemptStateOfEither(f)
}

func AttemptAuto[A any](f interface{}) *types.IOAttempt[A] {
	return types.NewAttemptAuto[A](f)
}

func AttemptOfUnit(f func()) *types.IOAttempt[*types.Unit] {
	return types.NewAttemptOfUnit[*types.Unit](f)
}

func AttemptStateOfUnit(f func(*state.State)) *types.IOAttempt[*types.Unit] {
	return types.NewAttemptStateOfUnit[*types.Unit](f)
}

func AttemptOfError[A any](f func() (A, error)) *types.IOAttempt[A] {
	return types.NewAttemptOfError(f)
}

func AttemptStateOfError[A any](f func(*state.State) (A, error)) *types.IOAttempt[A] {
	return types.NewAttemptStateOfError(f)
}

func AttemptExec[A any](f func(A)) *types.IOAttempt[A] {
	return types.NewAttemptExec(f)
}

func AttemptExecState[A any](f func(A, *state.State)) *types.IOAttempt[A] {
	return types.NewAttemptExecState(f)
}

func AttemptFlow[A any](f func(A) *result.Result[A]) *types.IOAttempt[A] {
	return types.NewAttemptFlowOfResult(f)
}

func AttemptFlowState[A any](f func(A, *state.State) *result.Result[A]) *types.IOAttempt[A] {
	return types.NewAttemptFlowStateOfResult(f)
}

func AttemptFlowOfResultOpiton[A any](f func(A) *result.Result[*option.Option[A]]) *types.IOAttempt[A] {
	return types.NewAttemptFlowOfResultOption(f)
}

func AttemptFlowStateOfResultOpiton[A any](f func(A, *state.State) *result.Result[*option.Option[A]]) *types.IOAttempt[A] {
	return types.NewAttemptFlowStateOfResultOption(f)
}

func Debug[A any](label string) *types.IODebug[A] {
	return types.NewDebug[A](label)
}

func MaybeFail[A any](f func(A) *result.Result[A]) *types.IOMaybeFail[A] {
	return types.NewMaybeFail[A](f)
}

func MaybeFailError[A any](f func(A) error) *types.IOMaybeFail[A] {
	return types.NewMaybeFailError[A](f)
}

func Filter[A any](f func(A) bool) *types.IOFilter[A] {
	return types.NewFilter[A](f)
}

func Exec[A any](f func(A)) *types.IOExec[A] {
	return types.NewExec[A](f)
}

func ExecState[A any](f func(A, *state.State)) *types.IOExec[A] {
	return types.NewExecState[A](f)
}

func FlatMap[A any, B any](f func(A) *types.IO[B]) *types.IOFlatMap[A, B] {
	return types.NewFlatMap[A, B](f)
}

func Map[A any, B any](f func(A) B) *types.IOMap[A, B] {
	return types.NewMap[A, B](f)
}

func PureVal[T any](value T) *types.IOPure[T] {
	return types.NewPureValue[T](value)
}

func Pure[T any](f func() T) *types.IOPure[T] {
	return types.NewPure[T](f)
}

func PureState[T any](f func(*state.State) T) *types.IOPure[T] {
	return types.NewPureState[T](f)
}

func Recover[A any](f func(error) A) *types.IORecover[A] {
	return types.NewRecover[A](f)
}

func SliceFilter[A any](f func(A) bool) *types.IOSliceFilter[A] {
	return types.NewSliceFilter[A](f)
}

func SliceFlatMap[A any, B any](f func(A) *types.IO[B]) *types.IOSliceFlatMap[A, B] {
	return types.NewSliceFlatMap[A, B](f)
}

func SliceForeach[A any](f func(A)) *types.IOSliceForeach[A] {
	return types.NewSliceForeach[A](f)
}

func SliceMap[A any, B any](f func(A) B) *types.IOSliceMap[A, B] {
	return types.NewSliceMap[A, B](f)
}

func Tap[A any](f func(A) bool) *types.IOTap[A] {
	return types.NewTap[A](f)
}

func Foreach[A any](f func(A)) *types.IOForeach[A] {
	return types.NewForeach[A](f)
}

func Or[A any](f func() A) *types.IOOr[A] {
	return types.NewOr[A](f)
}

func FailIfEmpty[A any](f func() error) *types.IOFailIfEmpty[A] {
	return types.NewFailIfEmpty[A](f)
}

func ExecIfEmpty[A any](f func()) *types.IOExecIfEmpty[A] {
	return types.NewExecIfEmpty[A](f)
}

func OrElse[A any](f func() *types.IO[A]) *types.IOOrElse[A] {
	return types.NewOrElse[A](f)
}

func CatchAll[A any](f func(error) *result.Result[*option.Option[A]]) *types.IOCatchAll[A] {
	return types.NewCatchAll[A](f)
}

func CatchAllOfResultOption[A any](f func(error) *result.Result[A]) *types.IOCatchAll[A] {
	return types.NewCatchAllOfResult[A](f)
}

func CatchAllOfOption[A any](f func(error) *option.Option[A]) *types.IOCatchAll[A] {
	return types.NewCatchAllOfOption[A](f)
}

func IOApp[T any](effects ...types.IORunnable) *runtime.IOApp[T] {
	return runtime.New[T](effects...)
}

func IOAppOfUnit(effects ...types.IORunnable) *runtime.IOApp[*types.Unit] {
	return runtime.New[*types.Unit](effects...)
}

func Suspend(vals ...types.IORunnable) *types.IOSuspended {
	return types.NewIOSuspended(vals...)
}

func Pipeline[T any]() *pipeline.Pipeline[T] {
	return pipeline.New[T]()
}
