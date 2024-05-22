package effect

type Value struct {
	f     func() interface{}
	value interface{}
}

func NewValue(f func() interface{}) *Value {
	return &Value{f: f}
}

func (this *Value) IsPure() bool {
	return true
}

func (this *Value) IsEffect() bool {
	return false
}

func (this *Value) Value() interface{} {
	return this.value
}

func (this *Value) GetResult() interface{} {
	return this.value
}

func (this *Value) Run() *Value {
	return &Value{f: this.f, value: this.f()}
}

func (this *Value) RunEffect() interface{} {
	return this.Run()
}

type Pure[T any] struct {
	f     func() T
	value T
}

func NewPure[T any](f func() T) *Pure[T] {
	return &Pure[T]{f: f}
}

func (this Pure[T]) Run() *Pure[T] {
	value := this.f()
	return &Pure[T]{f: this.f, value: value}
}

func (this Pure[T]) IsEffect() bool {
	return false
}

func (this Pure[T]) IsPure() bool {
	return true
}

func (this Pure[T]) GetResult() interface{} {
	return this.value
}

func (this Pure[T]) RunEffect() interface{} {
	return this.Run()
}
