package test

import (
	"github.com/mobilemindtec/go-io/either"
	"github.com/mobilemindtec/go-io/io"
	"github.com/mobilemindtec/go-io/option"
	"github.com/mobilemindtec/go-io/result"
	"github.com/mobilemindtec/go-io/state"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRuntime(t *testing.T) {

	rt := io.Runtime[PersonPtr]()

	rt.Effects(
		io.IO[string]().
			Pure(io.PureVal[PersonPtr](&Person{Name: "Ricardo"})).
			Map(io.Map[PersonPtr, string](
				func(ptr PersonPtr) string {
					return ptr.Name
				})),
		io.IO[PersonPtr]().
			Attempt(io.AttemptState(
				func(st *state.State) *result.Result[PersonPtr] {
					name := state.Var[string](st)
					return result.OfValue(&Person{Name: name})
				})),
	).Debug().UnsafeRun()

	assert.Equal(t, "Ricardo", rt.UnsafeYield().Name)
}

func TestAttemptEither(t *testing.T) {

	rt := io.Runtime[PersonPtr]()

	res :=
		rt.Effects(
			io.IO[PersonPtr]().
				Pure(io.PureVal[PersonPtr](new(Person))),
			io.IO[PersonValidator]().
				Attempt(io.AttemptStateOfResultEither(
					func(st *state.State) *result.Result[PersonValidator] {
						person := state.Var[PersonPtr](st)

						if len(person.Name) == 0 {
							return result.OfValue(either.Left[*Validation, PersonPtr](
								ValidationWith("name can't be empty")))
						}
						return result.OfValue(either.Right[*Validation, PersonPtr](person))
					})).
				MaybeFail(io.MaybeFailError(func(a PersonValidator) error {
					if a.IsLeft() {
						return a.Left()
					}
					return nil
				})).
				Map(io.Map[PersonValidator, PersonPtr](func(validator PersonValidator) PersonPtr {
					return validator.Right()
				})),
		).UnsafeRun()

	assert.True(t, res.IsError())
	assert.Equal(t, res.Error(), "validation error")

	validation, isValidation := res.Failure().(*Validation)

	assert.True(t, isValidation)
	assert.Equal(t, len(validation.Messages), 1)
	assert.Equal(t, validation.Messages[0], "name can't be empty")
}

func TestAttemptCustomError(t *testing.T) {

	rt := io.Runtime[PersonPtr]()

	res :=
		rt.Effects(
			io.IO[PersonPtr]().
				Pure(io.PureVal[PersonPtr](new(Person))),
			io.IO[*Validation]().
				Attempt(io.AttemptState(
					func(st *state.State) *result.Result[*Validation] {
						person := state.Var[PersonPtr](st)
						if len(person.Name) == 0 {
							return result.OfValue(
								ValidationWith("name can't be empty"))
						}
						return result.OfValue(NewValidation())
					})).
				MaybeFail(io.MaybeFailError(func(a *Validation) error {
					if a.NonEmpty() {
						return a
					}
					return nil
				})).
				Attempt(io.AttemptState[PersonPtr](func(s *state.State) *result.Result[PersonPtr] {
					return result.OfValue(state.Var[PersonPtr](s))
				})),
		).UnsafeRun()

	assert.True(t, res.IsError())
	assert.Equal(t, res.Error(), "validation error")

	validation, isValidation := res.Failure().(*Validation)

	assert.True(t, isValidation)
	assert.Equal(t, len(validation.Messages), 1)
	assert.Equal(t, validation.Messages[0], "name can't be empty")
}

func TestAttemptValidationOk(t *testing.T) {

	rt := io.Runtime[PersonPtr]()

	res :=
		rt.Effects(
			io.IO[PersonPtr]().
				Pure(io.PureVal[PersonPtr](&Person{Name: "Ricardo"})),
			io.IO[PersonPtr]().
				Attempt(io.AttemptState(
					func(st *state.State) *result.Result[*Validation] {
						person := state.Var[PersonPtr](st)
						if len(person.Name) == 0 {
							return result.OfValue(
								ValidationWith("name can't be empty"))
						}
						return result.OfValue(NewValidation())
					})).
				MaybeFail(io.MaybeFailError(func(a *Validation) error {
					if a.NonEmpty() {
						return a
					}
					return nil
				})).
				Attempt(io.AttemptState[PersonPtr](func(s *state.State) *result.Result[PersonPtr] {
					return result.OfValue(state.Var[PersonPtr](s))
				})),
		).Debug().UnsafeRun()

	assert.False(t, res.IsError())
	assert.Equal(t, res.Get().Get().Name, "Ricardo")
}

func TestAttemptAutoError(t *testing.T) {

	ret :=
		io.Runtime[int]().
			Effects(
				io.IO[int]().Pure(io.PureVal(1)),
				io.IO[int]().Pure(io.PureVal(2)),
				io.IO[int]().
					Attempt(io.AttemptAuto[int](func(x int, y int) (int, error) {
						return x + y, nil
					})),
			).UnsafeRun()

	assert.Equal(t, option.Some(3), ret.Get())
}

func TestAttemptAutoOption(t *testing.T) {

	ret :=
		io.Runtime[int]().
			Effects(
				io.IO[int]().Pure(io.PureVal(1)),
				io.IO[int]().Pure(io.PureVal(2)),
				io.IO[int]().
					Attempt(io.AttemptAuto[int](func(x int, y int) *option.Option[int] {
						return option.Some(x + y)
					})),
			).UnsafeRun()

	assert.Equal(t, option.Some(3), ret.Get())
}

func TestAttemptAutoResult(t *testing.T) {

	ret :=
		io.Runtime[int]().
			Effects(
				io.IO[int]().Pure(io.PureVal(1)),
				io.IO[int]().Pure(io.PureVal(2)),
				io.IO[int]().
					Attempt(io.AttemptAuto[int](func(x int, y int) *result.Result[int] {
						return result.OfValue(x + y)
					})),
			).UnsafeRun()

	assert.Equal(t, option.Some(3), ret.Get())
}

func TestAttemptAutoResultOption(t *testing.T) {
	ret :=
		io.Runtime[int]().
			Effects(
				io.IO[int]().Pure(io.PureVal(1)),
				io.IO[int]().Pure(io.PureVal(2)),
				io.IO[int]().
					Attempt(io.AttemptAuto[int](func(x int, y int) *result.Result[*option.Option[int]] {
						return result.OfValue(option.Some(x + y))
					})),
			).UnsafeRun()

	assert.Equal(t, option.Some(3), ret.Get())
}

func TestAttemptAutoEither(t *testing.T) {
	ret :=
		io.Runtime[int]().
			Effects(
				io.IO[int]().Pure(io.PureVal(1)),
				io.IO[int]().Pure(io.PureVal(2)),
				io.IO[int]().
					Attempt(io.AttemptAuto[int](func(x int, y int) *either.Either[error, int] {
						return either.Right[error, int](x + y)
					})).
					Map(io.Map[*either.Either[error, int], int](func(e *either.Either[error, int]) int {
						return e.Right()
					})),
			).Debug().UnsafeRun()

	assert.Equal(t, option.Some(3), ret.Get())
}
