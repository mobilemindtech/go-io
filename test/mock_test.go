package test

import "github.com/mobilemindtech/go-io/either"

type Address struct {
	Street   string `json:"street"`
	Number   string `json:"number"`
	District string `json:"district"`
	ZipCode  string `json:"zip_code"`
	State    string `json:"state"`
	City     string `json:"city"`
}
type User struct {
	Email   string   `json:"email"`
	Name    string   `json:"name"`
	Cpf     string   `json:"cpf"`
	Cnpj    string   `json:"cnpj"`
	Address *Address `json:"address"`
}

type Person struct {
	Name string
	Age  int
}

type PersonPtr = *Person

type Validation struct {
	Messages []string
}

func NewValidation() *Validation {
	return &Validation{Messages: []string{}}
}

func ValidationOk() *Validation {
	return NewValidation()
}

func ValidationWith(msgs ...string) *Validation {
	return NewValidation().AddMessage(msgs...)
}

func (this *Validation) AddMessage(msgs ...string) *Validation {
	for _, msg := range msgs {
		this.Messages = append(this.Messages, msg)
	}
	return this
}

func (this *Validation) Count() int {
	return len(this.Messages)
}

func (this *Validation) Empty() bool {
	return this.Count() == 0
}

func (this *Validation) NonEmpty() bool {
	return !this.Empty()
}

func (this *Validation) Error() string {
	return "validation error"
}

type PersonValidator = *either.Either[*Validation, PersonPtr]
