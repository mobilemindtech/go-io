package effect

import (
	"github.com/mobilemindtech/go-io/result"
	"github.com/mobilemindtech/go-io/util"
)

type IEffect interface {
	IsEffect() bool
	IsPure() bool
	GetResult() interface{}
	RunEffect() interface{}
}

type IResource interface {
	IsResource() bool
}

type IEffectT interface {
	ArgsCount() int
}

type Eff struct {
	f     func() interface{}
	value interface{}
}

func NewEff(f func() interface{}) *Eff {
	return &Eff{f: f}
}

func (this Eff) IsEffect() bool {
	return true
}

func (this Eff) IsPure() bool {
	return false
}

func (this Eff) GetResult() interface{} {
	if util.IsNil(this.value) {
		panic("Effect result not available")
	}
	return this.GetResult()
}

func (this Eff) Do() interface{} {
	value := this.f()
	if _, ok := value.(*result.IResult); !ok {
		panic("effect value should be a Result")
	}
	return &Eff{f: this.f, value: value}
}

type Effect[T any] struct {
	f      func() *result.Result[T]
	result *result.Result[T]
}

func NewEffect[T any](f func() *result.Result[T]) *Effect[T] {
	return &Effect[T]{f: f}
}

func (this *Effect[T]) Run() *Effect[T] {
	r := this.f()
	if r == nil {
		panic("effect can't return nil")
	}
	return &Effect[T]{result: r, f: this.f}
}

func (this *Effect[T]) Result() *result.Result[T] {
	if this.result == nil {
		panic("effect was not executed")
	}
	return this.result
}

func (this *Effect[T]) IsEffect() bool {
	return true
}

func (this *Effect[T]) IsPure() bool {
	return false
}

func (this *Effect[T]) GetResult() interface{} {
	return this.Result()
}

func (this *Effect[T]) RunEffect() interface{} {
	return this.Run()
}

// Effect with 1 arg
type EffectT1[T any, T1 any] struct {
	f      func(T1) *result.Result[T]
	result *result.Result[T]
}

func NewEffectT1[T any, T1 any](f func(T1) *result.Result[T]) *EffectT1[T, T1] {
	return &EffectT1[T, T1]{f: f}
}

func (this *EffectT1[T, T1]) ArgsCount() int {
	return 1
}

func (this *EffectT1[T, T1]) Run(v T1) *EffectT1[T, T1] {
	r := this.f(v)
	if r == nil {
		panic("effect can't return nil")
	}
	return &EffectT1[T, T1]{result: r, f: this.f}
}

func (this *EffectT1[T, T1]) Result() *result.Result[T] {
	if this.result == nil {
		panic("effect was not executed")
	}
	return this.result
}

// Effect with 2 args
type EffectT2[T any, T1 any, T2 any] struct {
	f      func(T1, T2) *result.Result[T]
	result *result.Result[T]
}

func NewEffectT2[T any, T1 any, T2 any](f func(T1, T2) *result.Result[T]) *EffectT2[T, T1, T2] {
	return &EffectT2[T, T1, T2]{f: f}
}

func (this *EffectT2[T, T1, T2]) ArgsCount() int {
	return 2
}

func (this *EffectT2[T, T1, T2]) Run(v1 T1, v2 T2) *EffectT2[T, T1, T2] {
	r := this.f(v1, v2)
	if r == nil {
		panic("effect can't return nil")
	}
	return &EffectT2[T, T1, T2]{result: r, f: this.f}
}

func (this *EffectT2[T, T1, T2]) Result() *result.Result[T] {
	if this.result == nil {
		panic("effect was not executed")
	}
	return this.result
}

// Effect with 3 args
type EffectT3[T any, T1 any, T2 any, T3 any] struct {
	f      func(T1, T2, T3) *result.Result[T]
	result *result.Result[T]
}

func NewEffectT3[T any, T1 any, T2 any, T3 any](f func(T1, T2, T3) *result.Result[T]) *EffectT3[T, T1, T2, T3] {
	return &EffectT3[T, T1, T2, T3]{f: f}
}

func (this *EffectT3[T, T1, T2, T3]) ArgsCount() int {
	return 3
}

func (this *EffectT3[T, T1, T2, T3]) Run(v1 T1, v2 T2, v3 T3) *EffectT3[T, T1, T2, T3] {
	r := this.f(v1, v2, v3)
	if r == nil {
		panic("effect can't return nil")
	}
	return &EffectT3[T, T1, T2, T3]{result: r, f: this.f}
}

func (this *EffectT3[T, T1, T2, T3]) Result() *result.Result[T] {
	if this.result == nil {
		panic("effect was not executed")
	}
	return this.result
}

// Effect with 4 args
type EffectT4[T any, T1 any, T2 any, T3 any, T4 any] struct {
	f      func(T1, T2, T3, T4) *result.Result[T]
	result *result.Result[T]
}

func NewEffectT4[T any, T1 any, T2 any, T3 any, T4 any](f func(T1, T2, T3, T4) *result.Result[T]) *EffectT4[T, T1, T2, T3, T4] {
	return &EffectT4[T, T1, T2, T3, T4]{f: f}
}

func (this *EffectT4[T, T1, T2, T3, T4]) ArgsCount() int {
	return 4
}

func (this *EffectT4[T, T1, T2, T3, T4]) Run(v1 T1, v2 T2, v3 T3, v4 T4) *EffectT4[T, T1, T2, T3, T4] {
	r := this.f(v1, v2, v3, v4)
	if r == nil {
		panic("effect can't return nil")
	}
	return &EffectT4[T, T1, T2, T3, T4]{result: r, f: this.f}
}

func (this *EffectT4[T, T1, T2, T3, T4]) Result() *result.Result[T] {
	if this.result == nil {
		panic("effect was not executed")
	}
	return this.result
}

// Effect with 5 args
type EffectT5[T any, T1 any, T2 any, T3 any, T4 any, T5 any] struct {
	f      func(T1, T2, T3, T4, T5) *result.Result[T]
	result *result.Result[T]
}

func NewEffectT5[T any, T1 any, T2 any, T3 any, T4 any, T5 any](f func(T1, T2, T3, T4, T5) *result.Result[T]) *EffectT5[T, T1, T2, T3, T4, T5] {
	return &EffectT5[T, T1, T2, T3, T4, T5]{f: f}
}

func (this *EffectT5[T, T1, T2, T3, T4, T5]) ArgsCount() int {
	return 5
}

func (this *EffectT5[T, T1, T2, T3, T4, T5]) Run(v1 T1, v2 T2, v3 T3, v4 T4, v5 T5) *EffectT5[T, T1, T2, T3, T4, T5] {
	r := this.f(v1, v2, v3, v4, v5)
	if r == nil {
		panic("effect can't return nil")
	}
	return &EffectT5[T, T1, T2, T3, T4, T5]{result: r, f: this.f}
}

func (this *EffectT5[T, T1, T2, T3, T4, T5]) Result() *result.Result[T] {
	if this.result == nil {
		panic("effect was not executed")
	}
	return this.result
}

type Resource[T any] struct {
	Open  func() *result.Result[T]
	Close func()
}

func NewResource[T any](open func() *result.Result[T], close func()) *Resource[T] {
	return &Resource[T]{Open: open, Close: close}
}

func (this Resource[T]) IsResource() bool {
	return true
}
