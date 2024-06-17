package test

import (
	"errors"
	"fmt"
	"github.com/mobilemindtec/go-io/io"
	"github.com/mobilemindtec/go-io/option"
	"github.com/mobilemindtec/go-io/result"
	"github.com/mobilemindtec/go-io/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPure(t *testing.T) {
	r :=
		types.NewIO[int]().
			Pure(io.PureVal(1)).
			UnsafeRun()

	assert.Equal(t, option.Some(1), r.Get())
}

func TestMap(t *testing.T) {
	r :=
		types.NewIO[string]().
			Pure(io.PureVal(1)).
			Map(io.Map[int, string](func(i int) string {
				return fmt.Sprintf("value is %v", i)
			})).
			UnsafeRun()

	assert.Equal(t, option.Some("value is 1"), r.Get())

}

func TestFilter(t *testing.T) {
	f1 :=
		types.NewIO[*Person]().
			Pure(io.PureVal(&Person{Age: 20})).
			Filter(io.Filter[*Person](func(p *Person) bool {
				return p.Age > 20
			})).
			UnsafeRun()

	f2 :=
		types.NewIO[*Person]().
			Pure(io.PureVal(&Person{Age: 20})).
			Filter(io.Filter[*Person](func(p *Person) bool {
				return p.Age > 10
			})).
			UnsafeRun()

	assert.True(t, f1.IsOk() && f1.Get().IsEmpty())
	assert.True(t, f2.IsOk() && f2.Get().NonEmpty())
	assert.Equal(t, option.Some(&Person{Age: 20}), f2.Get())
}

func TestFlatMap(t *testing.T) {
	r :=
		types.NewIO[string]().
			Pure(io.PureVal(1)).
			FlatMap(io.FlatMap[int, string](func(i int) types.IORunnable {
				return types.NewIO[string]().Pure(io.PureVal(fmt.Sprintf("value is %v", i)))
			})).
			UnsafeRun()

	assert.Equal(t, option.Some("value is 1"), r.Get())

}

func TestAttempt(t *testing.T) {
	r1 :=
		types.NewIO[string]().
			Pure(io.PureVal(1)).
			FlatMap(io.FlatMap[int, string](func(i int) types.IORunnable {
				return types.NewIO[string]().
					Attempt(io.Attempt[string](func() *result.Result[string] {
						return result.OfValue(fmt.Sprintf("success %v", i))
					}))
			})).UnsafeRun()

	r2 :=
		types.NewIO[string]().
			Pure(io.PureVal(1)).
			FlatMap(io.FlatMap[int, string](func(i int) types.IORunnable {
				return types.NewIO[string]().
					Attempt(io.AttemptOfError[string](func() (string, error) {
						return fmt.Sprintf("success %v", 1), nil
					}))
			})).
			UnsafeRun()

	r3 :=
		types.NewIO[string]().
			Pure(io.PureVal(1)).
			FlatMap(io.FlatMap[int, string](func(i int) types.IORunnable {
				return types.NewIO[string]().
					Attempt(io.AttemptOfError[string](func() (string, error) {
						return "", errors.New("ERROR!")
					}))
			})).
			UnsafeRun()

	assert.Equal(t, option.Some("success 1"), r1.Get())
	assert.Equal(t, option.Some("success 1"), r2.Get())
	assert.True(t, r3.IsError())
	assert.Equal(t, r3.Error(), "ERROR!")
}

func TestRecover(t *testing.T) {

	r :=
		types.NewIO[string]().
			Pure(io.PureVal(1)).
			FlatMap(io.FlatMap[int, string](func(i int) types.IORunnable {
				return types.NewIO[string]().
					Attempt(io.AttemptOfError[string](func() (string, error) {
						return "", errors.New("ERROR!")
					}))
			})).
			Recover(io.RecoverPure[string](func(_ error) string {
				return "recovered"
			})).
			UnsafeRun()

	assert.Equal(t, option.Some("recovered"), r.Get())
}

func TestFailWith(t *testing.T) {

	r1 :=
		types.NewIO[int]().
			Pure(io.PureVal(1)).
			MaybeFail(io.MaybeFail[int](func(i int) *result.Result[int] {
				return result.OfValue(i)
			})).
			UnsafeRun()

	r2 :=
		types.NewIO[int]().
			Pure(io.PureVal(1)).
			MaybeFail(io.MaybeFail[int](func(i int) *result.Result[int] {
				return result.OfError[int](errors.New("ERROR!"))
			})).
			UnsafeRun()

	r3 :=
		types.NewIO[int]().
			Pure(io.PureVal(1)).
			MaybeFail(io.MaybeFailError[int](func(i int) error {
				return errors.New("ERROR!")
			})).
			UnsafeRun()

	r4 :=
		types.NewIO[int]().
			Pure(io.PureVal(1)).
			MaybeFail(io.MaybeFailError[int](func(i int) error {
				return nil
			})).
			UnsafeRun()

	assert.Equal(t, option.Some(1), r1.Get())
	assert.Equal(t, r2.Error(), "ERROR!")
	assert.Equal(t, r3.Error(), "ERROR!")
	assert.Equal(t, option.Some(1), r4.Get())
}
