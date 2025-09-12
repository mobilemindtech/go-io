package ios

import (
	"fmt"
	"github.com/mobilemindtech/go-io/option"
	"github.com/mobilemindtech/go-io/result"
	"github.com/mobilemindtech/go-io/state"
	"github.com/mobilemindtech/go-io/types"
	"github.com/mobilemindtech/go-io/util"
	"log"
	"reflect"
	"runtime/debug"
)

type IOSliceAttemptEach[A any] struct {
	value      *result.Result[*option.Option[[]A]]
	prevEffect types.IOEffect
	feach      func(A) *result.Result[A]
	feachState func(A, *state.State) *result.Result[A]
	debug      bool
	debugInfo  *types.IODebugInfo
	state      *state.State
}

func NewSliceAttemptEach[A any](f func(A) *result.Result[A]) *IOSliceAttemptEach[A] {
	return &IOSliceAttemptEach[A]{feach: f}
}

func NewSliceAttemptEachWithState[A any](f func(A, *state.State) *result.Result[A]) *IOSliceAttemptEach[A] {
	return &IOSliceAttemptEach[A]{feachState: f}
}

func (this *IOSliceAttemptEach[A]) Lift() *types.IO[A] {
	return types.NewIO[A]().Effects(this)
}

func (this *IOSliceAttemptEach[A]) TypeIn() reflect.Type {
	return reflect.TypeFor[[]A]()
}

func (this *IOSliceAttemptEach[A]) TypeOut() reflect.Type {
	return reflect.TypeFor[[]A]()
}

func (this *IOSliceAttemptEach[A]) SetDebug(b bool) {
	this.debug = b
}

func (this *IOSliceAttemptEach[A]) SetState(st *state.State) {
	this.state = st
}

func (this *IOSliceAttemptEach[A]) SetDebugInfo(info *types.IODebugInfo) {
	this.debugInfo = info
}

func (this *IOSliceAttemptEach[A]) GetDebugInfo() *types.IODebugInfo {
	return this.debugInfo
}

func (this *IOSliceAttemptEach[A]) String() string {
	return fmt.Sprintf("SliceAttemptEach(%v)", this.value.String())
}

func (this *IOSliceAttemptEach[A]) SetPrevEffect(prev types.IOEffect) {
	this.prevEffect = prev
}

func (this *IOSliceAttemptEach[A]) GetPrevEffect() *option.Option[types.IOEffect] {
	return option.Of(this.prevEffect)
}

func (this *IOSliceAttemptEach[A]) GetResult() types.ResultOptionAny {
	return this.value.ToResultOfOption()
}

func (this *IOSliceAttemptEach[A]) UnsafeRun() types.IOEffect {
	var currEff interface{} = this
	prevEff := this.GetPrevEffect()
	this.value = result.OfValue(option.None[[]A]())

	if prevEff.NonEmpty() {
		r := prevEff.Get().GetResult()
		if r.IsError() {
			this.value = result.OfError[*option.Option[[]A]](r.Failure())
		} else if !r.Get().IsEmpty() {

			defer func() {
				if r := recover(); r != nil {

					err := types.NewIOError(fmt.Sprintf("%v", r), debug.Stack())

					if this.debug {
						if this.debugInfo != nil {
							log.Printf("[DEBUG IOSliceAttemptEach]=>> added in: %v:%v", this.debugInfo.Filename, this.debugInfo.Line)
						}
						log.Printf("[DEBUG IOSliceAttemptEach]=>> Error: %v\n", err.Error())
						log.Printf("[DEBUG IOSliceAttemptEach]=>> StackTrace: %v\n", err.StackTrace)
					}

					this.value = result.OfError[*option.Option[[]A]](err)
				}
			}()

			var results []A

			val := r.Get().GetValue()
			if effValue, ok := val.([]A); ok {

				for _, it := range effValue {

					var res *result.Result[A]

					if this.feach != nil {
						res = this.feach(it)
					} else {
						res = this.feachState(it, this.state)
					}

					if res.IsError() {
						this.value = result.OfError[*option.Option[[]A]](res.Failure())
						break
					} else {
						results = append(results, res.Get())
					}
				}

				if len(results) > 0 {
					this.value = result.OfValue(option.Some(results))
				}

			} else {
				util.PanicCastType("IOSliceAttemptEach",
					reflect.TypeOf(val), reflect.TypeFor[[]A]())
			}

		}
	}

	if this.debug {
		log.Printf("%v\n", this.String())
	}

	return currEff.(types.IOEffect)
}
