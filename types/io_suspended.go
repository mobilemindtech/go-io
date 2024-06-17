package types

import (
	"github.com/mobilemindtec/go-io/state"
	"reflect"
)

type IOSuspended struct {
	stack []IORunnable
}

func NewIOSuspended(vals ...IORunnable) *IOSuspended {
	s := &IOSuspended{stack: []IORunnable{}}
	return s.Suspend(vals...)
}

func (this *IOSuspended) Suspend(vals ...IORunnable) *IOSuspended {
	for _, eff := range vals {
		this.stack = append(this.stack, eff)
	}
	return this
}

func (this *IOSuspended) IOs() []IORunnable {
	return this.stack
}

// fake implements
func (this *IOSuspended) UnsafeRunIO() ResultOptionAny { return nil }
func (this *IOSuspended) GetVarName() string           { return "" }
func (this *IOSuspended) SetDebug(bool)                {}
func (this *IOSuspended) SetState(*state.State)        {}
func (this *IOSuspended) CheckTypesFlow()              {}
func (this *IOSuspended) IOType() reflect.Type         { return nil }
func (this *IOSuspended) GetLastEffect() IOEffect      { return nil }
func (this *IOSuspended) SetPrevEffect(IOEffect)       {}
