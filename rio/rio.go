package rio

import (
	"fmt"
	"github.com/mobilemindtec/go-io/option"
	"github.com/mobilemindtec/go-io/result"
	"github.com/mobilemindtec/go-io/types"
	"github.com/mobilemindtec/go-io/util"
	"log"
	"reflect"
	"runtime"
	"runtime/debug"
	"strings"
)

type RIOError struct {
	Message    string
	StackTrace string
	DebugInfo  string
	IOName     string
}

func NewRIOError(message string, stacktrace []byte) *RIOError {
	return &RIOError{Message: message, StackTrace: string(stacktrace)}
}

func (this RIOError) Error() string {
	return this.Message
}

type RIO interface {
	UnsafeRunIO() *result.Result[*option.Option[any]]
}

// IO computation
type IO[T any] struct {
	value       *result.Result[*option.Option[T]]
	debug_      bool
	name        string
	debugInfo   string
	computation func(*IO[T]) *IO[T]
}

func NewMaybeErrorIO[T any](res result.IResult) *IO[T] {
	if res.HasError() {
		return NewErrorIO[T](res.GetError())
	}
	return NewEmptyIO[T]()
}

func NewErrorIO[T any](err error) *IO[T] {
	return &IO[T]{value: result.OfError[*option.Option[T]](err)}
}

func NewEmptyIO[T any]() *IO[T] {
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
	} else if !this.IsEmpty() {
		return fmt.Sprintf("IO(%v)", this.value.Get())
	} else {
		return "IO(suspended)"
	}
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

func (this *IO[T]) Debug() *IO[T] {
	_, filename, line, _ := runtime.Caller(1)
	this.debug_ = true
	return this.withDebugInfo(filename, line)
}

func (this *IO[T]) As(name string) *IO[T] {
	this.name = name
	return this
}

func (this *IO[T]) withDebugInfo(filename string, lineNumber int) *IO[T] {
	this.debugInfo = fmt.Sprintf("add in %v:%v",
		getFileName(filename), lineNumber)
	return this
}

func (this *IO[T]) UnsafeRun() *IO[T] {

	if this.debug_ {
		_, filename, line, _ := runtime.Caller(1)
		log.Printf(">> DEBUG IO(%v)[%v] %v, call in %v:%v\n",
			this.name, reflect.TypeFor[T]().String(), this.debugInfo, getFileName(filename), line)
	}

	defer func() {
		if err := recover(); err != nil {
			log.Printf(">> DEBUG IO(%v)[%v] error: %v \n", this.name, reflect.TypeFor[T]().String(), err)
		}
	}()

	if this.computation != nil {
		return this.computation(this)
	}

	return this
}

func (this *IO[T]) UnsafeRunIO() *result.Result[*option.Option[any]] {
	return this.UnsafeRun().Get().ToResultOfOption()
}

func (this *IO[T]) PerformIO() *result.Result[*option.Option[T]] {
	return this.UnsafeRun().Get()
}

func suspend[T any](f func(*IO[T]) *IO[T]) *IO[T] {
	return &IO[T]{computation: f}

}

// Pure value
func Pure[T any](value T) *IO[T] {
	return suspend(func(_ *IO[T]) *IO[T] {
		return NewIO(value)
	}).As("Pure")
}

// PureF value from func
func PureF[T any](f func() T) *IO[T] {
	return suspend(func(_ *IO[T]) *IO[T] {
		return NewIO(f())
	}).As("PureF")
}

// Map computation
func Map[A, B any](io *IO[A], f func(A) B) *IO[B] {
	return suspend(func(_ *IO[B]) *IO[B] {
		ref := io.UnsafeRun()
		if ref.IsError() || ref.IsEmpty() {
			return NewMaybeErrorIO[B](ref.Get())
		}
		return NewIO(f(ref.UnsafeGet()))
	}).As("Map")
}

func MapToUnit[A any](io *IO[A]) *IO[*types.Unit] {
	return suspend(func(_ *IO[*types.Unit]) *IO[*types.Unit] {
		ref := io.UnsafeRun()
		if ref.IsError() || ref.IsEmpty() {
			return NewMaybeErrorIO[*types.Unit](ref.Get())
		}
		return NewIO(types.OfUnit())
	}).As("MapToUnit")
}

