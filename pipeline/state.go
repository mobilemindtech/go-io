package pipeline

import (
	"github.com/mobilemindtec/go-io/option"
	"reflect"
)

type Tuple struct {
	key string
	val interface{}
}

type StateItem struct {
	name  string
	value interface{}
	typ   reflect.Type
}

type State struct {
	items map[string]*StateItem
}

func NewState() *State {
	return &State{items: map[string]*StateItem{}}
}

func (this *State) SetVar(name string, value interface{}) *State {
	this.items[name] = &StateItem{name, value, reflect.TypeOf(value)}
	return this
}

func (this *State) Var(name string) interface{} {
	return this.items[name]
}

func (this *State) Delete(name string) *State {
	delete(this.items, name)
	return this
}

func (this *State) VarSafe(name string) *option.Option[any] {
	return option.Of[any](this.items[name])
}

func (this *State) Items() map[string]*StateItem {
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

func (this *State) ToCopy() *State {
	st := NewState()
	for k, v := range this.items {
		st.SetVar(k, v)
	}
	return st
}
