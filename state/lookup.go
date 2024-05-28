package state

import (
	"fmt"
	"reflect"
)

func LookupVar(state *State, argType reflect.Type) (string, reflect.Value) {

	var item reflect.Value
	var key string
	var found bool
	tuples := state.ToTuples()

	for i := len(tuples) - 1; i >= 0; i-- {
		tp := tuples[i]
		rtype := reflect.TypeOf(tp.Val)
		if rtype == argType {
			item = reflect.ValueOf(tp.Val)
			found = true
			key = tp.Key
			break
		}
	}

	if found {
		state.Delete(key)
		return key, item
	}

	if argType == reflect.TypeFor[*State]() {
		return "__state__", reflect.ValueOf(state)
	}

	state.Dump()
	panic(fmt.Sprintf("var type %v not found on state", argType))
}
