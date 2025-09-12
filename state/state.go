package state

import (
	"fmt"
	"github.com/mobilemindtech/go-io/option"
	"github.com/mobilemindtech/go-io/util"
	"log"
	"reflect"
)

type State struct {
	items map[string]interface{}
}

func NewState() *State {
	return &State{items: map[string]interface{}{}}
}

func (this *State) SetVar(name string, value interface{}) *State {
	this.items[name] = value
	return this
}

func (this *State) Var(name string) interface{} {
	return this.items[name]
}

func (this *State) Consume(name string) interface{} {
	val := this.items[name]
	delete(this.items, name)
	return val
}

func (this *State) Delete(name string) *State {
	delete(this.items, name)
	return this
}

func (this *State) VarSafe(name string) *option.Option[any] {
	return option.Of[any](this.items[name])
}

func (this *State) Items() map[string]interface{} {
	return this.items
}

func (this *State) Count() int {
	return len(this.items)
}

func (this *State) ToTuples() []*Tuple {
	var tuples []*Tuple
	for k, val := range this.items {
		tuples = append(tuples, &Tuple{k, val})
	}
	return tuples
}

func (this *State) Copy() *State {
	st := NewState()
	for k, v := range this.items {
		st.SetVar(k, v)
	}
	return st
}

func (this *State) Dump() {
	log.Printf("==> state dump start\n")
	for k, v := range this.items {
		log.Printf("%v=%v\n", k, v)
	}
	log.Printf("==> state dump end\n")
}

func ConsumeOf[T any](state *State, key string) T {
	val := state.Consume(key)
	if util.IsNotNil(val) {
		if v, ok := val.(T); ok {
			return v
		}
	}
	panic(fmt.Sprintf("var %v not found on state", key))
}

func VarOf[T any](state *State, key string) T {
	val := state.Var(key)
	if util.IsNotNil(val) {
		if v, ok := val.(T); ok {
			return v
		}
	}
	panic(fmt.Sprintf("var %v not found on state", key))
}

func Consume[T any](st *State) T {
	_, val := LookupVar(st, reflect.TypeFor[T](), true)
	return val.Interface().(T)
}

func Var[T any](st *State) T {
	_, val := LookupVar(st, reflect.TypeFor[T](), false)
	return val.Interface().(T)
}
