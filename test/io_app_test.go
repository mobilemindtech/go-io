package test

import (
	"github.com/mobilemindtech/go-io/either"
	"github.com/mobilemindtech/go-io/http"
	"github.com/mobilemindtech/go-io/io"
	"github.com/mobilemindtech/go-io/option"
	"github.com/mobilemindtech/go-io/result"
	"github.com/mobilemindtech/go-io/state"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIOApp(t *testing.T) {

	rt := io.IOApp[PersonPtr]()

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
	).UnsafeRun()

	assert.Equal(t, "Ricardo", rt.UnsafeYield().Name)
}

func TestAttemptEither(t *testing.T) {

	rt := io.IOApp[PersonPtr]()

	res :=
		rt.Effects(
			io.IO[PersonPtr]().
				Pure(io.PureVal[PersonPtr](new(Person))),
			io.IO[PersonValidator]().
				Attempt(io.AttemptState[PersonValidator](
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

	rt := io.IOApp[PersonPtr]()

	res :=
		rt.Effects(
			io.IO[PersonPtr]().
				Pure(io.PureVal[PersonPtr](new(Person))),
			io.IO[PersonPtr]().
				Pipe(io.PipeOfValue[PersonPtr, *Validation](
					func(person PersonPtr) *Validation {
						if len(person.Name) == 0 {
							return ValidationWith("name can't be empty")
						}
						return NewValidation()
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

	rt := io.IOApp[PersonPtr]()

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
		).UnsafeRun()

	assert.False(t, res.IsError())
	assert.Equal(t, res.Get().Get().Name, "Ricardo")
}

func TestAttemptAutoError(t *testing.T) {

	ret :=
		io.IOApp[int]().
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
		io.IOApp[int]().
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
		io.IOApp[int]().
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
		io.IOApp[int]().
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
		io.IOApp[int]().
			Effects(
				io.IO[int]().Pure(io.PureVal(1)),
				io.IO[int]().Pure(io.PureVal(2)),
				io.IO[int]().
					Attempt(io.AttemptAuto[*either.Either[error, int]](func(x int, y int) *either.Either[error, int] {
						return either.Right[error, int](x + y)
					})).
					Map(io.Map[*either.Either[error, int], int](func(e *either.Either[error, int]) int {
						return e.Right()
					})),
			).UnsafeRun()

	assert.Equal(t, option.Some(3), ret.Get())
}

func TestHttpIO(t *testing.T) {

	app := io.IOApp[int](
		http.
			NewClient[any, *User, any]().
			AsJSON().
			GetIO("http://4gym.com.br/tools-api/mock/user").
			Filter(io.Filter[*http.Response[*User, any]](
				func(r *http.Response[*User, any]) bool {
					return r.StatusCode == 200
				})),

		io.Map[*http.Response[*User, any], int](
			func(r *http.Response[*User, any]) int {
				return r.StatusCode
			}).Lift(),
	)

	assert.Equal(t, 200, app.UnsafeRun().Get().Get())
}
