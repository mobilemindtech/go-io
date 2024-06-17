package collections

import (
	"github.com/mobilemindtec/go-io/option"
)

type Stack[T any] struct {
	items []T
}

func NewStack[T any]() *Stack[T] {
	return &Stack[T]{[]T{}}
}

func (this *Stack[T]) IsEmpty() bool {
	return len(this.items) == 0
}

func (this *Stack[T]) IsNonEmpty() bool {
	return len(this.items) > 0
}

func (this *Stack[T]) Push(value T) *Stack[T] {
	this.items = append(this.items, value)
	return this
}

func (this *Stack[T]) UnsafePop() T {
	return this.Pop().OrNil()
}

func (this *Stack[T]) Pop() *option.Option[T] {
	if this.IsEmpty() {
		return option.None[T]()
	}
	size := len(this.items)
	top := this.items[size-1]
	this.items = this.items[:size-1]
	return option.Of(top)
}

func (this *Stack[T]) UnsafePeek() T {
	return this.Peek().OrNil()
}

func (this *Stack[T]) Peek() *option.Option[T] {
	if this.IsEmpty() {
		return option.None[T]()
	}
	size := len(this.items)
	top := this.items[size-1]
	return option.Of(top)
}

func (this *Stack[T]) Count() int {
	return len(this.items)
}

func (this *Stack[T]) GetItems() []T {
	return StackCopy(this).items
}

func (this *Stack[T]) Last() T {
	return this.items[0]
}

func StackCopy[T any](stack *Stack[T]) *Stack[T] {
	st := NewStack[T]()
	for _, it := range stack.items {
		st.items = append(st.items, it)
	}
	return st
}
