package io

import (
	"github.com/mobilemindtec/go-io/result"
	"github.com/mobilemindtec/go-io/runtime"
	"github.com/mobilemindtec/go-io/types"
)

func IO[T any](effs ...types.IOEffect) *types.IO[T] {
	return types.NewIO[T]().Effects(effs...)
}

func Attempt[T any](f func() *result.Result[T]) *types.IOAttempt[T] {
	return types.NewAttempt[T](f)
}

func Debug[A any](label string) *types.IODebug[A] {
	return types.NewDebug[A](label)
}

func FailWith[A any](f func(A) *result.Result[A]) *types.IOFailWith[A] {
	return types.NewFailWith[A](f)
}

func FailWithError[A any](f func(A) error) *types.IOFailWith[A] {
	return types.NewFailWithError[A](f)
}

func Filter[A any](f func(A) bool) *types.IOFilter[A] {
	return types.NewFilter[A](f)
}

func FlatMap[A any, B any](f func(A) *types.IO[B]) *types.IOFlatMap[A, B] {
	return types.NewFlatMap[A, B](f)
}

func Map[A any, B any](f func(A) B) *types.IOMap[A, B] {
	return types.NewMap[A, B](f)
}

func PureVal[T any](value T) *types.IOPure[T] {
	return types.NewPure[T](value)
}

func Pure[T any](f func() T) *types.IOPure[T] {
	return types.NewPureF[T](f)
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

func Runtime[T any]() *runtime.Runtime[T] {
	return runtime.NewRuntime[T]()
}