// FlatMap computation
func FlatMap[A, B any](io *IO[A], f func(A) *IO[B]) *IO[B] {
	return suspend(func(_ *IO[B]) *IO[B] {
		ref := io.UnsafeRun()
		if ref.IsError() || ref.IsEmpty() {
			return NewMaybeErrorIO[B](ref.Get())
		}
		return f(ref.UnsafeGet()).UnsafeRun()
	}).As("FlatMap")
}

// AndThan computation
func AndThan[A, B any](io *IO[A], f func() *IO[B]) *IO[B] {
	return suspend(func(_ *IO[B]) *IO[B] {
		ref := io.UnsafeRun()
		if ref.IsError() || ref.IsEmpty() {
			return NewMaybeErrorIO[B](ref.Get())
		}
		return f().UnsafeRun()
	}).As("AndThan")
}

func AndThanIO[A, B any](ioA *IO[A], ioB *IO[B]) *IO[B] {
	return suspend(func(_ *IO[B]) *IO[B] {
		ioA.UnsafeRun()
		return ioB.UnsafeRun()
	}).As("AndThanIO")
}

// Filter computation
func Filter[A any](io *IO[A], f func(A) bool) *IO[A] {
	return suspend(func(_ *IO[A]) *IO[A] {
		ref := io.UnsafeRun()
		if ref.IsError() || ref.IsEmpty() {
			return NewMaybeErrorIO[A](ref.Get())
		}
		if f(ref.UnsafeGet()) {
			return NewIO(ref.UnsafeGet())
		} else {
			return NewEmptyIO[A]()
		}
	}).As("Filter")
}

// Foreach computation
func Foreach[A any](io *IO[A], f func(A)) *IO[A] {
	return suspend(func(_ *IO[A]) *IO[A] {
		ref := io.UnsafeRun()
		if ref.IsEmpty() || ref.IsEmpty() {
			return NewMaybeErrorIO[A](ref.Get())
		}
		f(ref.UnsafeGet())
		return NewIO(ref.UnsafeGet())

	}).As("Foreach")
}

// OrElse computation
func OrElse[A any](io *IO[A], f func() *IO[A]) *IO[A] {
	return suspend(func(_ *IO[A]) *IO[A] {
		ref := io.UnsafeRun()
		if ref.IsError() {
			return NewErrorIO[A](ref.Get().Failure())
		}
		if ref.IsEmpty() {
			return f().UnsafeRun()
		} else {
			return NewIO(ref.UnsafeGet())
		}
	}).As("OrElse")
}

// Or computation
func Or[A any](io *IO[A], f func() A) *IO[A] {
	return suspend(func(_ *IO[A]) *IO[A] {
		ref := io.UnsafeRun()
		if ref.IsError() {
			return NewErrorIO[A](ref.Get().Failure())
		}
		if ref.IsEmpty() {
			return NewIO(f())
		} else {
			return NewIO(ref.UnsafeGet())
		}
	}).As("Or")
}

func IfEmpty[A any](io *IO[A], f func()) *IO[A] {
	return suspend(func(_ *IO[A]) *IO[A] {
		ref := io.UnsafeRun()
		if ref.IsError() {
			return NewErrorIO[A](ref.Get().Failure())
		}
		if ref.IsEmpty() {
			f()
		}
		return NewEmptyIO[A]()
	}).As("IfEmpty")
}

// Recover computation
func Recover[A any](io *IO[A], f func(error) A) *IO[A] {
	return suspend(func(_ *IO[A]) *IO[A] {
		ref := io.UnsafeRun()
		if ref.IsError() {
			return NewIO(f(ref.Get().GetError()))
		}
		return NewIOWithResult(ref.Get())
	}).As("Recover")
}

