package env

import "github.com/mobilemindtech/go-io/rio"

type EnvIO[R, A any] struct {
	thunk func(R) *rio.IO[A]
}

func (this *EnvIO[R, A]) Provide(env R) *rio.IO[A] {
	return this.thunk(env)
}

func Map[R, A, B any](env *EnvIO[R, A], f func(A) B) *EnvIO[R, B] {
	return &EnvIO[R, B]{
		func(r R) *rio.IO[B] {
			return rio.Map(env.thunk(r), f)
		},
	}
}
func FlatMap[R, A, B any](env *EnvIO[R, A], f func(A) *EnvIO[R, B]) *EnvIO[R, B] {
	return &EnvIO[R, B]{
		func(r R) *rio.IO[B] {
			return rio.FlatMap(
				env.thunk(r),
				func(a A) *rio.IO[B] {
					return f(a).thunk(r)
				})
		},
	}
}
