package types

import (
	"fmt"
	"github.com/mobilemindtec/go-io/collections"
	"github.com/mobilemindtec/go-io/option"
	"github.com/mobilemindtec/go-io/result"
	"github.com/mobilemindtec/go-io/state"
	"github.com/mobilemindtec/go-io/util"
	"log"
	"reflect"
	"runtime"
)

type IOUnit = *IO[*Unit]

type IO[T any] struct {
	stack      *collections.Stack[IOEffect]
	varName    string
	state      *state.State
	debug      bool
	prevEffect IOEffect
	lastEffect IOEffect
	//suspendedIOs []IORunnable
}

func NewIO[T any]() *IO[T] {
	return &IO[T]{stack: collections.NewStack[IOEffect](), state: state.NewState()}
}

/*
func (this *IO[T]) GetSuspended() []IORunnable {
	return this.suspendedIOs
}

func (this *IO[T]) WithSuspended(sp []IORunnable) *IO[T] {
	this.suspendedIOs = sp
	return this
}*/

func (this *IO[T]) IOType() reflect.Type {
	return reflect.TypeFor[T]()
}

func (this *IO[T]) UnLift() IOEffect {
	return this.Effect()
}

func (this *IO[T]) Effect() IOEffect {
	if this.stack.Count() != 1 {
		panic("can't transform IO to Effect. IO has many or none Effects")
	}
	return this.stack.UnsafePeek()
}

func (this *IO[T]) push(val IOEffect) *IO[T] {
	this.stack.
		Peek().
		IfNonEmpty(func(eff IOEffect) {
			val.SetPrevEffect(eff)
		})
	this.stack.Push(val)
	return this
}

func (this *IO[T]) GetLastEffect() IOEffect {
	return this.lastEffect
}

func (this *IO[T]) SetPrevEffect(eff IOEffect) {
	this.prevEffect = eff
}

func (this *IO[T]) As(name string) *IO[T] {
	this.varName = name
	return this
}

func (this *IO[T]) Effects(vals ...IOEffect) *IO[T] {
	for _, eff := range vals {
		this.push(eff)
	}
	return this
}

func (this *IO[T]) Pipe(val IOEffect) *IO[T] {
	_, filename, line, _ := runtime.Caller(1)
	val.SetDebugInfo(&IODebugInfo{Line: line, Filename: filename})
	this.push(val)
	return this
}

func (this *IO[T]) LoadVar(val IOEffect) *IO[T] {
	_, filename, line, _ := runtime.Caller(1)
	val.SetDebugInfo(&IODebugInfo{Line: line, Filename: filename})
	this.push(val)
	return this
}

func (this *IO[T]) Pure(val IOEffect) *IO[T] {
	_, filename, line, _ := runtime.Caller(1)
	val.SetDebugInfo(&IODebugInfo{Line: line, Filename: filename})
	this.push(val)
	return this
}

func (this *IO[T]) Map(val IOEffect) *IO[T] {
	_, filename, line, _ := runtime.Caller(1)
	val.SetDebugInfo(&IODebugInfo{Line: line, Filename: filename})
	this.push(val)
	return this
}

func (this *IO[T]) FlatMap(val IOEffect) *IO[T] {
	_, filename, line, _ := runtime.Caller(1)
	val.SetDebugInfo(&IODebugInfo{Line: line, Filename: filename})
	this.push(val)
	return this
}

func (this *IO[T]) AndThan(val IOEffect) *IO[T] {
	_, filename, line, _ := runtime.Caller(1)
	val.SetDebugInfo(&IODebugInfo{Line: line, Filename: filename})
	this.push(val)
	return this
}

func (this *IO[T]) AndThanMany(vals ...IOEffect) *IO[T] {
	_, filename, line, _ := runtime.Caller(1)
	for _, val := range vals {
		val.SetDebugInfo(&IODebugInfo{Line: line, Filename: filename})
		this.push(val)
	}
	return this
}

func (this *IO[T]) Recover(val IOEffect) *IO[T] {
	_, filename, line, _ := runtime.Caller(1)
	val.SetDebugInfo(&IODebugInfo{Line: line, Filename: filename})
	this.push(val)
	return this
}

func (this *IO[T]) MaybeFail(val IOEffect) *IO[T] {
	_, filename, line, _ := runtime.Caller(1)
	val.SetDebugInfo(&IODebugInfo{Line: line, Filename: filename})
	this.push(val)
	return this
}

