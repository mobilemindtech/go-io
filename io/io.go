package io

import (
	"github.com/mobilemindtec/go-io/either"
	"github.com/mobilemindtec/go-io/io/ios"
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

func Attempt[A any](f func() *result.Result[A]) *ios.IOAttempt[A] {
	return ios.NewAttempt(f)
}

func AttemptOfOption[A any](f func() *option.Option[A]) *ios.IOAttempt[A] {
	return ios.NewAttemptOfOption(f)
}

func AttemptOfResultOption[A any](f func() *result.Result[*option.Option[A]]) *ios.IOAttempt[A] {
	return ios.NewAttemptOfResultOption(f)
}

func AttemptState[A any](f func(*state.State) *result.Result[A]) *ios.IOAttempt[A] {
	return ios.NewAttemptState(f)
}

func AttemptRunIOIfEmpty[A any](f func(*state.State) types.IORunnable) *ios.IOAttempt[A] {
	return ios.NewAttemptRunIOIfEmpty[A](f)
}

func AttemptRunIO[A any](f func(*state.State) types.IORunnable) *ios.IOAttempt[A] {
	return ios.NewAttemptRunIO[A](f)
}

func AttemptPureState[A any](f func(*state.State) A) *ios.IOAttempt[A] {
	return ios.NewAttemptPureState(f)
}

func AttemptStateOfOption[A any](f func(*state.State) *option.Option[A]) *ios.IOAttempt[A] {
	return ios.NewAttemptStateOfOption(f)
}

func AttemptStateOfResultOption[A any](f func(*state.State) *result.Result[*option.Option[A]]) *ios.IOAttempt[A] {
	return ios.NewAttemptStateOfResultOption(f)
}

func AttemptOfResultEither[E error, A any](f func() *result.Result[*either.Either[E, A]]) *ios.IOAttempt[*either.Either[E, A]] {
	return ios.NewAttemptOfResultEither(f)
}

func AttemptStateOfResultEither[E error, A any](f func(*state.State) *result.Result[*either.Either[E, A]]) *ios.IOAttempt[*either.Either[E, A]] {
	return ios.NewAttemptStateOfResultEither(f)
}

func AttemptOfEither[E error, A any](f func() *either.Either[E, A]) *ios.IOAttempt[*either.Either[E, A]] {
	return ios.NewAttemptOfEither(f)
}

func AttemptStateOfEither[E error, A any](f func(*state.State) *either.Either[E, A]) *ios.IOAttempt[*either.Either[E, A]] {
	return ios.NewAttemptStateOfEither(f)
}

func AttemptAuto[A any](f interface{}) *ios.IOAttempt[A] {
	return ios.NewAttemptAuto[A](f)
}

func AttemptOfUnit(f func()) *ios.IOAttempt[*types.Unit] {
	return ios.NewAttemptOfUnit[*types.Unit](f)
}

func AttemptStateOfUnit(f func(*state.State)) *ios.IOAttempt[*types.Unit] {
	return ios.NewAttemptStateOfUnit[*types.Unit](f)
}

func AttemptOfError[A any](f func() (A, error)) *ios.IOAttempt[A] {
	return ios.NewAttemptOfError(f)
}

func AttemptStateOfError[A any](f func(*state.State) (A, error)) *ios.IOAttempt[A] {
	return ios.NewAttemptStateOfError(f)
}

func AttemptExec[A any](f func(A)) *ios.IOAttempt[A] {
	return ios.NewAttemptExec(f)
}

func AttemptExecState[A any](f func(A, *state.State)) *ios.IOAttempt[A] {
	return ios.NewAttemptExecState(f)
}

func AttemptExecIfEmpty[A any](f func()) *ios.IOAttempt[A] {
	return ios.NewAttemptExecIfEmpty[A](f)
}

func AttemptExecIfEmptyState[A any](f func(*state.State)) *ios.IOAttempt[A] {
	return ios.NewAttemptExecIfEmptyState[A](f)
}

func AttemptFlow[A any](f func(A) *result.Result[A]) *ios.IOAttempt[A] {
	return ios.NewAttemptFlowOfResult(f)
}

func AttemptFlowState[A any](f func(A, *state.State) *result.Result[A]) *ios.IOAttempt[A] {
	return ios.NewAttemptFlowStateOfResult(f)
}

func AttemptFlowOfResultOpiton[A any](f func(A) *result.Result[*option.Option[A]]) *ios.IOAttempt[A] {
	return ios.NewAttemptFlowOfResultOption(f)
}

func AttemptFlowStateOfResultOpiton[A any](f func(A, *state.State) *result.Result[*option.Option[A]]) *ios.IOAttempt[A] {
	return ios.NewAttemptFlowStateOfResultOption(f)
}

func Debug[A any](label string) *ios.IODebug[A] {
	return ios.NewDebug[A](label)
}

func MaybeFail[A any](f func(A) *result.Result[A]) *ios.IOMaybeFail[A] {
	return ios.NewMaybeFail[A](f)
}

func MaybeFailError[A any](f func(A) error) *ios.IOMaybeFail[A] {
	return ios.NewMaybeFailError[A](f)
}

func Filter[A any](f func(A) bool) *ios.IOFilter[A] {
	return ios.NewFilter[A](f)
}

func FlatMap[A any, B any](f func(A) types.IORunnable) *ios.IOFlatMap[A, B] {
	return ios.NewFlatMap[A, B](f)
}

func Map[A any, B any](f func(A) B) *ios.IOMap[A, B] {
	return ios.NewMap[A, B](f)
}

func PureVal[T any](value T) *ios.IOPure[T] {
	return ios.NewPureValue[T](value)
}

func Pure[T any](f func() T) *ios.IOPure[T] {
	return ios.NewPure[T](f)
}

func Recover[A any](f func(error) A) *ios.IORecover[A] {
	return ios.NewRecover[A](f)
}

func SliceFilter[A any](f func(A) bool) *ios.IOSliceFilter[A] {
	return ios.NewSliceFilter[A](f)
}

func SliceFlatMap[A any, B any](f func(A) types.IORunnable) *ios.IOSliceFlatMap[A, B] {
	return ios.NewSliceFlatMap[A, B](f)
}

func SliceForeach[A any](f func(A)) *ios.IOSliceForeach[A] {
	return ios.NewSliceForeach[A](f)
}

func SliceAttemptIfEmpty[A any](f func() *result.Result[[]A]) *ios.IOSliceAttemptIfEmpty[A] {
	return ios.NewSliceAttemptIfEmpty[A](f)
}

func SliceMap[A any, B any](f func(A) B) *ios.IOSliceMap[A, B] {
	return ios.NewSliceMap[A, B](f)
}

func AsSliceOf[A any]() *ios.IOAsSlice[A] {
	return ios.NewAsSliceOf[A]()
}

func Tap[A any](f func(A) bool) *ios.IOTap[A] {
	return ios.NewTap[A](f)
}

func Foreach[A any](f func(A)) *ios.IOForeach[A] {
	return ios.NewForeach[A](f)
}

func Or[A any](f func() A) *ios.IOOr[A] {
	return ios.NewOr[A](f)
}

func FailIfEmpty[A any](f func() error) *ios.IOFailIfEmpty[A] {
	return ios.NewFailIfEmpty[A](f)
}

func OrElse[A any](f func() types.IORunnable) *ios.IOOrElse[A] {
	return ios.NewOrElse[A](f)
}

func CatchAll[A any](f func(error) *result.Result[*option.Option[A]]) *ios.IOCatchAll[A] {
	return ios.NewCatchAll[A](f)
}

func CatchAllOfResult[A any](f func(error) *result.Result[A]) *ios.IOCatchAll[A] {
	return ios.NewCatchAllOfResult[A](f)
}

func CatchAllOfOption[A any](f func(error) *option.Option[A]) *ios.IOCatchAll[A] {
	return ios.NewCatchAllOfOption[A](f)
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
