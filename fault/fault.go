package fault

import (
	"errors"
	"fmt"
)

func AnyToError(err any) error {
	switch x := err.(type) {
	case string:
		return errors.New(x)
	case error:
		return x
	default:
		return fmt.Errorf("unknown panic: %v", err)
	}
}

// Check panic if error
func Check(err error, msg ...string) {
	if err != nil {
		if len(msg) > 0 {
			panic(fmt.Errorf(msg[0], err))
		}
		panic(err)
	}
}

func OrPanicF[T any](f func() (T, error), msg ...string) T {
	t, err := f()
	if err != nil {
		if len(msg) > 0 {
			panic(fmt.Errorf(msg[0], err))
		}
		panic(err)
	}
	return t
}

func OrPanic[T any](val T, err error) T {
	if err != nil {
		panic(err)
	}
	return val
}

func OrPanicWith[T any](val T, err error) func(string) T {
	return func(msg string) T {
		if err != nil {
			if len(msg) > 0 {
				panic(fmt.Errorf(msg, err))
			}
			panic(err)
		}
		return val
	}
}
