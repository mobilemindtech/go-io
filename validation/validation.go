package validation

import (
	"fmt"
	"strings"
)

type Validation interface {
	IsSuccess() bool
	IsFailure() bool
	GetErrors() map[string]string
}

type Success struct {
}

func NewSuccess() Validation {
	return new(Success)
}

func (this *Success) IsSuccess() bool {
	return true
}

func (this *Success) IsFailure() bool {
	return false
}

func (this *Success) GetErrors() map[string]string {
	panic("success ha not errors")
}

type Failure struct {
	Errors map[string]string
}

func WithErrors(errs map[string]string) Validation {
	return &Failure{Errors: errs}
}

func NewFailure() Validation {
	return &Failure{Errors: map[string]string{}}
}

func (this *Failure) IsSuccess() bool {
	return false
}

func (this *Failure) IsFailure() bool {
	return true
}

func (this *Failure) GetErrors() map[string]string {
	return this.Errors
}

func (this *Failure) Error() string {
	var items []string
	for k, v := range this.Errors {
		items = append(items, fmt.Sprintf("%v: %v", k, v))
	}
	return strings.Join(items, ", ")
}
