package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/mobilemindtec/go-io/io"
	"github.com/mobilemindtec/go-io/option"
	"github.com/mobilemindtec/go-io/result"
	"github.com/mobilemindtec/go-io/types"
	gio "io"
	"log"
	"net/http"
	"reflect"
)

type HttpMethod string

const (
	GET    HttpMethod = "GET"
	POST   HttpMethod = "POST"
	PUT    HttpMethod = "PUT"
	DELETE HttpMethod = "DELETE"
	PATCH  HttpMethod = "PATCH"
	HEAD   HttpMethod = "HEAD"
)

type HttpEncoder[T any] interface {
	Encode(T) *result.Result[[]byte]
}

type HttpDecoder[T any] interface {
	Decode(data []byte) *result.Result[T]
}

type JsonEncoder[T any] struct {
}

func NewJsonEncoder[T any]() *JsonEncoder[T] {
	return &JsonEncoder[T]{}
}

func (this *JsonEncoder[T]) Encode(data T) *result.Result[[]byte] {
	return result.Try(func() ([]byte, error) {
		return json.Marshal(data)
	})
}

type JsonDecoder[T any] struct {
}

func NewJsonDecoder[T any]() *JsonDecoder[T] {
	return &JsonDecoder[T]{}
}

func (this *JsonDecoder[T]) Decode(data []byte) *result.Result[T] {
	return result.Try(func() (T, error) {
		typOf := reflect.TypeFor[T]()
		if typOf.Kind() == reflect.Pointer {
			typOf = typOf.Elem()
			val := reflect.New(typOf).Interface()
			return val.(T), json.Unmarshal(data, val)
		} else {
			val := reflect.New(typOf).Elem().Interface().(T)
			return val, json.Unmarshal(data, &val)
		}
	})
}

type Response[T any, E any] struct {
	Value       *option.Option[T] `json:"value"`
	StatusCode  int `json:"status_code"`
	RawBody     []byte `json:"-"`
	ErrorEntity *option.Option[E] `json:"error_entity"`
}

type HttpClient[T any, R any, E any] struct {
	debug        bool
	encoder      HttpEncoder[T]
	decoder      HttpDecoder[R]
	errorDecoder HttpDecoder[E]
	headers      map[string]string
}

func NewClient[T any, R any, E any]() *HttpClient[T, R, E] {
	return &HttpClient[T, R, E]{headers: map[string]string{}}
}

func (this *HttpClient[T, R, E]) Debug() *HttpClient[T, R, E] {
	this.debug = true
	return this
}

func (this *HttpClient[T, R, E]) AsJSON() *HttpClient[T, R, E] {
	this.headers["Content-Type"] = "application/json"
	this.headers["Accept"] = "application/json"
	this.encoder = NewJsonEncoder[T]()
	this.decoder = NewJsonDecoder[R]()
	this.errorDecoder = NewJsonDecoder[E]()
	return this
}

func (this *HttpClient[T, R, E]) SetErrorDecoder(decoder HttpDecoder[E]) *HttpClient[T, R, E] {
	this.errorDecoder = decoder
	return this
}

func (this *HttpClient[T, R, E]) Header(name string, value string) *HttpClient[T, R, E] {
	this.headers[name] = value
	return this
}

func (this *HttpClient[T, R, E]) Headers(vals ...string) *HttpClient[T, R, E] {

	if len(vals)%2 != 0 {
		log.Printf("headers should be even\n")
	}

	var name string
	for i, arg := range vals {
		if i%2 == 1 {
			this.headers[name] = arg
		} else {
			name = arg
		}
	}

	return this
}

func (this *HttpClient[T, R, E]) getPayloadOrNone(payload ...T) *option.Option[T] {
	if len(payload) > 0 {
		return option.Some(payload[0])
	}
	return option.None[T]()
}

func (this *HttpClient[T, R, E]) GetIO(url string, payload ...T) *types.IO[*Response[R, E]] {
	return io.IO[*Response[R, E]](
		io.Attempt(func() *result.Result[*Response[R, E]] {
			return this.Get(url, payload...)
		}))
}

func (this *HttpClient[T, R, E]) PostIO(url string, payload ...T) *types.IO[*Response[R, E]] {
	return io.IO[*Response[R, E]](
		io.Attempt(func() *result.Result[*Response[R, E]] {
			return this.Post(url, payload...)
		}))
}

func (this *HttpClient[T, R, E]) PutIO(url string, payload ...T) *types.IO[*Response[R, E]] {
	return io.IO[*Response[R, E]](
		io.Attempt(func() *result.Result[*Response[R, E]] {
			return this.Put(url, payload...)
		}))
}

func (this *HttpClient[T, R, E]) DeleteIO(url string, payload ...T) *types.IO[*Response[R, E]] {
	return io.IO[*Response[R, E]](
		io.Attempt(func() *result.Result[*Response[R, E]] {
			return this.Delete(url, payload...)
		}))
}

func (this *HttpClient[T, R, E]) PatchIO(url string, payload ...T) *types.IO[*Response[R, E]] {
	return io.IO[*Response[R, E]](
		io.Attempt(func() *result.Result[*Response[R, E]] {
			return this.Patch(url, payload...)
		}))
}

