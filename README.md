# IO Monad Go Lang

Attention. This is an experimental project that I use in my company. The fact that Go Lang does not support generics in struct functions makes the whole thing more insecure. So I started with building a Pipeline and then went to the IO transformations that seem more interesting to me. I ended up implementing two versions, one with functions and one with struct. Below is a minimum documentation of the use of Pipeline and the two IO implementations.

### Pipeline

Pipeline should return:
- any
- (any, error)
- (string, any) -> var name, type
- *result.Result[any]
- *option.Option[any]
- *result.Result[*option.Option[any]]
  Any return of error or None stop pipeline.
  The last computation should be return a same type of Pipeline[T] generic type.

```go

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

	assert.Equal(t, 7, res.Get().Get())
}

```


### RIO

Experimental IO operations using functions

```go
rio.Pure[T any](value T) *IO[T]
rio.PureF[T any](f func() T) *IO[T]
rio.Map[A, B any](io *IO[A], f func(A) B) *IO[B]
rio.FlatMap[A, B any](io *IO[A], f func(A) *IO[B]) *IO[B]
rio.AndThan[A, B any](io *IO[A], f func() *IO[B]) *IO[B]
rio.Filter[A any](io *IO[A], f func(A) bool) *IO[A]
rio.Foreach[A any](io *IO[A], f func(A)) *IO[A]
rio.OrElse[A any](io *IO[A], f func() *IO[A]) *IO[A]
rio.Or[A any](io *IO[A], f func() A) *IO[A]
rio.Recover[A any](io *IO[A], f func(error) A) *IO[A]
rio.RecoverIO[A any](io *IO[A], f func(error) *IO[A]) *IO[A]
rio.CatchAll[A any](io *IO[A], f func(error)) *IO[A]
rio.Ensure[A any](io *IO[A], f func()) *IO[A]
rio.Debug[A any](io *IO[A], label ...string) *IO[A]
rio.Attempt[A any](f func() *result.Result[A]) *IO[A]
rio.Pipe2[A, B, T any](a *IO[A], b *IO[B], f func(A, B) *IO[T]) *IO[T]
rio.Pipe3[A, B, C, T any](a *IO[A], b *IO[B], c *IO[C], f func(A, B, C) *IO[T]) *IO[T]
rio.Pipe4[A, B, C, D, T any](a *IO[A], b *IO[B], c *IO[C], d *IO[D], f func(A, B, C, D) *IO[T]) *IO[T]
rio.Pipe5[A, B, C, D, E, T any](a *IO[A], b *IO[B], c *IO[C], d *IO[D], e *IO[E], f func(A, B, C, D, E) *IO[T]) *IO[T]
rio.Pipe6[A, B, C, D, E, F, T any](a *IO[A], b *IO[B], c *IO[C], d *IO[D], e *IO[E], f *IO[F], fn func(A, B, C, D, E, F) *IO[T]) *IO[T]
rio.Pipe7[A, B, C, D, E, F, G, T any](a *IO[A], b *IO[B], c *IO[C], d *IO[D], e *IO[E], f *IO[F], g *IO[G], fn func(A, B, C, D, E, F, G) *IO[T]) *IO[T]
rio.Pipe8[A, B, C, D, E, F, G, H, T any](a *IO[A], b *IO[B], c *IO[C], d *IO[D], e *IO[E], f *IO[F], g *IO[G], h *IO[H], fn func(A, B, C, D, E, F, G, H) *IO[T]) *IO[T]
rio.Pipe9[A, B, C, D, E, F, G, H, I, T any](a *IO[A], b *IO[B], c *IO[C], d *IO[D], e *IO[E], f *IO[F], g *IO[G], h *IO[H], i *IO[I], fn func(A, B, C, D, E, F, G, H, I) *IO[T]) *IO[T]
rio.Pipe10[A, B, C, D, E, F, G, H, I, J, T any](a *IO[A], b *IO[B], c *IO[C], d *IO[D], e *IO[E], f *IO[F], g *IO[G], h *IO[H], i *IO[I], j *IO[J], fn func(A, B, C, D, E, F, G, H, I, J) *IO[T]) *IO[T]
rio.UnsafeRun[T any](io *IO[T]) (r *result.Result[*option.Option[T]])
```
```golang

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

```
### IO

Experimental IO operations using "magic" app




The IO[T] agregate a set of effects. The result of the last effect should the same of IO generic type.

For exemple:


```go
io.IOApp[string](

	io.IO[string](	// IO operation
		io.Pure(func() int { // effect 1 
            return 1
        }).
		Map(func(i int) string { // effect 1
            return fmt.Sprintf("%v", i)
        }),
	),

).UnsafeRun()	
``` 

See as the last effect return string, that is the same type of IO operation.

```go
io.IOApp[string](
	
	io.IO[string](	// IO operation
		io.Pure(func() int { // effect 1
				return 1
		}).
		FailIf(func(i int) error { // effect 2
            if i == 0 {
                return fmt.Errorf("value can't be zero")
            }
            return nil
        }).
		Map(func(i int) string { // effect 3
            return fmt.Sprintf("%v", i)
        }),
	),

).UnsafeRun()
``` 

It is also possible apply transformations between distinct IO's

```go

io.IOApp[string](

	io.IO[int](	// IO operation

		io.Pure(func() int { // effect 1 
            return 1
        }).
		FailIf(func(i int) error { // this effect can be arbitrary type
            if i == 0 {
                return fmt.Errorf("value can't be zero")
            }
            return nil
        }),
	),

	io.IO[string]( // IO receive result of last IO
		io.Map(func(i int) string {
            return fmt.Sprintf("%v", i)
        }),
	),

).UnsafeRun()
``` 

Http example:

```go
	app := io.IOApp[int](
		http.
			NewClient[any, *User, any]().
			AsJSON().
			GetIO("http://4gym.com.br/tools-api/mock/user"). // IO[*http.Response[*User, any]]
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
```


The IOApp have an internal cache to store IO results. Using Attempt* effects, we can get a stored value by type. Or else
we can get stored type using any Pipe* function. For exemple

```go

io.IOApp[int] (

	io.PureVal(1).Lift(),
	io.PureVal(2).Lift(),
	io.Pipe2(func(x int, y int) int {
		// x = 2, y = 1
		return x + y
	}),

).UnsafeRun()

``` 

or

```go

io.IOApp[int] (

	io.PureVal(1).Lift(),
	io.PureVal(2).Lift(),
	io.AttemptValueState(func(s *state.State) int {
		x := state.Consume[int](s) // 2
		y := state.Consume[int](s) // 1
		return x + y
	}),

).UnsafeRun()

```

`state.Consume` ensures, in this case, that var is used only one time into effect. If that's not a concern, as in most cases, you can use `state.Var[T](s)`.


The stored values are consumed from last to first in both cases.

Finally, you can return a suspended computation

```go

func SuspendedComputation() *types.IOSuspended {
	return types.Suspend(
		io.PureVal(1).Lift(),
		io.PureVal(2).Lift(),
	)
}

func MyApp() {
	
	io.IOApp[int] (

		SuspendedComputation(),

		io.AttemptValueState(func(s *state.State) int {
			x := state.Consume[int](s) // 2
			y := state.Consume[int](s) // 1
			return x + y
		}),

	).UnsafeRun()

}

```

All Effects:

