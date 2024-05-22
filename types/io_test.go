package types

import (
	"errors"
	"fmt"
	"github.com/mobilemindtec/go-io/result"
	"github.com/stretchr/testify/assert"
	"testing"
)

type Person struct {
	Name string
	Age  int
}

func TestPure(t *testing.T) {
	r :=
		NewIO[int]().
			Pure(NewPure(1)).
			UnsafeRun()

	assert.Equal(t, r.Get(), 1)
}

func TestMap(t *testing.T) {
	r :=
		NewIO[string]().
			Pure(NewPure(1)).
			Map(NewMap[int, string](func(i int) string {
				return fmt.Sprintf("value is %v", i)
			})).
			UnsafeRun()

	assert.Equal(t, r.Get(), "value is 1")

}

func TestFilter(t *testing.T) {
	f1 :=
		NewIO[*Person]().
			Pure(NewPure(&Person{Age: 20})).
			Filter(NewFilter[*Person](func(p *Person) bool {
				return p.Age > 20
			})).
			UnsafeRun()

	f2 :=
		NewIO[*Person]().
			Pure(NewPure(&Person{Age: 20})).
			Filter(NewFilter[*Person](func(p *Person) bool {
				return p.Age > 10
			})).
			UnsafeRun()

	assert.True(t, f1.IsOk() && f1.OptionEmpty())
	assert.True(t, f2.IsOk() && f2.OptionNonEmpty())
	assert.Equal(t, f2.Get(), &Person{Age: 20})
}

func TestFlatMap(t *testing.T) {
	r :=
		NewIO[string]().
			Pure(NewPure(1)).
			FlatMap(NewFlatMap[int, string](func(i int) *IO[string] {
				return NewIO[string]().Pure(NewPure(fmt.Sprintf("value is %v", i)))
			})).
			UnsafeRun()

	assert.Equal(t, r.Get(), "value is 1")

}

func TestAttempt(t *testing.T) {
	r1 :=
		NewIO[string]().
			Pure(NewPure(1)).
			FlatMap(NewFlatMap[int, string](func(i int) *IO[string] {
				return NewIO[string]().
					Attempt(NewAttempt[string](func() *result.Result[string] {
						return result.OfValue(fmt.Sprintf("success %v", i))
					}))
			})).
			UnsafeRun()

	r2 :=
		NewIO[string]().
			Pure(NewPure(1)).
			FlatMap(NewFlatMap[int, string](func(i int) *IO[string] {
				return NewIO[string]().
					Attempt(NewAttemptTry[string](func() (string, error) {
						return fmt.Sprintf("success %v", 1), nil
					}))
			})).
			UnsafeRun()

	r3 :=
		NewIO[string]().
			Pure(NewPure(1)).
			FlatMap(NewFlatMap[int, string](func(i int) *IO[string] {
				return NewIO[string]().
					Attempt(NewAttemptTry[string](func() (string, error) {
						return "", errors.New("ERROR!")
					}))
			})).
			UnsafeRun()

	assert.Equal(t, r1.Get(), "success 1")
	assert.Equal(t, r2.Get(), "success 1")
	assert.True(t, r3.IsError())
	assert.Equal(t, r3.Error().Error(), "ERROR!")
}

func TestRecover(t *testing.T) {

	r :=
		NewIO[string]().
			Pure(NewPure(1)).
			FlatMap(NewFlatMap[int, string](func(i int) *IO[string] {
				return NewIO[string]().
					Attempt(NewAttemptTry[string](func() (string, error) {
						return "", errors.New("ERROR!")
					}))
			})).
			Recover(NewRecover[string](func(_ error) string {
				return "recovered"
			})).
			UnsafeRun()

	assert.Equal(t, r.Get(), "recovered")
}

func TestFailWith(t *testing.T) {

	r1 :=
		NewIO[int]().
			Pure(NewPure(1)).
			FailWith(NewFailWith[int](func(i int) *result.Result[int] {
				return result.OfValue(i)
			})).
			UnsafeRun()

	r2 :=
		NewIO[int]().
			Pure(NewPure(1)).
			FailWith(NewFailWith[int](func(i int) *result.Result[int] {
				return result.OfError[int](errors.New("ERROR!"))
			})).
			UnsafeRun()

	r3 :=
		NewIO[int]().
			Pure(NewPure(1)).
			FailWith(NewFailWithError[int](func(i int) error {
				return errors.New("ERROR!")
			})).
			UnsafeRun()

	r4 :=
		NewIO[int]().
			Pure(NewPure(1)).
			FailWith(NewFailWithError[int](func(i int) error {
				return nil
			})).
			UnsafeRun()

	assert.Equal(t, r1.Get(), 1)
	assert.Equal(t, r2.Error().Error(), "ERROR!")
	assert.Equal(t, r3.Error().Error(), "ERROR!")
	assert.Equal(t, r4.Get(), 1)
}
