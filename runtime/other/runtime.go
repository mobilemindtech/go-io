package other

import (
	"fmt"
	eff "github.com/mobilemindtec/go-io/effect"
	"github.com/mobilemindtec/go-io/option"
	"github.com/mobilemindtec/go-io/result"
	"github.com/mobilemindtec/go-io/util"

	"reflect"
	"strings"
)

type RuntimeStatus int

const (
	AllDone RuntimeStatus = iota + 1
	Exit
	Error
)

type Func struct {
	Info *util.FuncInfo
	F    interface{}
}

type IRuntime interface {
	GetResults() []interface{}
	SetResults([]interface{})
}

type Runtime[T any] struct {
	stack      []interface{}
	resources  []interface{}
	results    []interface{}
	onError    func(error)
	onComplete func(T)
	onExit     func()
}

func New[T any]() *Runtime[T] {
	return &Runtime[T]{stack: []interface{}{}, resources: []interface{}{}}
}

func (this *Runtime[T]) GetResults() []interface{} {
	return this.results
}
func (this *Runtime[T]) SetResults(rs []interface{}) {
	this.results = rs
}

func (this *Runtime[T]) Effect(effects ...interface{}) *Runtime[T] {

	for _, ef := range effects {

		_, isEffecf := ef.(eff.IEffect)
		_, isResource := ef.(eff.IResource)
		_, isEffecfT := ef.(eff.IEffectT)

		if isResource {
			this.resources = append(this.resources, ef)
		} else if isEffecf || isEffecfT {
			this.stack = append(this.stack, ef)
		} else if util.IsFunc(ef) {
			f := &Func{
				Info: util.NewFuncInfo(ef),
				F:    ef,
			}

			if f.Info.ReturnCount != 1 {
				panic(fmt.Sprintf("the func should be return a *option.Result[T], but have %v",
					f.Info.ReturnCount))
			}

			rType := f.Info.ReturnTypes[0]
			if !strings.Contains(rType.String(), "*option.Result") {
				panic(fmt.Sprintf("the func should be return a *option.Result[T], but have %v",
					rType.Name()))
			}

			this.stack = append(this.stack, f)
		} else {
			panic("value should be a Resource Effect, EffectT, Pure or func")
		}
	}

	return this
}

func (this *Runtime[T]) Computation(cpt *Computation) *Runtime[T] {
	this.stack = append(this.stack, cpt)
	return this
}

func (this *Runtime[T]) Pure(f func() interface{}) *Runtime[T] {
	this.stack = append(this.stack, eff.NewValue(f))
	return this
}

func (this *Runtime[T]) EffectF(f func() interface{}) {
	this.stack = append(this.stack, eff.NewEff(f))
}

func (this *Runtime[T]) RunUnsafe() *result.Result[*option.Option[T]] {

	//var newStack []interface{}

	for _, it := range this.stack {
		if ef, ok := it.(eff.IEffect); ok {
			result := ef.RunEffect().(eff.IEffect).GetResult()
			fmt.Printf("result = %#v\n", result)
			this.results = append(this.results, result)
		} else {
			switch it.(type) {
			case *Computation:
				it.(*Computation).Run()
				break
			case *Test:
				r := it.(*Test).Run()
				if !r.IsTrue() {
					return result.OfValue(option.None[T]())
				}
				break
			}
		}
	}

	size := len(this.results)
	if size > 0 {
		last := this.results[size-1]
		if v, ok := resultToRaw[T](last); ok {
			return result.OfValue(option.Of(v))
		}
	}

	return result.OfValue(option.None[T]())
}

func resultToRaw[T any](value interface{}) (T, bool) {
	if r, ok := value.(*result.Result[T]); ok {
		return r.OrNil(), true
	} else if v, ok := value.(T); ok {
		return v, true
	}
	var t T
	return t, false
}

func GetValue[T any](rt IRuntime) T {

	results := rt.GetResults()
	size := len(results)
	var test T
	typ := reflect.TypeOf(test)
	for i := size; i >= 0; i++ {
		it := results[i]
		if v, ok := resultToRaw[T](it); ok {
			return v
		}
	}
	panic(fmt.Sprintf("value not found for %v", typ))
}

func ConsumeValue[T any](rt IRuntime) T {
	results := rt.GetResults()
	size := len(results)
	var test T
	typ := reflect.TypeOf(test)
	var newResults = make([]interface{}, size-1)
	var rawResult T
	var found bool
	for i := 0; i < size; i++ {
		it := results[i]

		if !found {
			if v, ok := resultToRaw[T](it); ok {
				rawResult = v
				found = true
				continue
			}
		}

		newResults = append(newResults, it)
	}

	if !found {
		panic(fmt.Sprintf("value not found for %v", typ))
	}

	rt.SetResults(newResults)

	return rawResult
}