```go
io.AndThan[A any](f func() types.IORunnable) *ios.IOAndThan[A]	
io.Then[A any](f func(A) A) *ios.IOThen[A]
io.Ensure[A any](f func()) *ios.IOEnsure[A]
io.Debug[A any](label string) *ios.IODebug[A]
io.MaybeFail[A any](f func(A) *result.Result[A]) *ios.IOMaybeFail[A]
io.MaybeFailError[A any](f func(A) error) *ios.IOMaybeFail[A]
io.Filter[A any](f func(A) bool) *ios.IOFilter[A]
io.FlatMap[A any, B any](f func(A) types.IORunnable) *ios.IOFlatMap[A, B]
io.Map[A any, B any](f func(A) B) *ios.IOMap[A, B]
io.PureVal[T any](value T) *ios.IOPure[T]
io.Pure[T any](f func() T) *ios.IOPure[T]
io.RecoverPure[A any](f func(error) A) *ios.IORecover[A]
io.Recover[A any](f func(error) *result.Result[A]) *ios.IORecover[A]
io.RecoverOption[A any](f func(error) *option.Option[A]) *ios.IORecover[A]
io.RecoverResultOption[A any](f func(error) *result.Result[*option.Option[A]]) *ios.IORecover[A]
io.SliceFilter[A any](f func(A) bool) *ios.IOSliceFilter[A]
io.SliceFlatMap[A any, B any](f func(A) types.IORunnable) *ios.IOSliceFlatMap[A, B]
io.SliceForeach[A any](f func(A)) *ios.IOSliceForeach[A]
io.SliceAttemptOrElse[A any](f func() *result.Result[[]A]) *ios.IOSliceAttemptOrElse[A]
io.SliceAttemptOrElseWithState[A any](f func(*state.State) *result.Result[[]A]) *ios.IOSliceAttemptOrElse[A]
io.SliceMap[A any, B any](f func(A) B) *ios.IOSliceMap[A, B]
io.SliceOr[A any](f func() []A) *ios.IOSliceOr[A]
io.SliceOrElse[A any](f func() types.IORunnable) *ios.IOSliceOrElse[A]
io.AsSliceOf[A any]() *ios.IOAsSlice[A]
io.Tap[A any](f func(A) bool) *ios.IOTap[A]
io.Foreach[A any](f func(A)) *ios.IOForeach[A]
io.Or[A any](f func() A) *ios.IOOr[A]
io.FailIfEmpty[A any](f func() error) *ios.IOFailIfEmpty[A]
io.FailWith[A any](f func() error) *ios.IOFailWith[A]
io.FailIf[A any](f func(A) error) *ios.IOFailIf[A]
io.OrElse[A any](f func() types.IORunnable) *ios.IOOrElse[A]
io.CatchAll[A any](f func(error)) *ios.IOCatchAll[A]
io.Nohup[A any]() *ios.IONohup[A]
io.Attempt[A any](f func() *result.Result[A]) *ios.IOAttempt[A]
io.AttemptOfOption[A any](f func() *option.Option[A]) *ios.IOAttempt[A]
io.AttemptOfResultOption[A any](f func() *result.Result[*option.Option[A]]) *ios.IOAttempt[A]
io.AttemptState[A any](f func(*state.State) *result.Result[A]) *ios.IOAttempt[A]
io.AttemptStateOfOption[A any](f func(*state.State) *option.Option[A]) *ios.IOAttempt[A]
io.AttemptStateOfResultOption[A any](f func(*state.State) *result.Result[*option.Option[A]]) *ios.IOAttempt[A]
io.AttemptOfUnit[A any](f func()) *ios.IOAttempt[A]
io.AttemptStateOfUnit(f func(*state.State)) *ios.IOAttempt[*types.Unit]
io.AttemptOfError[A any](f func() (A, error)) *ios.IOAttempt[A]
io.AttemptStateOfError[A any](f func(*state.State) (A, error)) *ios.IOAttempt[A]
io.AttemptPureState[A any](f func(*state.State) A) *ios.IOAttempt[A]
io.AttemptAndThanWithState[A any](f func(*state.State) types.IORunnable) *ios.IOAttemptAndThan[A]
io.AttemptAndThan[A any](f func() types.IORunnable) *ios.IOAttemptAndThan[A]
io.AttemptRunIOWithState[A any](f func(*state.State) types.IORunnable) *ios.IOAttemptAndThan[A]
io.AttemptRunIO[A any](f func() types.IORunnable) *ios.IOAttemptAndThan[A]
io.AttemptAuto[A any](f interface{}) *ios.IOAttemptAuto[A]
io.AttemptExec[A any](f func(A)) *ios.IOAttemptExec[A]
io.AttemptExecWithState[A any](f func(A, *state.State)) *ios.IOAttemptExec[A]
io.AttemptExecOrElse[A any](f func()) *ios.IOAttemptExecOrElse[A]
io.AttemptExecOrElseWithState[A any](f func(*state.State)) *ios.IOAttemptExecOrElse[A]
io.AttemptOrElseWithState[A any](f func(*state.State) types.IORunnable) *ios.IOAttemptOrElse[A]io.AttemptOrElse[A any](f func() types.IORunnable) *ios.IOAttemptOrElse[A]
io.AttemptThen[A any](f func(A) *result.Result[A]) *ios.IOAttemptThen[A]
io.AttemptThenWithState[A any](f func(A, *state.State) *result.Result[A]) *ios.IOAttemptThen[A]
io.AttemptThenOption[A any](f func(A) *result.Result[*option.Option[A]]) *ios.IOAttemptThen[A]
io.AttemptThenOptionWithState[A any](f func(A, *state.State) *result.Result[*option.Option[A]]) *ios.IOAttemptThen[A]
io.AttemptThenIO[A any](f func(A) types.IORunnable) *ios.IOAttemptThen[A]
io.AttemptThenIOWithState[A any](f func(A, *state.State) types.IORunnable) *ios.IOAttemptThen[A]
io.PipeIO[A, T any](f func(A) types.IORunnable) *ios.IOPipe[A, T]
io.Pipe[A, T any](f func(A) *result.Result[*option.Option[T]]) *ios.IOPipe[A, T]
io.PipeOfValue[A, T any](f func(A) T) *ios.IOPipe[A, T]
io.PipeOfResult[A, T any](f func(A) *result.Result[T]) *ios.IOPipe[A, T]
io.PipeOfOption[A, T any](f func(A) *option.Option[T]) *ios.IOPipe[A, T]
io.Pipe2IO[A, B, T any](f func(A, B) types.IORunnable) *ios.IOPipe2[A, B, T]
io.Pipe2[A, B, T any](f func(A, B) *result.Result[*option.Option[T]]) *ios.IOPipe2[A, B, T]
io.Pipe2OfValue[A, B, T any](f func(A, B) T) *ios.IOPipe2[A, B, T]
io.Pipe2OfResult[A, B, T any](f func(A, B) *result.Result[T]) *ios.IOPipe2[A, B, T]
io.Pipe2OfOption[A, B, T any](f func(A, B) *option.Option[T]) *ios.IOPipe2[A, B, T]
io.Pipe3IO[A, B, C, T any](f func(A, B, C) types.IORunnable) *ios.IOPipe3[A, B, C, T]
io.Pipe3[A, B, C, T any](f func(A, B, C) *result.Result[*option.Option[T]]) *ios.IOPipe3[A, B, C, T]
io.Pipe3OfValue[A, B, C, T any](f func(A, B, C) T) *ios.IOPipe3[A, B, C, T]
io.Pipe3OfResult[A, B, C, T any](f func(A, B, C) *result.Result[T]) *ios.IOPipe3[A, B, C, T]
io.Pipe3OfOption[A, B, C, T any](f func(A, B, C) *option.Option[T]) *ios.IOPipe3[A, B, C, T]
io.Pipe4IO[A, B, C, D, T any](f func(A, B, C, D) types.IORunnable) *ios.IOPipe4[A, B, C, D, T]
io.Pipe4[A, B, C, D, T any](f func(A, B, C, D) *result.Result[*option.Option[T]]) *ios.IOPipe4[A, B, C, D, T]
io.Pipe4OfValue[A, B, C, D, T any](f func(A, B, C, D) T) *ios.IOPipe4[A, B, C, D, T]
io.Pipe4OfResult[A, B, C, D, T any](f func(A, B, C, D) *result.Result[T]) *ios.IOPipe4[A, B, C, D, T]
io.Pipe4OfOption[A, B, C, D, T any](f func(A, B, C, D) *option.Option[T]) *ios.IOPipe4[A, B, C, D, T]
io.Pipe5IO[A, B, C, D, E, T any](f func(A, B, C, D, E) types.IORunnable) *ios.IOPipe5[A, B, C, D, E, T]
io.Pipe5[A, B, C, D, E, T any](f func(A, B, C, D, E) *result.Result[*option.Option[T]]) *ios.IOPipe5[A, B, C, D, E, T]
io.Pipe5OfValue[A, B, C, D, E, T any](f func(A, B, C, D, E) T) *ios.IOPipe5[A, B, C, D, E, T]
io.Pipe5OfResult[A, B, C, D, E, T any](f func(A, B, C, D, E) *result.Result[T]) *ios.IOPipe5[A, B, C, D, E, T]
io.Pipe5OfOption[A, B, C, D, E, T any](f func(A, B, C, D, E) *option.Option[T]) *ios.IOPipe5[A, B, C, D, E, T]
io.Pipe6IO[A, B, C, D, E, F, T any](f func(A, B, C, D, E, F) types.IORunnable) *ios.IOPipe6[A, B, C, D, E, F, T]
io.Pipe6[A, B, C, D, E, F, T any](f func(A, B, C, D, E, F) *result.Result[*option.Option[T]]) *ios.IOPipe6[A, B, C, D, E, F, T]
io.Pipe6OfValue[A, B, C, D, E, F, T any](f func(A, B, C, D, E, F) T) *ios.IOPipe6[A, B, C, D, E, F, T]
io.Pipe6OfResult[A, B, C, D, E, F, T any](f func(A, B, C, D, E, F) *result.Result[T]) *ios.IOPipe6[A, B, C, D, E, F, T]
io.Pipe6OfOption[A, B, C, D, E, F, T any](f func(A, B, C, D, E, F) *option.Option[T]) *ios.IOPipe6[A, B, C, D, E, F, T]
io.Pipe7IO[A, B, C, D, E, F, G, T any](f func(A, B, C, D, E, F, G) types.IORunnable) *ios.IOPipe7[A, B, C, D, E, F, G, T]
io.Pipe7[A, B, C, D, E, F, G, T any](f func(A, B, C, D, E, F, G) *result.Result[*option.Option[T]]) *ios.IOPipe7[A, B, C, D, E, F, G, T]
io.Pipe7OfValue[A, B, C, D, E, F, G, T any](f func(A, B, C, D, E, F, G) T) *ios.IOPipe7[A, B, C, D, E, F, G, T]
io.Pipe7OfResult[A, B, C, D, E, F, G, T any](f func(A, B, C, D, E, F, G) *result.Result[T]) *ios.IOPipe7[A, B, C, D, E, F, G, T]
io.Pipe7OfOption[A, B, C, D, E, F, G, T any](f func(A, B, C, D, E, F, G) *option.Option[T]) *ios.IOPipe7[A, B, C, D, E, F, G, T]
io.Pipe8IO[A, B, C, D, E, F, G, H, T any](f func(A, B, C, D, E, F, G, H) types.IORunnable) *ios.IOPipe8[A, B, C, D, E, F, G, H, T]
io.Pipe8[A, B, C, D, E, F, G, H, T any](f func(A, B, C, D, E, F, G, H) *result.Result[*option.Option[T]]) *ios.IOPipe8[A, B, C, D, E, F, G, H, T]
io.Pipe8OfValue[A, B, C, D, E, F, G, H, T any](f func(A, B, C, D, E, F, G, H) T) *ios.IOPipe8[A, B, C, D, E, F, G, H, T]
io.Pipe8OfResult[A, B, C, D, E, F, G, H, T any](f func(A, B, C, D, E, F, G, H) *result.Result[T]) *ios.IOPipe8[A, B, C, D, E, F, G, H, T]
io.Pipe8OfOption[A, B, C, D, E, F, G, H, T any](f func(A, B, C, D, E, F, G, H) *option.Option[T]) *ios.IOPipe8[A, B, C, D, E, F, G, H, T]
io.Pipe9IO[A, B, C, D, E, F, G, H, I, T any](f func(A, B, C, D, E, F, G, H, I) types.IORunnable) *ios.IOPipe9[A, B, C, D, E, F, G, H, I, T]
io.Pipe9[A, B, C, D, E, F, G, H, I, T any](f func(A, B, C, D, E, F, G, H, I) *result.Result[*option.Option[T]]) *ios.IOPipe9[A, B, C, D, E, F, G, H, I, T]
io.Pipe9OfValue[A, B, C, D, E, F, G, H, I, T any](f func(A, B, C, D, E, F, G, H, I) T) *ios.IOPipe9[A, B, C, D, E, F, G, H, I, T]
io.Pipe9OfResult[A, B, C, D, E, F, G, H, I, T any](f func(A, B, C, D, E, F, G, H, I) *result.Result[T]) *ios.IOPipe9[A, B, C, D, E, F, G, H, I, T]
io.Pipe9OfOption[A, B, C, D, E, F, G, H, I, T any](f func(A, B, C, D, E, F, G, H, I) *option.Option[T]) *ios.IOPipe9[A, B, C, D, E, F, G, H, I, T]
io.Pipe10IO[A, B, C, D, E, F, G, H, I, J, T any](f func(A, B, C, D, E, F, G, H, I, J) types.IORunnable) *ios.IOPipe10[A, B, C, D, E, F, G, H, I, J, T]
io.Pipe10[A, B, C, D, E, F, G, H, I, J, T any](f func(A, B, C, D, E, F, G, H, I, J) *result.Result[*option.Option[T]]) *ios.IOPipe10[A, B, C, D, E, F, G, H, I, J, T]
io.Pipe10OfValue[A, B, C, D, E, F, G, H, I, J, T any](f func(A, B, C, D, E, F, G, H, I, J) T) *ios.IOPipe10[A, B, C, D, E, F, G, H, I, J, T]
io.Pipe10OfResult[A, B, C, D, E, F, G, H, I, J, T any](f func(A, B, C, D, E, F, G, H, I, J) *result.Result[T]) *ios.IOPipe10[A, B, C, D, E, F, G, H, I, J, T]
io.Pipe10OfOption[A, B, C, D, E, F, G, H, I, J, T any](f func(A, B, C, D, E, F, G, H, I, J) *option.Option[T]) *ios.IOPipe10[A, B, C, D, E, F, G, H, I, J, T]
```
