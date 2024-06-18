package rio

import (
	"fmt"
	"github.com/mobilemindtec/go-io/option"
	"github.com/mobilemindtec/go-io/result"
	"github.com/mobilemindtec/go-io/util"
	"log"
	"runtime/debug"
)

// IO computation
type IO[T any] struct {
	value       *result.Result[*option.Option[T]]
	Computation func() *IO[T]
}

func MaybeErrorIO[T any](res result.IResult) *IO[T] {
	if res.HasError() {
		return ErrorIO[T](res.GetError())
	}
	return EmptyIO[T]()
}

func ErrorIO[T any](err error) *IO[T] {
	return &IO[T]{value: result.OfError[*option.Option[T]](err)}
}

func EmptyIO[T any]() *IO[T] {
	return &IO[T]{value: result.OfValue(option.None[T]())}
}

func NewIO[T any](value T) *IO[T] {
	return &IO[T]{value: result.OfValue(option.Of(value))}
}

func NewIOWithResult[T any](value *result.Result[*option.Option[T]]) *IO[T] {
	return &IO[T]{value: value}
}

func (this *IO[T]) String() string {
	if this.IsEmpty() {
		return "IO(empty)"
	} else if this.IsError() {
		return fmt.Sprintf("IO(%v)", this.value.Error())
	} else if this.HasValue() {
		return fmt.Sprintf("IO(%v)", this.value.Get())
	} else {
		return "IO(suspended)"
	}
}
func (this *IO[T]) HasValue() bool {
	return !this.IsError() && !this.IsEmpty()
}

func (this *IO[T]) IsEmpty() bool {
	return this.value != nil && this.value.IsOk() && this.value.Get().IsNone()
}

func (this *IO[T]) IsError() bool {
	return this.value != nil && this.value.IsError()
}

func (this *IO[T]) Get() *result.Result[*option.Option[T]] {
	return this.value
}

func (this *IO[T]) UnsafeGet() T {

	if util.IsNil(this.value) {
		panic("value is nil, probably the IO computation is not was executed")
	}

	if this.IsError() {
		panic("IO is on error state")
	}

	if this.IsEmpty() {
		panic("IO is on empty state")
	}

	return this.value.Get().Get()
}

func (this *IO[T]) UnsafeRun() *IO[T] {
	return this.Computation()
}

func suspend[T any](f func() *IO[T]) *IO[T] {
	return &IO[T]{Computation: f}
}

// Pure value
func Pure[T any](value T) *IO[T] {
	return suspend[T](func() *IO[T] {
		return NewIO(value)
	})
}

// PureF value from func
func PureF[T any](f func() T) *IO[T] {
	return suspend[T](func() *IO[T] {
		return NewIO(f())
	})
}

// Map computation
func Map[A, B any](io *IO[A], f func(A) B) *IO[B] {
	return suspend[B](func() *IO[B] {
		ref := io.UnsafeRun()
		if !ref.HasValue() {
			return MaybeErrorIO[B](ref.Get())
		}
		return NewIO(f(ref.UnsafeGet()))
	})
}

// FlatMap computation
func FlatMap[A, B any](io *IO[A], f func(A) *IO[B]) *IO[B] {
	return suspend[B](func() *IO[B] {
		ref := io.UnsafeRun()
		if !ref.HasValue() {
			return MaybeErrorIO[B](ref.Get())
		}
		return f(ref.UnsafeGet()).UnsafeRun()
	})
}

// AndThan computation
func AndThan[A, B any](io *IO[A], f func() *IO[B]) *IO[B] {
	return suspend[B](func() *IO[B] {
		ref := io.UnsafeRun()
		if !ref.HasValue() {
			return MaybeErrorIO[B](ref.Get())
		}
		return f().UnsafeRun()
	})
}

// Filter computation
func Filter[A any](io *IO[A], f func(A) bool) *IO[A] {
	return suspend[A](func() *IO[A] {
		ref := io.UnsafeRun()
		if !ref.HasValue() {
			return MaybeErrorIO[A](ref.Get())
		}
		if f(ref.UnsafeGet()) {
			return NewIO(ref.UnsafeGet())
		} else {
			return EmptyIO[A]()
		}
	})
}

// Foreach computation
func Foreach[A any](io *IO[A], f func(A)) *IO[A] {
	return suspend[A](func() *IO[A] {
		ref := io.UnsafeRun()
		if !ref.HasValue() {
			return MaybeErrorIO[A](ref.Get())
		}
		f(ref.UnsafeGet())
		return NewIO(ref.UnsafeGet())

	})
}

