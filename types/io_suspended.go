package types

import (
	"github.com/mobilemindtec/go-io/state"
	"reflect"
)

type IOSuspended[T any] struct {
	stack []IORunnable
}

func NewIOSuspended[T any](vals ...IORunnable) *IOSuspended[T] {
	s := &IOSuspended[T]{stack: []IORunnable{}}
	return s.Suspend(vals...)
}

func (this *IOSuspended[T]) Suspend(vals ...IORunnable) *IOSuspended[T] {
	for _, eff := range vals {
		this.stack = append(this.stack, eff)
	}
	return this
}

func (this *IOSuspended[T]) IOs() []IORunnable {
	return this.stack
}

// fake implements
func (this *IOSuspended[T]) UnsafeRunIO() ResultOptionAny { return nil }
func (this *IOSuspended[T]) GetVarName() string           { return "" }
func (this *IOSuspended[T]) SetDebug(bool)                {}
func (this *IOSuspended[T]) SetState(*state.State)        {}
func (this *IOSuspended[T]) CheckTypesFlow()              {}
func (this *IOSuspended[T]) IOType() reflect.Type         { return nil }
func (this *IOSuspended[T]) GetLastEffect() IOEffect      { return nil }
func (this *IOSuspended[T]) SetPrevEffect(IOEffect)       {}
