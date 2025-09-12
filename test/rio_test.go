package test

import (
	"errors"
	"fmt"
	"github.com/mobilemindtech/go-io/http"
	"github.com/mobilemindtech/go-io/result"
	"github.com/mobilemindtech/go-io/rio"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestNewIO(t *testing.T) {

	perforMap := rio.Map(
		rio.Pure(1), func(a int) string {
			return fmt.Sprintf("%v", a)
		})

	assert.Equal(t, "1", rio.UnsafeRun(perforMap).Get().Get())

	flatMap := rio.FlatMap(
		rio.Map(
			rio.Attempt(func() *result.Result[int] {
				return result.OfError[int](errors.New("error"))
			}),
			func(a int) string {
				return fmt.Sprintf("%v", a)
			}),
		func(s string) *rio.IO[int] {
			i, _ := strconv.Atoi(s)
			return rio.Pure(i)
		})

	assert.Equal(t, "error", rio.UnsafeRun(flatMap).GetError().Error())
}

func TestNewIOPipe(t *testing.T) {

	ioA := rio.Pure("Ricardo")
	ioB := rio.Pure("Bocchi")
	ioC := rio.Pure(37)

	pipe2IO := rio.FlatMap2(ioA, ioB, func(a string, b string) *rio.IO[string] {
		return rio.Pure(fmt.Sprintf("%v %v", a, b))
	})

	assert.Equal(t, "Ricardo Bocchi", rio.UnsafeRun(pipe2IO).Get().Get())

	pipeIO := rio.FlatMap3(ioA, ioB, ioC, func(a string, b string, c int) *rio.IO[string] {
		return rio.Pure(fmt.Sprintf("%v %v age %v", a, b, c))
	})

	assert.Equal(t, "Ricardo Bocchi age 37", rio.UnsafeRun(pipeIO).Get().Get())

}

func TestNewIOPipeMap(t *testing.T) {

	ioA := rio.Pure("Ricardo")
	ioB := rio.Pure("Bocchi")
	ioC := rio.Pure(37)

	pipe2IO := rio.Map2(ioA, ioB, func(a string, b string) string {
		return fmt.Sprintf("%v %v", a, b)
	})

	assert.Equal(t, "Ricardo Bocchi", rio.UnsafeRun(pipe2IO).Get().Get())

	pipeIO := rio.Map3(ioA, ioB, ioC, func(a string, b string, c int) string {
		return fmt.Sprintf("%v %v age %v", a, b, c)
	})

	assert.Equal(t, "Ricardo Bocchi age 37", rio.UnsafeRun(pipeIO).Get().Get())

}

func TestHttpRIO(t *testing.T) {

	reqIO := http.
		NewClient[any, *User, any]().
		AsJSON().
		GetRIO("http://4gym.com.br/tools-api/mock/user")

	filterIO := rio.Filter[*http.Response[*User, any]](
		reqIO,
		func(r *http.Response[*User, any]) bool {
			return r.StatusCode == 200
		})

	mapIO := rio.Map[*http.Response[*User, any], int](
		filterIO,
		func(r *http.Response[*User, any]) int {
			return r.StatusCode
		})

	assert.Equal(t, 200, rio.UnsafeRun(mapIO).Get().Get())

}
