package runtime

import (
	"github.com/mobilemindtec/go-io/result"
	"github.com/mobilemindtec/go-io/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

type PersonPtr = *types.Person

func TestRuntime(t *testing.T) {

	rt := NewRuntime[PersonPtr]()

	rt.Effects(
		types.New[string]().
			Pure(types.NewPure[PersonPtr](&types.Person{Name: "Ricardo"})).
			Map(types.NewMap[PersonPtr, string](
				func(ptr PersonPtr) string {
					return ptr.Name
				})).
			As("name"),
		types.New[PersonPtr]().
			Attempt(types.NewAttempt(
				func() *result.Result[PersonPtr] {
					name, _ := Var[string](rt, "name")
					return result.OfValue(&types.Person{Name: name})
				})).
			As("person"),
	).UnsafeRun()

	assert.Equal(t, "Ricardo", rt.Yield().Name)
}
