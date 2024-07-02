package io

import (
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

func AndThan[A any](f func() *types.IO[A]) *ios.IOAndThan[A] {
	return ios.NewAndThan(f)
}

func AndThanIO[A any](otherIO *types.IO[A]) *ios.IOAndThan[A] {
	return ios.NewAndThanIO(otherIO)
}

func Then[A any](f func(A) A) *ios.IOThen[A] {
	return ios.NewThen[A](f)
}

func Ensure[A any](f func()) *ios.IOEnsure[A] {
	return ios.NewEnsure[A](f)
}

func Debug[A any](label string) *ios.IODebug[A] {
	return ios.NewDebug[A](label)
}

func Error[A any](err error) *ios.IOError[A] {
	return ios.NewError[A](err)
}

func ErrorIO[A any](err error) *types.IO[A] {
	return ios.NewError[A](err).Lift()
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

func FlatMap[A, B any](f func(A) *types.IO[B]) *ios.IOFlatMap[A, B] {
	return ios.NewFlatMap[A, B](f)
}

func FlatMap1[A, B any](ioA *types.IO[A], f func(A) *types.IO[B]) *ios.IOFlatMap[A, B] {
	return ios.NewFlatMapIO[A, B](ioA, f)
}

func FlatMap2[A, B, T any](ioA *types.IO[A], ioB *types.IO[B], f func(A, B) *types.IO[T]) *ios.IOFlatMap2[A, B, T] {
	return ios.NewFlatMap2[A, B, T](ioA, ioB, f)
}

func FlatMap3[A, B, C, T any](ioA *types.IO[A], ioB *types.IO[B], ioC *types.IO[C], f func(A, B, C) *types.IO[T]) *ios.IOFlatMap3[A, B, C, T] {
	return ios.NewFlatMap3[A, B, C, T](ioA, ioB, ioC, f)
}

func FlatMap4[A, B, C, D, T any](ioA *types.IO[A], ioB *types.IO[B], ioC *types.IO[C], ioD *types.IO[D], f func(A, B, C, D) *types.IO[T]) *ios.IOFlatMap4[A, B, C, D, T] {
	return ios.NewFlatMap4[A, B, C, D, T](ioA, ioB, ioC, ioD, f)
}

func FlatMap5[A, B, C, D, E, T any](ioA *types.IO[A], ioB *types.IO[B], ioC *types.IO[C], ioD *types.IO[D], ioE *types.IO[E], f func(A, B, C, D, E) *types.IO[T]) *ios.IOFlatMap5[A, B, C, D, E, T] {
	return ios.NewFlatMap5[A, B, C, D, E, T](ioA, ioB, ioC, ioD, ioE, f)
}

func Map[A, B any](f func(A) B) *ios.IOMap[A, B] {
	return ios.NewMap[A, B](f)
}

func MapToUnit[A any]() *ios.IOMap[A, *types.Unit] {
	return ios.NewMap[A, *types.Unit](func(a A) *types.Unit {
		return types.OfUnit()
	})
}

func PureVal[T any](value T) *ios.IOPure[T] {
	return ios.NewPureValue[T](value)
}

func Pure[T any](f func() T) *ios.IOPure[T] {
	return ios.NewPure[T](f)
}

func RecoverPure[A any](f func(error) A) *ios.IORecover[A] {
	return ios.NewRecoverPure[A](f)
}

func Recover[A any](f func(error) *result.Result[A]) *ios.IORecover[A] {
	return ios.NewRecover[A](f)
}

func RecoverOption[A any](f func(error) *option.Option[A]) *ios.IORecover[A] {
	return ios.NewRecoverOption[A](f)
}

func RecoverResultOption[A any](f func(error) *result.Result[*option.Option[A]]) *ios.IORecover[A] {
	return ios.NewRecoverResultOption[A](f)
}

func SliceFilter[A any](f func(A) bool) *ios.IOSliceFilter[A] {
	return ios.NewSliceFilter[A](f)
}

func SliceFlatMap[A, B any](f func(A) *types.IO[B]) *ios.IOSliceFlatMap[A, B] {
	return ios.NewSliceFlatMap[A, B](f)
}

func SliceForeach[A any](f func(A)) *ios.IOSliceForeach[A] {
	return ios.NewSliceForeach[A](f)
}

func SliceAttemptOrElse[A any](f func() *result.Result[[]A]) *ios.IOSliceAttemptOrElse[A] {
	return ios.NewSliceAttemptOrElse[A](f)
}

func SliceAttemptOrElseWithState[A any](f func(*state.State) *result.Result[[]A]) *ios.IOSliceAttemptOrElse[A] {
	return ios.NewSliceAttemptOrElseWithState[A](f)
}

func SliceAttemptEach[A any](f func(A) *result.Result[A]) *ios.IOSliceAttemptEach[A] {
	return ios.NewSliceAttemptEach[A](f)
}

func SliceAttemptEachWithState[A any](f func(A, *state.State) *result.Result[A]) *ios.IOSliceAttemptEach[A] {
	return ios.NewSliceAttemptEachWithState[A](f)
}

func SliceAttempt[A any](f func([]A) *result.Result[[]A]) *ios.IOSliceAttempt[A] {
	return ios.NewSliceAttempt[A](f)
}

func SliceAttemptWithState[A any](f func([]A, *state.State) *result.Result[[]A]) *ios.IOSliceAttempt[A] {
	return ios.NewSliceAttemptWithState[A](f)
}

func SliceMap[A, B any](f func(A) B) *ios.IOSliceMap[A, B] {
	return ios.NewSliceMap[A, B](f)
}

func SliceOr[A any](f func() []A) *ios.IOSliceOr[A] {
	return ios.NewSliceOr[A](f)
}

func SliceOrElse[A any](f func() *types.IO[A]) *ios.IOSliceOrElse[A] {
	return ios.NewSliceOrElse[A](f)
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

func FailIfEmptyUnit(f func() error) *ios.IOFailIfEmpty[*types.Unit] {
	return ios.NewFailIfEmptyUnit(f)
}

func FailWith[A any](f func() error) *ios.IOFailWith[A] {
	return ios.NewFailWith[A](f)
}

func FailIf[A any](f func(A) error) *ios.IOFailIf[A] {
	return ios.NewFailIf[A](f)
}

func OrElse[A any](f func() *types.IO[A]) *ios.IOOrElse[A] {
	return ios.NewOrElse[A](f)
}

func CatchAll[A any](f func(error)) *ios.IOCatchAll[A] {
	return ios.NewCatchAll[A](f)
}

func Nohup[A any]() *ios.IONohup[A] {
	return ios.NewNohup[A]()
}

func NohupIO[A any]() *types.IO[A] {
	return ios.NewNohup[A]().Lift()
}

func Unit() *ios.IOUnit {
	return ios.NewUnit()
}

func LoadVar[A any]() *ios.IOLoadVar[A] {
	return ios.NewLoadVar[A]()
}

func Attempt[A any](f func() *result.Result[A]) *ios.IOAttempt[A] {
	return ios.NewAttempt[A](f)
}

func AttemptFlatMap[A, B any](f func(A, *state.State) *types.IO[B]) *ios.IOAttemptFlatMap[A, B] {
	return ios.NewAttemptFlatMap[A, B](f)
}

func AttemptOfOption[A any](f func() *option.Option[A]) *ios.IOAttempt[A] {
	return ios.NewAttemptOfOption[A](f)
}

func AttemptOfResultOption[A any](f func() *result.Result[*option.Option[A]]) *ios.IOAttempt[A] {
	return ios.NewAttemptOfResultOption[A](f)
}

func AttemptState[A any](f func(*state.State) *result.Result[A]) *ios.IOAttempt[A] {
	return ios.NewAttemptState[A](f)
}

func AttemptStateOfOption[A any](f func(*state.State) *option.Option[A]) *ios.IOAttempt[A] {
	return ios.NewAttemptStateOfOption[A](f)
}

func AttemptStateOfResultOption[A any](f func(*state.State) *result.Result[*option.Option[A]]) *ios.IOAttempt[A] {
	return ios.NewAttemptStateOfResultOption[A](f)
}

func AttemptOfUnit(f func()) *ios.IOAttempt[*types.Unit] {
	return ios.NewAttemptOfUnit(f)
}

func AttemptStateOfUnit(f func(*state.State)) *ios.IOAttempt[*types.Unit] {
	return ios.NewAttemptStateOfUnit(f)
}

func AttemptOfError[A any](f func() (A, error)) *ios.IOAttempt[A] {
	return ios.NewAttemptOfError[A](f)
}

func AttemptStateOfError[A any](f func(*state.State) (A, error)) *ios.IOAttempt[A] {
	return ios.NewAttemptStateOfError[A](f)
}

func AttemptValueState[A any](f func(*state.State) A) *ios.IOAttempt[A] {
	return ios.NewAttemptValueState[A](f)
}

func AttemptAndThanWithState[A any](f func(*state.State) *types.IO[A]) *ios.IOAttemptAndThan[A] {
	return ios.NewAttemptAndThanWithState[A](f)
}

func AttemptAndThan[A any](f func() *types.IO[A]) *ios.IOAttemptAndThan[A] {
	return ios.NewAttemptAndThan[A](f)
}

func AttemptRunIOWithState[A any](f func(*state.State) *types.IO[A]) *ios.IOAttemptAndThan[A] {
	return ios.NewAttemptRunIOWithState[A](f)
}

func AttemptRunIO[A any](f func() *types.IO[A]) *ios.IOAttemptAndThan[A] {
	return ios.NewAttemptRunIO[A](f)
}

func AttemptAuto[A any](f interface{}) *ios.IOAttemptAuto[A] {
	return ios.NewAttemptAuto[A](f)
}

func AttemptExec[A any](f func(A)) *ios.IOAttemptExec[A] {
	return ios.NewAttemptExec[A](f)
}

func AttemptExecWithState[A any](f func(A, *state.State)) *ios.IOAttemptExec[A] {
	return ios.NewAttemptExecWithState[A](f)
}

func AttemptExecOrElse[A any](f func()) *ios.IOAttemptExecOrElse[A] {
	return ios.NewAttemptExecOrElse[A](f)
}

func AttemptExecOrElseWithState[A any](f func(*state.State)) *ios.IOAttemptExecOrElse[A] {
	return ios.NewAttemptExecOrElseWithState[A](f)
}

func AttemptExecOrElseOfUnit(f func()) *ios.IOAttemptExecOrElse[*types.Unit] {
	return ios.NewAttemptExecOrElseOfUnit(f)
}

func AttemptExecOrElseWithStateOfUnit(f func(*state.State)) *ios.IOAttemptExecOrElse[*types.Unit] {
	return ios.NewAttemptExecOrElseWithStateOfUnit(f)
}

func AttemptOrElseWithState[A any](f func(*state.State) *types.IO[A]) *ios.IOAttemptOrElse[A] {
	return ios.NewAttemptOrElseWithState[A](f)
}

func AttemptOrElse[A any](f func() *types.IO[A]) *ios.IOAttemptOrElse[A] {
	return ios.NewAttemptOrElse[A](f)
}

func AttemptOrElseOfResult[A any](f func() *result.Result[A]) *ios.IOAttemptOrElse[A] {
	return ios.NewAttemptOrElseOfResult[A](f)
}

func AttemptOrElseOfResultWithState[A any](f func(*state.State) *result.Result[A]) *ios.IOAttemptOrElse[A] {
	return ios.NewAttemptOrElseOfResultWithState[A](f)
}

func AttemptOrElseOfResultOption[A any](f func() *result.Result[*option.Option[A]]) *ios.IOAttemptOrElse[A] {
	return ios.NewAttemptOrElseOfResultOption[A](f)
}

func AttemptOrElseOfResultOptionWithState[A any](f func(*state.State) *result.Result[*option.Option[A]]) *ios.IOAttemptOrElse[A] {
	return ios.NewAttemptOrElseOfResultOptionWithState[A](f)
}

func AttemptThen[A any](f func(A) *result.Result[A]) *ios.IOAttemptThen[A] {
	return ios.NewAttemptThen[A](f)
}

func AttemptThenWithState[A any](f func(A, *state.State) *result.Result[A]) *ios.IOAttemptThen[A] {
	return ios.NewAttemptThenWithState[A](f)
}

func AttemptThenOption[A any](f func(A) *result.Result[*option.Option[A]]) *ios.IOAttemptThen[A] {
	return ios.NewAttemptThenOption[A](f)
}

func AttemptThenOptionWithState[A any](f func(A, *state.State) *result.Result[*option.Option[A]]) *ios.IOAttemptThen[A] {
	return ios.NewAttemptThenOptionWithState[A](f)
}

func AttemptThenIO[A any](f func(A) *types.IO[A]) *ios.IOAttemptThen[A] {
	return ios.NewAttemptThenIO[A](f)
}

func AttemptThenIOWithState[A any](f func(A, *state.State) *types.IO[A]) *ios.IOAttemptThen[A] {
	return ios.NewAttemptThenIOWithState[A](f)
}

func PipeIO[A, T any](f func(A) *types.IO[T]) *ios.IOPipe[A, T] {
	return ios.NewPipeIO[A, T](f)
}

func Pipe[A, T any](f func(A) *result.Result[*option.Option[T]]) *ios.IOPipe[A, T] {
	return ios.NewPipe[A, T](f)
}

func PipeOfValue[A, T any](f func(A) T) *ios.IOPipe[A, T] {
	return ios.NewPipeOfValue[A, T](f)
}

func PipeOfResult[A, T any](f func(A) *result.Result[T]) *ios.IOPipe[A, T] {
	return ios.NewPipeOfResult[A, T](f)
}

func PipeOfOption[A, T any](f func(A) *option.Option[T]) *ios.IOPipe[A, T] {
	return ios.NewPipeOfOption[A, T](f)
}

func Pipe2IO[A, B, T any](f func(A, B) *types.IO[T]) *ios.IOPipe2[A, B, T] {
	return ios.NewPipe2IO[A, B, T](f)
}

func Pipe2[A, B, T any](f func(A, B) *result.Result[*option.Option[T]]) *ios.IOPipe2[A, B, T] {
	return ios.NewPipe2[A, B, T](f)
}

func Pipe2OfValue[A, B, T any](f func(A, B) T) *ios.IOPipe2[A, B, T] {
	return ios.NewPipe2OfValue[A, B, T](f)
}

func Pipe2OfResult[A, B, T any](f func(A, B) *result.Result[T]) *ios.IOPipe2[A, B, T] {
	return ios.NewPipe2OfResult[A, B, T](f)
}

func Pipe2OfOption[A, B, T any](f func(A, B) *option.Option[T]) *ios.IOPipe2[A, B, T] {
	return ios.NewPipe2OfOption[A, B, T](f)
}

func Pipe3IO[A, B, C, T any](f func(A, B, C) *types.IO[T]) *ios.IOPipe3[A, B, C, T] {
	return ios.NewPipe3IO[A, B, C, T](f)
}

func Pipe3[A, B, C, T any](f func(A, B, C) *result.Result[*option.Option[T]]) *ios.IOPipe3[A, B, C, T] {
	return ios.NewPipe3[A, B, C, T](f)
}

func Pipe3OfValue[A, B, C, T any](f func(A, B, C) T) *ios.IOPipe3[A, B, C, T] {
	return ios.NewPipe3OfValue[A, B, C, T](f)
}

func Pipe3OfResult[A, B, C, T any](f func(A, B, C) *result.Result[T]) *ios.IOPipe3[A, B, C, T] {
	return ios.NewPipe3OfResult[A, B, C, T](f)
}

func Pipe3OfOption[A, B, C, T any](f func(A, B, C) *option.Option[T]) *ios.IOPipe3[A, B, C, T] {
	return ios.NewPipe3OfOption[A, B, C, T](f)
}

func Pipe4IO[A, B, C, D, T any](f func(A, B, C, D) *types.IO[T]) *ios.IOPipe4[A, B, C, D, T] {
	return ios.NewPipe4IO[A, B, C, D, T](f)
}

func Pipe4[A, B, C, D, T any](f func(A, B, C, D) *result.Result[*option.Option[T]]) *ios.IOPipe4[A, B, C, D, T] {
	return ios.NewPipe4[A, B, C, D, T](f)
}

func Pipe4OfValue[A, B, C, D, T any](f func(A, B, C, D) T) *ios.IOPipe4[A, B, C, D, T] {
	return ios.NewPipe4OfValue[A, B, C, D, T](f)
}

func Pipe4OfResult[A, B, C, D, T any](f func(A, B, C, D) *result.Result[T]) *ios.IOPipe4[A, B, C, D, T] {
	return ios.NewPipe4OfResult[A, B, C, D, T](f)
}

func Pipe4OfOption[A, B, C, D, T any](f func(A, B, C, D) *option.Option[T]) *ios.IOPipe4[A, B, C, D, T] {
	return ios.NewPipe4OfOption[A, B, C, D, T](f)
}

func Pipe5IO[A, B, C, D, E, T any](f func(A, B, C, D, E) *types.IO[T]) *ios.IOPipe5[A, B, C, D, E, T] {
	return ios.NewPipe5IO[A, B, C, D, E, T](f)
}

func Pipe5[A, B, C, D, E, T any](f func(A, B, C, D, E) *result.Result[*option.Option[T]]) *ios.IOPipe5[A, B, C, D, E, T] {
	return ios.NewPipe5[A, B, C, D, E, T](f)
}

func Pipe5OfValue[A, B, C, D, E, T any](f func(A, B, C, D, E) T) *ios.IOPipe5[A, B, C, D, E, T] {
	return ios.NewPipe5OfValue[A, B, C, D, E, T](f)
}

func Pipe5OfResult[A, B, C, D, E, T any](f func(A, B, C, D, E) *result.Result[T]) *ios.IOPipe5[A, B, C, D, E, T] {
	return ios.NewPipe5OfResult[A, B, C, D, E, T](f)
}

func Pipe5OfOption[A, B, C, D, E, T any](f func(A, B, C, D, E) *option.Option[T]) *ios.IOPipe5[A, B, C, D, E, T] {
	return ios.NewPipe5OfOption[A, B, C, D, E, T](f)
}

func Pipe6IO[A, B, C, D, E, F, T any](f func(A, B, C, D, E, F) *types.IO[T]) *ios.IOPipe6[A, B, C, D, E, F, T] {
	return ios.NewPipe6IO[A, B, C, D, E, F, T](f)
}

func Pipe6[A, B, C, D, E, F, T any](f func(A, B, C, D, E, F) *result.Result[*option.Option[T]]) *ios.IOPipe6[A, B, C, D, E, F, T] {
	return ios.NewPipe6[A, B, C, D, E, F, T](f)
}

func Pipe6OfValue[A, B, C, D, E, F, T any](f func(A, B, C, D, E, F) T) *ios.IOPipe6[A, B, C, D, E, F, T] {
	return ios.NewPipe6OfValue[A, B, C, D, E, F, T](f)
}

func Pipe6OfResult[A, B, C, D, E, F, T any](f func(A, B, C, D, E, F) *result.Result[T]) *ios.IOPipe6[A, B, C, D, E, F, T] {
	return ios.NewPipe6OfResult[A, B, C, D, E, F, T](f)
}

func Pipe6OfOption[A, B, C, D, E, F, T any](f func(A, B, C, D, E, F) *option.Option[T]) *ios.IOPipe6[A, B, C, D, E, F, T] {
	return ios.NewPipe6OfOption[A, B, C, D, E, F, T](f)
}

func Pipe7IO[A, B, C, D, E, F, G, T any](f func(A, B, C, D, E, F, G) *types.IO[T]) *ios.IOPipe7[A, B, C, D, E, F, G, T] {
	return ios.NewPipe7IO[A, B, C, D, E, F, G, T](f)
}

func Pipe7[A, B, C, D, E, F, G, T any](f func(A, B, C, D, E, F, G) *result.Result[*option.Option[T]]) *ios.IOPipe7[A, B, C, D, E, F, G, T] {
	return ios.NewPipe7[A, B, C, D, E, F, G, T](f)
}

func Pipe7OfValue[A, B, C, D, E, F, G, T any](f func(A, B, C, D, E, F, G) T) *ios.IOPipe7[A, B, C, D, E, F, G, T] {
	return ios.NewPipe7OfValue[A, B, C, D, E, F, G, T](f)
}

func Pipe7OfResult[A, B, C, D, E, F, G, T any](f func(A, B, C, D, E, F, G) *result.Result[T]) *ios.IOPipe7[A, B, C, D, E, F, G, T] {
	return ios.NewPipe7OfResult[A, B, C, D, E, F, G, T](f)
}

func Pipe7OfOption[A, B, C, D, E, F, G, T any](f func(A, B, C, D, E, F, G) *option.Option[T]) *ios.IOPipe7[A, B, C, D, E, F, G, T] {
	return ios.NewPipe7OfOption[A, B, C, D, E, F, G, T](f)
}

func Pipe8IO[A, B, C, D, E, F, G, H, T any](f func(A, B, C, D, E, F, G, H) *types.IO[T]) *ios.IOPipe8[A, B, C, D, E, F, G, H, T] {
	return ios.NewPipe8IO[A, B, C, D, E, F, G, H, T](f)
}

func Pipe8[A, B, C, D, E, F, G, H, T any](f func(A, B, C, D, E, F, G, H) *result.Result[*option.Option[T]]) *ios.IOPipe8[A, B, C, D, E, F, G, H, T] {
	return ios.NewPipe8[A, B, C, D, E, F, G, H, T](f)
}

func Pipe8OfValue[A, B, C, D, E, F, G, H, T any](f func(A, B, C, D, E, F, G, H) T) *ios.IOPipe8[A, B, C, D, E, F, G, H, T] {
	return ios.NewPipe8OfValue[A, B, C, D, E, F, G, H, T](f)
}

func Pipe8OfResult[A, B, C, D, E, F, G, H, T any](f func(A, B, C, D, E, F, G, H) *result.Result[T]) *ios.IOPipe8[A, B, C, D, E, F, G, H, T] {
	return ios.NewPipe8OfResult[A, B, C, D, E, F, G, H, T](f)
}

func Pipe8OfOption[A, B, C, D, E, F, G, H, T any](f func(A, B, C, D, E, F, G, H) *option.Option[T]) *ios.IOPipe8[A, B, C, D, E, F, G, H, T] {
	return ios.NewPipe8OfOption[A, B, C, D, E, F, G, H, T](f)
}

func Pipe9IO[A, B, C, D, E, F, G, H, I, T any](f func(A, B, C, D, E, F, G, H, I) *types.IO[T]) *ios.IOPipe9[A, B, C, D, E, F, G, H, I, T] {
	return ios.NewPipe9IO[A, B, C, D, E, F, G, H, I, T](f)
}

func Pipe9[A, B, C, D, E, F, G, H, I, T any](f func(A, B, C, D, E, F, G, H, I) *result.Result[*option.Option[T]]) *ios.IOPipe9[A, B, C, D, E, F, G, H, I, T] {
	return ios.NewPipe9[A, B, C, D, E, F, G, H, I, T](f)
}

func Pipe9OfValue[A, B, C, D, E, F, G, H, I, T any](f func(A, B, C, D, E, F, G, H, I) T) *ios.IOPipe9[A, B, C, D, E, F, G, H, I, T] {
	return ios.NewPipe9OfValue[A, B, C, D, E, F, G, H, I, T](f)
}

func Pipe9OfResult[A, B, C, D, E, F, G, H, I, T any](f func(A, B, C, D, E, F, G, H, I) *result.Result[T]) *ios.IOPipe9[A, B, C, D, E, F, G, H, I, T] {
	return ios.NewPipe9OfResult[A, B, C, D, E, F, G, H, I, T](f)
}

func Pipe9OfOption[A, B, C, D, E, F, G, H, I, T any](f func(A, B, C, D, E, F, G, H, I) *option.Option[T]) *ios.IOPipe9[A, B, C, D, E, F, G, H, I, T] {
	return ios.NewPipe9OfOption[A, B, C, D, E, F, G, H, I, T](f)
}

func Pipe10IO[A, B, C, D, E, F, G, H, I, J, T any](f func(A, B, C, D, E, F, G, H, I, J) *types.IO[T]) *ios.IOPipe10[A, B, C, D, E, F, G, H, I, J, T] {
	return ios.NewPipe10IO[A, B, C, D, E, F, G, H, I, J, T](f)
}

func Pipe10[A, B, C, D, E, F, G, H, I, J, T any](f func(A, B, C, D, E, F, G, H, I, J) *result.Result[*option.Option[T]]) *ios.IOPipe10[A, B, C, D, E, F, G, H, I, J, T] {
	return ios.NewPipe10[A, B, C, D, E, F, G, H, I, J, T](f)
}

func Pipe10OfValue[A, B, C, D, E, F, G, H, I, J, T any](f func(A, B, C, D, E, F, G, H, I, J) T) *ios.IOPipe10[A, B, C, D, E, F, G, H, I, J, T] {
	return ios.NewPipe10OfValue[A, B, C, D, E, F, G, H, I, J, T](f)
}

func Pipe10OfResult[A, B, C, D, E, F, G, H, I, J, T any](f func(A, B, C, D, E, F, G, H, I, J) *result.Result[T]) *ios.IOPipe10[A, B, C, D, E, F, G, H, I, J, T] {
	return ios.NewPipe10OfResult[A, B, C, D, E, F, G, H, I, J, T](f)
}

func Pipe10OfOption[A, B, C, D, E, F, G, H, I, J, T any](f func(A, B, C, D, E, F, G, H, I, J) *option.Option[T]) *ios.IOPipe10[A, B, C, D, E, F, G, H, I, J, T] {
	return ios.NewPipe10OfOption[A, B, C, D, E, F, G, H, I, J, T](f)
}

func IOApp[T any](effects ...types.IORunnable) *runtime.IOApp[T] {
	return runtime.New[T](effects...)
}

func IOAppOfUnit(effects ...types.IORunnable) *runtime.IOApp[*types.Unit] {
	return runtime.New[*types.Unit](effects...)
}

/*
func Suspend[T any](vals ...types.IORunnable) *types.IOSuspended[T] {
	return types.NewIOSuspended[T](vals...)
}

func SuspendOfUnit(vals ...types.IORunnable) *types.IOSuspended[*types.Unit] {
	return types.NewIOSuspended[*types.Unit](vals...)
}*/

func Pipeline[T any]() *pipeline.Pipeline[T] {
	return pipeline.New[T]()
}
