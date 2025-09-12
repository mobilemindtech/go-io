package either

import (
	"fmt"
	"github.com/mobilemindtech/go-io/option"
	"github.com/mobilemindtech/go-io/result"
)

type IEither interface {
	IsLeft() bool
	IsRight() bool
	LeftAny() interface{}
	RightAny() interface{}
	String() string
	Error() string
}

type _Left[T error] struct {
	value T
}

func _newLeft[T error](value T) *_Left[T] {
	return &_Left[T]{value: value}
}

func (this _Left[T]) Get() T {
	return this.value
}

func (this _Left[T]) String() string {
	return fmt.Sprintf("Left(%v)", this.value)
}

type _Right[T any] struct {
	value T
}

func _newRight[T any](value T) *_Right[T] {
	return &_Right[T]{value: value}
}

func (this _Right[T]) Get() T {
	return this.value
}

func (this _Right[T]) String() string {
	return fmt.Sprintf("Left(%v)", this.value)
}

type Either[A error, B any] struct {
	left  *_Left[A]
	right *_Right[B]
}

func Left[A error, B any](value A) *Either[A, B] {
	return &Either[A, B]{left: _newLeft(value)}
}

func Right[A error, B any](value B) *Either[A, B] {
	return &Either[A, B]{right: _newRight(value)}
}

func (this Either[A, B]) Left() A {
	if this.IsRight() {
		panic("invalid call Left of Right")
	}
	return this.left.Get()
}
func (this Either[A, B]) Right() B {
	if this.IsLeft() {
		panic("invalid call GetRight of Left")
	}
	return this.right.Get()
}

func (this Either[A, B]) IsLeft() bool {
	return this.left != nil
}

func (this Either[A, B]) IsRight() bool {
	return !this.IsLeft()
}

func (this Either[A, B]) LeftAny() interface{} {
	return this.Left()
}

func (this Either[A, B]) RightAny() interface{} {
	return this.Right()
}

func (this Either[A, B]) Error() string {
	return this.Left().Error()
}

func (this Either[A, B]) String() string {
	if this.IsLeft() {
		return this.left.String()
	}
	return this.right.String()
}

func (this Either[A, B]) IfLeft(f func(A)) Either[A, B] {
	if this.IsLeft() {
		f(this.Left())
	}
	return this
}

func (this Either[A, B]) IfRight(f func(B)) Either[A, B] {
	if this.IsRight() {
		f(this.Right())
	}
	return this
}

func (this Either[A, B]) ToLeftOption() *option.Option[A] {
	if this.IsLeft() {
		return option.Some(this.Left())
	}
	return option.None[A]()
}

func (this Either[A, B]) ToRightOption() *option.Option[B] {
	if this.IsRight() {
		return option.Some(this.Right())
	}
	return option.None[B]()
}

type EitherE[A any] struct {
	*Either[error, A]
}

func (this EitherE[A]) ToResult() *result.Result[A] {
	if this.IsLeft() {
		return result.OfError[A](this.Left())
	}
	return result.OfValue(this.Right())
}

func (this *EitherE[A]) Error() string {
	return this.Left().Error()
}

func LeftE[A any](value error) *EitherE[A] {
	return &EitherE[A]{Left[error, A](value)}
}

func RightE[A any](value A) *EitherE[A] {
	return &EitherE[A]{Right[error, A](value)}
}