// RecoverIO computation
func RecoverIO[A any](io *IO[A], f func(error) *IO[A]) *IO[A] {
	return suspend(func(_ *IO[A]) *IO[A] {
		ref := io.UnsafeRun()
		if ref.IsError() {
			return f(ref.Get().GetError()).UnsafeRun()
		}
		return NewIOWithResult(ref.Get())
	}).As("RecoverIO")
}

// CatchAll computation
func CatchAll[A any](io *IO[A], f func(error)) *IO[A] {
	return suspend(func(_ *IO[A]) *IO[A] {
		ref := io.UnsafeRun()
		if ref.IsError() {
			f(ref.Get().GetError())
		}
		return NewIOWithResult(ref.Get())
	}).As("CatchAll")
}

// Ensure computation
func Ensure[A any](io *IO[A], f func()) *IO[A] {
	return suspend(func(_ *IO[A]) *IO[A] {
		ref := io.UnsafeRun()
		f()
		return NewIOWithResult(ref.Get())
	}).As("Ensure")
}

// EnsureUnit
func EnsureUnit(f func()) *IO[*types.Unit] {
	return suspend(func(_ *IO[*types.Unit]) *IO[*types.Unit] {
		f()
		return NewIO(types.OfUnit())
	}).As("EnsureUnit")
}

// EnsureIO
func EnsureIO[T any](io *IO[T], f func()) *IO[T] {
	return suspend(func(_ *IO[T]) *IO[T] {
		f()
		return io.UnsafeRun()
	}).As("EnsureIO")
}

// Debug computation
func Debug[A any](io *IO[A], label ...string) *IO[A] {
	return suspend(func(_ *IO[A]) *IO[A] {
		ref := io.UnsafeRun()
		if len(label) > 0 {
			log.Printf("DEBUG IO[%v]>> %v", label[0], ref)
		} else {
			log.Printf("DEBUG IO>> %v", ref)
		}
		return NewIOWithResult(ref.Get())
	}).As("Debug")
}

// Attempt computation
func Attempt[A any](f func() *result.Result[A]) *IO[A] {
	return suspend(func(that *IO[A]) (io *IO[A]) {

		defer func() {
			if err := recover(); err != nil {
				io = catchErrorForAttempt[A](err, that)
			}
		}()

		res := f()
		if res.IsOk() {
			io = NewIO(res.Get())
			return
		}
		io = NewMaybeErrorIO[A](res)
		return
	}).As("Attempt")
}

// AttemptWith computation
func AttemptWith[A, B any](ioA *IO[A], f func(A) *result.Result[B]) *IO[B] {
	return suspend(func(that *IO[B]) (io *IO[B]) {

		defer func() {
			if err := recover(); err != nil {
				io = catchErrorForAttempt[B](err, that)
			}
		}()

		resultIO := ioA.UnsafeRun()

		if resultIO.IsError() {
			io = NewErrorIO[B](resultIO.Get().Failure())
			return
		}

		if resultIO.IsEmpty() {
			io = NewEmptyIO[B]()
			return
		}

		res := f(resultIO.UnsafeGet())
		if res.IsOk() {
			io = NewIO(res.Get())
			return
		}
		io = NewMaybeErrorIO[B](res)
		return
	}).As("AttemptWith")
}

// AttemptWith computation
func AttemptWithOption[A, B any](ioA *IO[A], f func(A) *result.Result[*option.Option[B]]) *IO[B] {
	return suspend(func(that *IO[B]) (io *IO[B]) {

		defer func() {
			if err := recover(); err != nil {
				io = catchErrorForAttempt[B](err, that)
			}
		}()

		resultIO := ioA.UnsafeRun()

		if resultIO.IsError() {
			io = NewErrorIO[B](resultIO.Get().Failure())
			return
		}

		if resultIO.IsEmpty() {
			io = NewEmptyIO[B]()
			return
		}

		io = NewIOWithResult(f(resultIO.UnsafeGet()))
		return
	}).As("AttemptWithOption")
}

// FlatMap2 computation
func FlatMap2[A, B, T any](a *IO[A], b *IO[B], f func(A, B) *IO[T]) *IO[T] {
	return suspend(func(that *IO[T]) *IO[T] {
		return FlatMap(a, func(valA A) *IO[T] {
			return FlatMap(b, func(valB B) *IO[T] {
				return f(valA, valB)
			})
		}).UnsafeRun()
	}).As("FlatMap2")
}