// OrElse computation
func OrElse[A any](io *IO[A], f func() *IO[A]) *IO[A] {
	return suspend[A](func() *IO[A] {
		ref := io.UnsafeRun()
		if !ref.HasValue() {
			return MaybeErrorIO[A](ref.Get())
		}
		if ref.IsEmpty() {
			return f().UnsafeRun()
		} else {
			return EmptyIO[A]()
		}
	})
}

// Or computation
func Or[A any](io *IO[A], f func() A) *IO[A] {
	return suspend[A](func() *IO[A] {
		ref := io.UnsafeRun()
		if !ref.HasValue() {
			return MaybeErrorIO[A](ref.Get())
		}
		if ref.IsEmpty() {
			return NewIO(f())
		} else {
			return EmptyIO[A]()
		}
	})
}

// Recover computation
func Recover[A any](io *IO[A], f func(error) A) *IO[A] {
	return suspend[A](func() *IO[A] {
		ref := io.UnsafeRun()
		if ref.IsError() {
			return NewIO(f(ref.Get().GetError()))
		}
		return NewIOWithResult(ref.Get())
	})
}

// RecoverIO computation
func RecoverIO[A any](io *IO[A], f func(error) *IO[A]) *IO[A] {
	return suspend[A](func() *IO[A] {
		ref := io.UnsafeRun()
		if ref.IsError() {
			return f(ref.Get().GetError()).UnsafeRun()
		}
		return NewIOWithResult(ref.Get())
	})
}

// CatchAll computation
func CatchAll[A any](io *IO[A], f func(error)) *IO[A] {
	return suspend[A](func() *IO[A] {
		ref := io.UnsafeRun()
		f(ref.Get().GetError())
		return NewIOWithResult(ref.Get())
	})
}

// Ensure computation
func Ensure[A any](io *IO[A], f func()) *IO[A] {
	return suspend[A](func() *IO[A] {
		ref := io.UnsafeRun()
		f()
		return NewIOWithResult(ref.Get())
	})
}

// Debug computation
func Debug[A any](io *IO[A], label ...string) *IO[A] {
	return suspend[A](func() *IO[A] {
		ref := io.UnsafeRun()
		if len(label) > 0 {
			log.Printf("DEBUG IO[%v]>> %v", label[0], ref)
		} else {
			log.Printf("DEBUG IO>> %v", ref)
		}
		return NewIOWithResult(ref.Get())
	})
}

// Attempt computation
func Attempt[A any](f func() *result.Result[A]) *IO[A] {
	return suspend[A](func() (io *IO[A]) {

		defer func() {
			if err := recover(); err != nil {

				log.Printf(string(debug.Stack()))

				switch err.(type) {
				case error:
					io = ErrorIO[A](err.(error))
					break
				default:
					io = ErrorIO[A](fmt.Errorf("%v", err))
				}
			}
		}()

		res := f()
		if res.IsOk() {
			io = NewIO(res.Get())
			return
		}
		io = MaybeErrorIO[A](res)
		return
	})
}

// Pipe2 computation
func Pipe2[A, B, T any](a *IO[A], b *IO[B], f func(A, B) *IO[T]) *IO[T] {
	return suspend[T](func() *IO[T] {
		return FlatMap(a, func(valA A) *IO[T] {
			return FlatMap(b, func(valB B) *IO[T] {
				return f(valA, valB)
			})
		}).UnsafeRun()
	})
}

// Pipe3 computation
func Pipe3[A, B, C, T any](a *IO[A], b *IO[B], c *IO[C], f func(A, B, C) *IO[T]) *IO[T] {
	return suspend[T](func() *IO[T] {
		return Pipe2(a, b, func(valA A, valB B) *IO[T] {
			return FlatMap(c, func(valC C) *IO[T] {
				return f(valA, valB, valC)
			})
		}).UnsafeRun()
	})
}

// Pipe4 computation
func Pipe4[A, B, C, D, T any](a *IO[A], b *IO[B], c *IO[C], d *IO[D], f func(A, B, C, D) *IO[T]) *IO[T] {
	return suspend[T](func() *IO[T] {
		return Pipe3(a, b, c, func(valA A, valB B, valC C) *IO[T] {
			return FlatMap(d, func(valD D) *IO[T] {
				return f(valA, valB, valC, valD)
			})
		}).UnsafeRun()
	})
}