func (this *IO[T]) Filter(val IOEffect) *IO[T] {
	_, filename, line, _ := runtime.Caller(1)
	val.SetDebugInfo(&IODebugInfo{Line: line, Filename: filename})
	this.push(val)
	return this
}

func (this *IO[T]) Tap(val IOEffect) *IO[T] {
	_, filename, line, _ := runtime.Caller(1)
	val.SetDebugInfo(&IODebugInfo{Line: line, Filename: filename})
	this.push(val)
	return this
}

func (this *IO[T]) Or(val IOEffect) *IO[T] {
	_, filename, line, _ := runtime.Caller(1)
	val.SetDebugInfo(&IODebugInfo{Line: line, Filename: filename})
	this.push(val)
	return this
}

func (this *IO[T]) OrElse(val IOEffect) *IO[T] {
	_, filename, line, _ := runtime.Caller(1)
	val.SetDebugInfo(&IODebugInfo{Line: line, Filename: filename})
	this.push(val)
	return this
}

func (this *IO[T]) Debug(val IOEffect) *IO[T] {
	_, filename, line, _ := runtime.Caller(1)
	val.SetDebugInfo(&IODebugInfo{Line: line, Filename: filename})
	this.push(val)
	return this
}

func (this *IO[T]) Ensure(val IOEffect) *IO[T] {
	_, filename, line, _ := runtime.Caller(1)
	val.SetDebugInfo(&IODebugInfo{Line: line, Filename: filename})
	this.push(val)
	return this
}

func (this *IO[T]) SliceForeach(val IOEffect) *IO[T] {
	_, filename, line, _ := runtime.Caller(1)
	val.SetDebugInfo(&IODebugInfo{Line: line, Filename: filename})
	this.push(val)
	return this
}

func (this *IO[T]) SliceMap(val IOEffect) *IO[T] {
	_, filename, line, _ := runtime.Caller(1)
	val.SetDebugInfo(&IODebugInfo{Line: line, Filename: filename})
	this.push(val)
	return this
}

func (this *IO[T]) SliceFlatMap(val IOEffect) *IO[T] {
	_, filename, line, _ := runtime.Caller(1)
	val.SetDebugInfo(&IODebugInfo{Line: line, Filename: filename})
	this.push(val)
	return this
}

func (this *IO[T]) SliceFilter(val IOEffect) *IO[T] {
	_, filename, line, _ := runtime.Caller(1)
	val.SetDebugInfo(&IODebugInfo{Line: line, Filename: filename})
	this.push(val)
	return this
}

func (this *IO[T]) SliceAttemptOr(val IOEffect) *IO[T] {
	_, filename, line, _ := runtime.Caller(1)
	val.SetDebugInfo(&IODebugInfo{Line: line, Filename: filename})
	this.push(val)
	return this
}

func (this *IO[T]) SliceAttemptOrElse(val IOEffect) *IO[T] {
	_, filename, line, _ := runtime.Caller(1)
	val.SetDebugInfo(&IODebugInfo{Line: line, Filename: filename})
	this.push(val)
	return this
}

func (this *IO[T]) AsSlice(val IOEffect) *IO[T] {
	_, filename, line, _ := runtime.Caller(1)
	val.SetDebugInfo(&IODebugInfo{Line: line, Filename: filename})
	this.push(val)
	return this
}

func (this *IO[T]) Attempt(val IOEffect) *IO[T] {
	_, filename, line, _ := runtime.Caller(1)
	val.SetDebugInfo(&IODebugInfo{Line: line, Filename: filename})
	this.push(val)
	return this
}

func (this *IO[T]) Exec(val IOEffect) *IO[T] {
	_, filename, line, _ := runtime.Caller(1)
	val.SetDebugInfo(&IODebugInfo{Line: line, Filename: filename})
	this.push(val)
	return this
}

func (this *IO[T]) ExecIfEmpty(val IOEffect) *IO[T] {
	_, filename, line, _ := runtime.Caller(1)
	val.SetDebugInfo(&IODebugInfo{Line: line, Filename: filename})
	this.push(val)
	return this
}

func (this *IO[T]) FailIfEmpty(val IOEffect) *IO[T] {
	_, filename, line, _ := runtime.Caller(1)
	val.SetDebugInfo(&IODebugInfo{Line: line, Filename: filename})
	this.push(val)
	return this
}

