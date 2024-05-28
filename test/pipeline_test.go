package test

import (
	"errors"
	"github.com/mobilemindtec/go-io/option"
	"github.com/mobilemindtec/go-io/pipeline"
	"github.com/mobilemindtec/go-io/result"
	"github.com/mobilemindtec/go-io/state"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPipelineSimpleSum(t *testing.T) {
	res :=
		pipeline.New[int]().
			Next(func() int {
				return 5
			}).
			Next(func() int {
				return 2
			}).
			Next(func(x int, y int) int {
				return x + y
			}).
			UnsafeRun()

	assert.Equal(t, option.Some(7), res.Get())
}

func TestPipelineSimpleSumWithResult(t *testing.T) {
	res :=
		pipeline.New[int]().
			Next(func() int {
				return 5
			}).
			Next(func() int {
				return 2
			}).
			Next(func(x int, y int) *result.Result[int] {
				return result.OfValue(x + y)
			}).
			UnsafeRun()

	assert.Equal(t, option.Some(7), res.Get())
}

func TestPipelineSimpleFail(t *testing.T) {
	res :=
		pipeline.New[int]().
			Next(func() int {
				return 5
			}).
			Next(func() int {
				return 2
			}).
			Next(func(x int, y int) int {
				return x + y
			}).
			Next(func() *result.Result[int] {
				return result.OfError[int](errors.New("pipeline error"))
			}).
			UnsafeRun()

	assert.True(t, res.IsError())
	assert.Equal(t, "pipeline error", res.Error())
}

func TestPipelineSimpleSumWithState(t *testing.T) {
	res :=
		pipeline.New[int]().
			Next(func() int {
				return 5
			}).
			Next(func() int {
				return 2
			}).
			Next(func(st *state.State) int {
				x := state.Var[int](st)
				y := state.Var[int](st)
				return x + y
			}).
			UnsafeRun()

	assert.Equal(t, option.Some(7), res.Get())
}

func TestPipelineSimpleSumWithOpitonSome(t *testing.T) {
	res :=
		pipeline.New[int]().
			Next(func() int {
				return 5
			}).
			Next(func() int {
				return 2
			}).
			Next(func(x int, y int) *option.Option[int] {
				return option.Of(x + y)
			}).
			UnsafeRun()

	assert.Equal(t, option.Some(7), res.Get())
}

func TestPipelineSimpleSumWithOpitonNone(t *testing.T) {
	res :=
		pipeline.New[*int]().
			Next(func() *int {
				x := 5
				return &x
			}).
			Next(func() *int {
				x := 2
				return &x
			}).
			Next(func(x *int, y *int) *option.Option[*int] {
				return option.None[*int]()
			}).
			UnsafeRun()

	assert.False(t, res.IsError())
	assert.True(t, res.Get().Empty())
}