// Pipe5 computation
func Pipe5[A, B, C, D, E, T any](a *IO[A], b *IO[B], c *IO[C], d *IO[D], e *IO[E], f func(A, B, C, D, E) *IO[T]) *IO[T] {
	return suspend[T](func() *IO[T] {
		return Pipe4(a, b, c, d, func(valA A, valB B, valC C, valD D) *IO[T] {
			return FlatMap(e, func(valE E) *IO[T] {
				return f(valA, valB, valC, valD, valE)
			})
		}).UnsafeRun()
	})
}

// Pipe6 computation
func Pipe6[A, B, C, D, E, F, T any](a *IO[A], b *IO[B], c *IO[C], d *IO[D], e *IO[E], f *IO[F], fn func(A, B, C, D, E, F) *IO[T]) *IO[T] {
	return suspend[T](func() *IO[T] {
		return Pipe5(a, b, c, d, e, func(valA A, valB B, valC C, valD D, valE E) *IO[T] {
			return FlatMap(f, func(valF F) *IO[T] {
				return fn(valA, valB, valC, valD, valE, valF)
			})
		}).UnsafeRun()
	})
}

// Pipe7 computation
func Pipe7[A, B, C, D, E, F, G, T any](a *IO[A], b *IO[B], c *IO[C], d *IO[D], e *IO[E], f *IO[F], g *IO[G], fn func(A, B, C, D, E, F, G) *IO[T]) *IO[T] {
	return suspend[T](func() *IO[T] {
		return Pipe6(a, b, c, d, e, f, func(valA A, valB B, valC C, valD D, valE E, valF F) *IO[T] {
			return FlatMap(g, func(valG G) *IO[T] {
				return fn(valA, valB, valC, valD, valE, valF, valG)
			})
		}).UnsafeRun()
	})
}

// Pipe8 computation
func Pipe8[A, B, C, D, E, F, G, H, T any](a *IO[A], b *IO[B], c *IO[C], d *IO[D], e *IO[E], f *IO[F], g *IO[G], h *IO[H], fn func(A, B, C, D, E, F, G, H) *IO[T]) *IO[T] {
	return suspend[T](func() *IO[T] {
		return Pipe7(a, b, c, d, e, f, g, func(valA A, valB B, valC C, valD D, valE E, valF F, valG G) *IO[T] {
			return FlatMap(h, func(valH H) *IO[T] {
				return fn(valA, valB, valC, valD, valE, valF, valG, valH)
			})
		}).UnsafeRun()
	})
}

// Pipe9 computation
func Pipe9[A, B, C, D, E, F, G, H, I, T any](a *IO[A], b *IO[B], c *IO[C], d *IO[D], e *IO[E], f *IO[F], g *IO[G], h *IO[H], i *IO[I], fn func(A, B, C, D, E, F, G, H, I) *IO[T]) *IO[T] {
	return suspend[T](func() *IO[T] {
		return Pipe8(a, b, c, d, e, f, g, h, func(valA A, valB B, valC C, valD D, valE E, valF F, valG G, valH H) *IO[T] {
			return FlatMap(i, func(valI I) *IO[T] {
				return fn(valA, valB, valC, valD, valE, valF, valG, valH, valI)
			})
		}).UnsafeRun()
	})
}

// Pipe10 computation
func Pipe10[A, B, C, D, E, F, G, H, I, J, T any](a *IO[A], b *IO[B], c *IO[C], d *IO[D], e *IO[E], f *IO[F], g *IO[G], h *IO[H], i *IO[I], j *IO[J], fn func(A, B, C, D, E, F, G, H, I, J) *IO[T]) *IO[T] {
	return suspend[T](func() *IO[T] {
		return Pipe9(a, b, c, d, e, f, g, h, i, func(valA A, valB B, valC C, valD D, valE E, valF F, valG G, valH H, valI I) *IO[T] {
			return FlatMap(j, func(valJ J) *IO[T] {
				return fn(valA, valB, valC, valD, valE, valF, valG, valH, valI, valJ)
			})
		}).UnsafeRun()
	})
}

// UnsafeRun run IO computations
func UnsafeRun[T any](io *IO[T]) (r *result.Result[*option.Option[T]]) {

	defer func() {
		if err := recover(); err != nil {
			switch err.(type) {
			case error:
				r = result.OfError[*option.Option[T]](err.(error))
				break
			default:
				r = result.OfError[*option.Option[T]](fmt.Errorf("%v", err))
			}
		}
	}()

	r = io.UnsafeRun().Get()
	return
}