func (this *HttpClient[T, R, E]) HeadIO(url string, payload ...T) *types.IO[*Response[R, E]] {
	return io.IO[*Response[R, E]](
		io.Attempt(func() *result.Result[*Response[R, E]] {
			return this.Head(url, payload...)
		}))
}

func (this *HttpClient[T, R, E]) RequestIO(url string, method HttpMethod, payload *option.Option[T]) *types.IO[*Response[R, E]] {
	return io.IO[*Response[R, E]](
		io.Attempt(func() *result.Result[*Response[R, E]] {
			return this.Request(url, method, payload)
		}))
}

func (this *HttpClient[T, R, E]) Get(url string, payload ...T) *result.Result[*Response[R, E]] {
	return this.Request(url, GET, this.getPayloadOrNone(payload...))
}

func (this *HttpClient[T, R, E]) Delete(url string, payload ...T) *result.Result[*Response[R, E]] {
	return this.Request(url, DELETE, this.getPayloadOrNone(payload...))
}

func (this *HttpClient[T, R, E]) Patch(url string, payload ...T) *result.Result[*Response[R, E]] {
	return this.Request(url, PATCH, this.getPayloadOrNone(payload...))
}

func (this *HttpClient[T, R, E]) Head(url string, payload ...T) *result.Result[*Response[R, E]] {
	return this.Request(url, HEAD, this.getPayloadOrNone(payload...))
}

func (this *HttpClient[T, R, E]) Post(url string, payload ...T) *result.Result[*Response[R, E]] {
	return this.Request(url, POST, this.getPayloadOrNone(payload...))
}

func (this *HttpClient[T, R, E]) Put(url string, payload ...T) *result.Result[*Response[R, E]] {
	return this.Request(url, PUT, this.getPayloadOrNone(payload...))
}

func (this *HttpClient[T, R, E]) Request(url string, method HttpMethod, data *option.Option[T]) *result.Result[*Response[R, E]] {

	var req *http.Request
	var err error
	client := new(http.Client)

	if this.encoder == nil {
		if reflect.TypeFor[T]().Kind() != reflect.String {
			return result.OfError[*Response[R, E]](fmt.Errorf("encoder is required"))
		}
	}

	if this.debug {
		log.Printf("URL %v, METHOD = %v\n", url, method)
	}

	var payload *bytes.Buffer

	if data.NonEmpty() {
		if this.encoder != nil {
			res := this.encoder.Encode(data.Get())
			if res.IsError() {
				return result.OfError[*Response[R, E]](fmt.Errorf("payload encode error: %v", res.Failure().Error()))
			}
			payload = bytes.NewBuffer(res.Get())
		} else {
			payload = bytes.NewBufferString(data.GetValue().(string))
		}
		if this.debug {
			log.Printf("PAYLOAD = %v\n", payload.String())
		}
	}

	req, err = http.NewRequest(string(method), url, payload)

	if err != nil {
		return result.OfError[*Response[R, E]](fmt.Errorf("create request error: %v", err))
	}

	if this.debug {
		log.Printf("HEADERS = %v\n", this.headers)
	}

	for k, v := range this.headers {
		req.Header.Add(k, v)
	}

	res, err := client.Do(req)

	if err != nil {
		return result.OfError[*Response[R, E]](fmt.Errorf("do request error: %v", err))
	}

	defer res.Body.Close()

	body, err := gio.ReadAll(res.Body)

	if err != nil {
		return result.OfError[*Response[R, E]](fmt.Errorf("read reponse error: %v", err))
	}

	if this.debug {
		log.Printf("RESPONSE = %v\n", string(body))
	}

	switch res.StatusCode {

	case 200:

		if this.decoder != nil {

			decoded := this.decoder.Decode(body)

			if decoded.IsError() {
				return result.OfError[*Response[R, E]](
					fmt.Errorf("response decode error: %v", decoded.Failure().Error()))
			} else {
				return result.OfValue(&Response[R, E]{
					Value:       option.Of(decoded.Get()),
					StatusCode:  res.StatusCode,
					RawBody:     body,
					ErrorEntity: option.None[E](),
				})
			}

		} else if reflect.TypeFor[R]().Kind() == reflect.String {

			str := reflect.ValueOf(string(body)).Interface().(R)
			return result.OfValue(&Response[R, E]{
				Value:       option.Of(str),
				StatusCode:  res.StatusCode,
				RawBody:     body,
				ErrorEntity: option.None[E](),
			})

		} else {
			return result.OfValue(&Response[R, E]{
				Value:       option.None[R](),
				StatusCode:  res.StatusCode,
				RawBody:     body,
				ErrorEntity: option.None[E](),
			})
		}

	default:

		if this.errorDecoder != nil {
			decoded := this.errorDecoder.Decode(body)

			if decoded.IsError() {
				return result.OfError[*Response[R, E]](
					fmt.Errorf("decode error response error: %v", decoded.Failure().Error()))
			} else {
				return result.OfValue(&Response[R, E]{
					Value:       option.None[R](),
					ErrorEntity: option.Of(decoded.Get()),
					RawBody:     body,
					StatusCode:  res.StatusCode,
				})
			}
		}

		return result.OfValue(&Response[R, E]{
			Value:       option.None[R](),
			ErrorEntity: option.None[E](),
			RawBody:     body,
			StatusCode:  res.StatusCode,
		})
	}
}