// FlatMap3 computation
func FlatMap3[A, B, C, T any](a *IO[A], b *IO[B], c *IO[C], f func(A, B, C) *IO[T]) *IO[T] {
	return suspend(func(_ *IO[T]) *IO[T] {
		return FlatMap2(a, b, func(valA A, valB B) *IO[T] {
			return FlatMap(c, func(valC C) *IO[T] {
				return f(valA, valB, valC)
			})
		}).UnsafeRun()
	}).As("FlatMap3")
}

// FlatMap4 computation
func FlatMap4[A, B, C, D, T any](a *IO[A], b *IO[B], c *IO[C], d *IO[D], f func(A, B, C, D) *IO[T]) *IO[T] {
	return suspend(func(_ *IO[T]) *IO[T] {
		return FlatMap3(a, b, c, func(valA A, valB B, valC C) *IO[T] {
			return FlatMap(d, func(valD D) *IO[T] {
				return f(valA, valB, valC, valD)
			})
		}).UnsafeRun()
	}).As("FlatMap4")
}

// FlatMap5 computation
func FlatMap5[A, B, C, D, E, T any](a *IO[A], b *IO[B], c *IO[C], d *IO[D], e *IO[E], f func(A, B, C, D, E) *IO[T]) *IO[T] {
	return suspend(func(_ *IO[T]) *IO[T] {
		return FlatMap4(a, b, c, d, func(valA A, valB B, valC C, valD D) *IO[T] {
			return FlatMap(e, func(valE E) *IO[T] {
				return f(valA, valB, valC, valD, valE)
			})
		}).UnsafeRun()
	}).As("FlatMap5")
}

// FlatMap6 computation
func FlatMap6[A, B, C, D, E, F, T any](a *IO[A], b *IO[B], c *IO[C], d *IO[D], e *IO[E], f *IO[F], fn func(A, B, C, D, E, F) *IO[T]) *IO[T] {
	return suspend(func(_ *IO[T]) *IO[T] {
		return FlatMap5(a, b, c, d, e, func(valA A, valB B, valC C, valD D, valE E) *IO[T] {
			return FlatMap(f, func(valF F) *IO[T] {
				return fn(valA, valB, valC, valD, valE, valF)
			})
		}).UnsafeRun()
	}).As("FlatMap6")
}

// FlatMap7 computation
func FlatMap7[A, B, C, D, E, F, G, T any](a *IO[A], b *IO[B], c *IO[C], d *IO[D], e *IO[E], f *IO[F], g *IO[G], fn func(A, B, C, D, E, F, G) *IO[T]) *IO[T] {
	return suspend(func(_ *IO[T]) *IO[T] {
		return FlatMap6(a, b, c, d, e, f, func(valA A, valB B, valC C, valD D, valE E, valF F) *IO[T] {
			return FlatMap(g, func(valG G) *IO[T] {
				return fn(valA, valB, valC, valD, valE, valF, valG)
			})
		}).UnsafeRun()
	}).As("FlatMap7")
}

// FlatMap8 computation
func FlatMap8[A, B, C, D, E, F, G, H, T any](a *IO[A], b *IO[B], c *IO[C], d *IO[D], e *IO[E], f *IO[F], g *IO[G], h *IO[H], fn func(A, B, C, D, E, F, G, H) *IO[T]) *IO[T] {
	return suspend(func(_ *IO[T]) *IO[T] {
		return FlatMap7(a, b, c, d, e, f, g, func(valA A, valB B, valC C, valD D, valE E, valF F, valG G) *IO[T] {
			return FlatMap(h, func(valH H) *IO[T] {
				return fn(valA, valB, valC, valD, valE, valF, valG, valH)
			})
		}).UnsafeRun()
	}).As("FlatMap8")
}