func (this *IO[T]) FailIf(val IOEffect) *IO[T] {
	_, filename, line, _ := runtime.Caller(1)
	val.SetDebugInfo(&IODebugInfo{Line: line, Filename: filename})
	this.push(val)
	return this
}

func (this *IO[T]) Foreach(val IOEffect) *IO[T] {
	_, filename, line, _ := runtime.Caller(1)
	val.SetDebugInfo(&IODebugInfo{Line: line, Filename: filename})
	this.push(val)
	return this
}

func (this *IO[T]) CatchAll(val IOEffect) *IO[T] {
	_, filename, line, _ := runtime.Caller(1)
	val.SetDebugInfo(&IODebugInfo{Line: line, Filename: filename})
	this.push(val)
	return this
}

func (this *IO[T]) runStackIO(currEff IOEffect, sp int) IOEffect {
	if this.stack.IsNonEmpty() {
		eff := this.stack.UnsafePop()
		this.runStackIO(eff, sp-1)
	}

	if this.debug {
		log.Printf("IO>> UnsafeRun IO(Name=%v,SP=%v) %v",
			this.varName, sp, reflect.TypeOf(currEff))
	}

	if stf, ok := currEff.(IOStateful); ok {
		stf.SetState(this.state)
	}

	if this.debug {
		currEff.SetDebug(this.debug)
	}

	r := currEff.UnsafeRun()

	//if this.debug {
	//	log.Printf("IO>> UnsafeRun IO(%v) before %v = %v", this.varName,reflect.TypeOf(currEff).Name(), r.String())
	//}

	return r
}

func (this *IO[T]) UnsafeRun() *result.Result[*option.Option[T]] {

	if this.debug {
		log.Printf("IO>> run stack IO(%v) with %v operations, prevEffect = %v", this.varName, this.stack.Count(), this.prevEffect)
	}

	// last to execute
	lastEff := this.stack.UnsafePop()
	var firstEff IOEffect
	if this.stack.IsNonEmpty() {
		firstEff = this.stack.Last()
	} else {
		firstEff = lastEff
	}

	this.lastEffect = lastEff
	firstEff.SetPrevEffect(this.prevEffect)

	effResult := this.runStackIO(lastEff, this.stack.Count())
	r := effResult.GetResult()

	if r.IsError() {
		return result.OfError[*option.Option[T]](r.GetError())
	}

	if util.CanNil(reflect.ValueOf(r.GetValue()).Kind()) &&
		(util.IsNil(r.GetValue()) || r.Get().Empty()) {
		return result.OfValue(option.None[T]())
	}

	val := r.Get().GetValue()
	if v, ok := val.(T); ok {
		return result.OfValue(option.Some(v))
	}
	typOf := reflect.TypeFor[T]()
	panic(fmt.Sprintf("can't cast %v to IO(%v) result type %v", r.GetValue(), this.varName, typOf))
}

func (this *IO[T]) CheckTypesFlow() {

	var lastTypeOut reflect.Type
	var lastIO reflect.Type

	for _, it := range this.stack.GetItems() {

		if lastTypeOut == nil {
			lastTypeOut = it.TypeOut()
			lastIO = reflect.TypeOf(it).Elem()
		} else {

			if it.TypeIn() == reflect.TypeFor[*Unit]() {
				if it.TypeOut() != reflect.TypeFor[*Unit]() {
					lastTypeOut = it.TypeOut()
				}
				continue
			}

			if lastTypeOut != it.TypeIn() {

				curr := reflect.TypeOf(it).Elem()

				panic(fmt.Errorf("IO %v expect type is %v, but last IO %v result type is %v",
					curr.Name(), it.TypeIn(), lastIO.Name(), lastTypeOut))
			}

			lastTypeOut = it.TypeOut()
			lastIO = reflect.TypeOf(it).Elem()
		}
	}
}

func (this *IO[T]) SetState(st *state.State) {
	this.state = st
}

func (this *IO[T]) SetDebug(b bool) {
	this.debug = b
}

func (this *IO[T]) DebugOn() *IO[T] {
	this.SetDebug(true)
	return this
}

func (this *IO[T]) UnsafeRunIO() ResultOptionAny {
	return this.UnsafeRun().ToResultOfOption()
}

func (this *IO[T]) GetVarName() string {
	return this.varName
}
