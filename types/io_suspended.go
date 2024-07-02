package types


/*
type IOSuspended[T any] struct {
	stack []IORunnable
}

func NewIOSuspended[T any](vals ...IORunnable) *IOSuspended[T] {
	s := &IOSuspended[T]{stack: []IORunnable{}}
	return s.Suspend(vals...)
}

// Suspended Add suspended IO operations linked this to last IO
func (this *IOSuspended[T]) Suspend(vals ...IORunnable) *IOSuspended[T] {
	for _, eff := range vals {
		this.stack = append(this.stack, eff)
	}
	return this
}

func (this *IOSuspended[T]) IOs() []IORunnable {
	return this.stack
}
*/
/*
// IO get last IO with suspended IO operations
func (this *IOSuspended[T]) ToIO() *IO[T] {
	l := len(this.stack)
	return this.stack[l-1].(*IO[T]).WithSuspended(this.stack[0 : l-1])
}*/

// fake implements
/*
func (this *IOSuspended[T]) UnsafeRunIO() ResultOptionAny { return nil }
func (this *IOSuspended[T]) GetVarName() string           { return "" }
func (this *IOSuspended[T]) SetDebug(bool)                {}
func (this *IOSuspended[T]) SetState(*state.State)        {}
func (this *IOSuspended[T]) CheckTypesFlow()              {}
func (this *IOSuspended[T]) IOType() reflect.Type         { return nil }
func (this *IOSuspended[T]) GetLastEffect() IOEffect      { return nil }
func (this *IOSuspended[T]) SetPrevEffect(IOEffect)       {}
func (this *IOSuspended[T]) GetSuspended() []IORunnable   { return nil }
*/