// FlatMap9 computation
func FlatMap9[A, B, C, D, E, F, G, H, I, T any](a *IO[A], b *IO[B], c *IO[C], d *IO[D], e *IO[E], f *IO[F], g *IO[G], h *IO[H], i *IO[I], fn func(A, B, C, D, E, F, G, H, I) *IO[T]) *IO[T] {
	return suspend(func(_ *IO[T]) *IO[T] {
		return FlatMap8(a, b, c, d, e, f, g, h, func(valA A, valB B, valC C, valD D, valE E, valF F, valG G, valH H) *IO[T] {
			return FlatMap(i, func(valI I) *IO[T] {
				return fn(valA, valB, valC, valD, valE, valF, valG, valH, valI)
			})
		}).UnsafeRun()
	}).As("FlatMap9")
}

// FlatMap10 computation
func FlatMap10[A, B, C, D, E, F, G, H, I, J, T any](a *IO[A], b *IO[B], c *IO[C], d *IO[D], e *IO[E], f *IO[F], g *IO[G], h *IO[H], i *IO[I], j *IO[J], fn func(A, B, C, D, E, F, G, H, I, J) *IO[T]) *IO[T] {
	return suspend(func(_ *IO[T]) *IO[T] {
		return FlatMap9(a, b, c, d, e, f, g, h, i, func(valA A, valB B, valC C, valD D, valE E, valF F, valG G, valH H, valI I) *IO[T] {
			return FlatMap(j, func(valJ J) *IO[T] {
				return fn(valA, valB, valC, valD, valE, valF, valG, valH, valI, valJ)
			})
		}).UnsafeRun()
	}).As("FlatMap10")
}

// Map2 computation
func Map2[A, B, T any](a *IO[A], b *IO[B], f func(A, B) T) *IO[T] {
	return suspend(func(_ *IO[T]) *IO[T] {
		return FlatMap(a, func(valA A) *IO[T] {
			return FlatMap(b, func(valB B) *IO[T] {
				return NewIO(f(valA, valB))
			})
		}).UnsafeRun()
	}).As("Map2")
}

// Map3 computation
func Map3[A, B, C, T any](a *IO[A], b *IO[B], c *IO[C], f func(A, B, C) T) *IO[T] {
	return suspend(func(_ *IO[T]) *IO[T] {
		return FlatMap2(a, b, func(valA A, valB B) *IO[T] {
			return FlatMap(c, func(valC C) *IO[T] {
				return NewIO(f(valA, valB, valC))
			})
		}).UnsafeRun()
	}).As("Map3")
}

// Map4 computation
func Map4[A, B, C, D, T any](a *IO[A], b *IO[B], c *IO[C], d *IO[D], f func(A, B, C, D) T) *IO[T] {
	return suspend(func(_ *IO[T]) *IO[T] {
		return FlatMap3(a, b, c, func(valA A, valB B, valC C) *IO[T] {
			return FlatMap(d, func(valD D) *IO[T] {
				return NewIO(f(valA, valB, valC, valD))
			})
		}).UnsafeRun()
	}).As("Map4")
}

// Map5 computation
func Map5[A, B, C, D, E, T any](a *IO[A], b *IO[B], c *IO[C], d *IO[D], e *IO[E], f func(A, B, C, D, E) T) *IO[T] {
	return suspend(func(_ *IO[T]) *IO[T] {
		return FlatMap4(a, b, c, d, func(valA A, valB B, valC C, valD D) *IO[T] {
			return FlatMap(e, func(valE E) *IO[T] {
				return NewIO(f(valA, valB, valC, valD, valE))
			})
		}).UnsafeRun()
	}).As("Map5")
}

// Map6 computation
func Map6[A, B, C, D, E, F, T any](a *IO[A], b *IO[B], c *IO[C], d *IO[D], e *IO[E], f *IO[F], fn func(A, B, C, D, E, F) T) *IO[T] {
	return suspend(func(_ *IO[T]) *IO[T] {
		return FlatMap5(a, b, c, d, e, func(valA A, valB B, valC C, valD D, valE E) *IO[T] {
			return FlatMap(f, func(valF F) *IO[T] {
				return NewIO(fn(valA, valB, valC, valD, valE, valF))
			})
		}).UnsafeRun()
	}).As("Map6")
}

