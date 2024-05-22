package other

import (
	eff "github.com/mobilemindtec/go-io/effect"
	"github.com/mobilemindtec/go-io/option"
	"github.com/mobilemindtec/go-io/result"
	"testing"
)

func TestRuntime(t *testing.T) {
	rt := New[int]()
	result := rt.
		Effect(
			eff.NewEffect(func() *result.Result[int] { return result.OfValue(1) }),
			eff.NewEffect(func() *result.Result[int] {
				return result.OfValue(2)
			}),
			eff.NewPure(func() int {
				v1 := ConsumeValue[int](rt)
				v2 := ConsumeValue[int](rt)
				return v1 + v2
			}),
		).
		RunUnsafe()

	result.
		IfError(func(err error) {
			t.Error(err)
		}).
		IfOk(func(opt *option.Option[int]) {
			opt.
				IfEmpty(func() {
					t.Error("empty result")
				}).
				IfNonEmpty(func(i int) {
					if i != 3 {
						t.Errorf("expected 3, but found %v", i)
					}
				})
		})

}