// Map7 computation
func Map7[A, B, C, D, E, F, G, T any](a *IO[A], b *IO[B], c *IO[C], d *IO[D], e *IO[E], f *IO[F], g *IO[G], fn func(A, B, C, D, E, F, G) T) *IO[T] {
	return suspend(func(_ *IO[T]) *IO[T] {
		return FlatMap6(a, b, c, d, e, f, func(valA A, valB B, valC C, valD D, valE E, valF F) *IO[T] {
			return FlatMap(g, func(valG G) *IO[T] {
				return NewIO(fn(valA, valB, valC, valD, valE, valF, valG))
			})
		}).UnsafeRun()
	}).As("Map7")
}

// Map8 computation
func Map8[A, B, C, D, E, F, G, H, T any](a *IO[A], b *IO[B], c *IO[C], d *IO[D], e *IO[E], f *IO[F], g *IO[G], h *IO[H], fn func(A, B, C, D, E, F, G, H) T) *IO[T] {
	return suspend(func(_ *IO[T]) *IO[T] {
		return FlatMap7(a, b, c, d, e, f, g, func(valA A, valB B, valC C, valD D, valE E, valF F, valG G) *IO[T] {
			return FlatMap(h, func(valH H) *IO[T] {
				return NewIO(fn(valA, valB, valC, valD, valE, valF, valG, valH))
			})
		}).UnsafeRun()
	}).As("Map8")
}

// Map9 computation
func Map9[A, B, C, D, E, F, G, H, I, T any](a *IO[A], b *IO[B], c *IO[C], d *IO[D], e *IO[E], f *IO[F], g *IO[G], h *IO[H], i *IO[I], fn func(A, B, C, D, E, F, G, H, I) T) *IO[T] {
	return suspend(func(_ *IO[T]) *IO[T] {
		return FlatMap8(a, b, c, d, e, f, g, h, func(valA A, valB B, valC C, valD D, valE E, valF F, valG G, valH H) *IO[T] {
			return FlatMap(i, func(valI I) *IO[T] {
				return NewIO(fn(valA, valB, valC, valD, valE, valF, valG, valH, valI))
			})
		}).UnsafeRun()
	}).As("Map9")
}

// Map10 computation
func Map10[A, B, C, D, E, F, G, H, I, J, T any](a *IO[A], b *IO[B], c *IO[C], d *IO[D], e *IO[E], f *IO[F], g *IO[G], h *IO[H], i *IO[I], j *IO[J], fn func(A, B, C, D, E, F, G, H, I, J) T) *IO[T] {
	return suspend(func(_ *IO[T]) *IO[T] {
		return FlatMap9(a, b, c, d, e, f, g, h, i, func(valA A, valB B, valC C, valD D, valE E, valF F, valG G, valH H, valI I) *IO[T] {
			return FlatMap(j, func(valJ J) *IO[T] {
				return NewIO(fn(valA, valB, valC, valD, valE, valF, valG, valH, valI, valJ))
			})
		}).UnsafeRun()
	}).As("Map10")
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

func UnwrapOption[T any](opt *option.Option[T]) T {
	return opt.Get()
}

func OfUnit() *types.Unit {
	return types.OfUnit()
}

func getFileName(name string) string {
	sp := strings.Split(name, "/")
	return sp[len(sp)-1]
}

func catchErrorForAttempt[A any](err any, io *IO[A]) *IO[A] {

	stacktrace := string(debug.Stack())

	if io.debug_ {
		log.Printf(">> DEBUG IO(%v)[%v] %v\n",
			io.name, reflect.TypeFor[A]().String(), io.debugInfo)
		log.Printf(">> DEBUG IO(%v)\n\n%v\n\n", io.name, stacktrace)
	}

	rioError := &RIOError{
		Message:    fmt.Sprintf("%v", err),
		StackTrace: stacktrace,
		DebugInfo:  io.debugInfo,
		IOName:     io.name,
	}

	return NewErrorIO[A](rioError)
}
