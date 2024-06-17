package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/mobilemindtec/go-io/option"
	"github.com/mobilemindtec/go-io/result"
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
	StatusCode  int               `json:"status_code"`
	RawBody     []byte            `json:"-"`
	ErrorEntity *option.Option[E] `json:"error_entity"`
}

func (this *Response[T, E]) Body() string {
	return string(this.RawBody)
}

type HttpClient[Req, Resp, Err any] struct {
	debug        bool
	encoder      HttpEncoder[Req]
	decoder      HttpDecoder[Resp]
	errorDecoder HttpDecoder[Err]
	headers      map[string]string
}

func NewClient[Req, Resp, Err any]() *HttpClient[Req, Resp, Err] {
	return &HttpClient[Req, Resp, Err]{headers: map[string]string{}}
}

func (this *HttpClient[Req, Resp, Err]) Debug() *HttpClient[Req, Resp, Err] {
	this.debug = true
	return this
}

func (this *HttpClient[Req, Resp, Err]) AsJSON() *HttpClient[Req, Resp, Err] {
	this.headers["Content-Type"] = "application/json"
	this.headers["Accept"] = "application/json"
	this.encoder = NewJsonEncoder[Req]()
	this.decoder = NewJsonDecoder[Resp]()
	this.errorDecoder = NewJsonDecoder[Err]()
	return this
}

func (this *HttpClient[Req, Resp, Err]) SetErrorDecoder(decoder HttpDecoder[Err]) *HttpClient[Req, Resp, Err] {
	this.errorDecoder = decoder
	return this
}

func (this *HttpClient[Req, Resp, Err]) Header(name string, value string) *HttpClient[Req, Resp, Err] {
	this.headers[name] = value
	return this
}

func (this *HttpClient[Req, Resp, Err]) Headers(vals ...string) *HttpClient[Req, Resp, Err] {

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

func (this *HttpClient[Req, Resp, Err]) getPayloadOrNone(payload ...Req) *option.Option[Req] {
	if len(payload) > 0 {
		return option.Some(payload[0])
	}
	return option.None[Req]()
}

func (this *HttpClient[Req, Resp, Err]) Get(url string, payload ...Req) *result.Result[*Response[Resp, Err]] {
	return this.Request(url, GET, this.getPayloadOrNone(payload...))
}

func (this *HttpClient[Req, Resp, Err]) Delete(url string, payload ...Req) *result.Result[*Response[Resp, Err]] {
	return this.Request(url, DELETE, this.getPayloadOrNone(payload...))
}

func (this *HttpClient[Req, Resp, Err]) Patch(url string, payload ...Req) *result.Result[*Response[Resp, Err]] {
	return this.Request(url, PATCH, this.getPayloadOrNone(payload...))
}

func (this *HttpClient[Req, Resp, Err]) Head(url string, payload ...Req) *result.Result[*Response[Resp, Err]] {
	return this.Request(url, HEAD, this.getPayloadOrNone(payload...))
}

func (this *HttpClient[Req, Resp, Err]) Post(url string, payload ...Req) *result.Result[*Response[Resp, Err]] {
	return this.Request(url, POST, this.getPayloadOrNone(payload...))
}

func (this *HttpClient[Req, Resp, Err]) Put(url string, payload ...Req) *result.Result[*Response[Resp, Err]] {
	return this.Request(url, PUT, this.getPayloadOrNone(payload...))
}

func (this *HttpClient[Req, Resp, Err]) Request(url string, method HttpMethod, data *option.Option[Req]) *result.Result[*Response[Resp, Err]] {

	var req *http.Request
	var err error
	client := new(http.Client)

	if this.encoder == nil {
		if reflect.TypeFor[Req]().Kind() != reflect.String {
			return result.OfError[*Response[Resp, Err]](fmt.Errorf("encoder is required"))
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
				return result.OfError[*Response[Resp, Err]](fmt.Errorf("payload encode error: %v", res.Failure().Error()))
			}
			payload = bytes.NewBuffer(res.Get())
		} else {
			payload = bytes.NewBufferString(data.GetValue().(string))
		}
		if this.debug {
			log.Printf("PAYLOAD = %v\n", payload.String())
		}
		req, err = http.NewRequest(string(method), url, payload)
	} else {
		req, err = http.NewRequest(string(method), url, nil)
	}

	if err != nil {
		return result.OfError[*Response[Resp, Err]](fmt.Errorf("create request error: %v", err))
	}

	if this.debug {
		log.Printf("HEADERS = %v\n", this.headers)
	}

	for k, v := range this.headers {
		req.Header.Add(k, v)
	}

	res, err := client.Do(req)

	if err != nil {
		return result.OfError[*Response[Resp, Err]](fmt.Errorf("do request error: %v", err))
	}

	defer res.Body.Close()

	body, err := gio.ReadAll(res.Body)

	if err != nil {
		return result.OfError[*Response[Resp, Err]](fmt.Errorf("read reponse error: %v", err))
	}

	if this.debug {
		log.Printf("RESPONSE = %v\n", string(body))
	}

	switch res.StatusCode {

	case 200:

		if this.decoder != nil {

			decoded := this.decoder.Decode(body)

			if decoded.IsError() {
				return result.OfError[*Response[Resp, Err]](
					fmt.Errorf("response decode error: %v", decoded.Failure().Error()))
			} else {
				return result.OfValue(&Response[Resp, Err]{
					Value:       option.Of(decoded.Get()),
					StatusCode:  res.StatusCode,
					RawBody:     body,
					ErrorEntity: option.None[Err](),
				})
			}

		} else if reflect.TypeFor[Resp]().Kind() == reflect.String {

			str := reflect.ValueOf(string(body)).Interface().(Resp)
			return result.OfValue(&Response[Resp, Err]{
				Value:       option.Of(str),
				StatusCode:  res.StatusCode,
				RawBody:     body,
				ErrorEntity: option.None[Err](),
			})

		} else {
			return result.OfValue(&Response[Resp, Err]{
				Value:       option.None[Resp](),
				StatusCode:  res.StatusCode,
				RawBody:     body,
				ErrorEntity: option.None[Err](),
			})
		}

	default:

		if this.errorDecoder != nil {
			decoded := this.errorDecoder.Decode(body)

			if decoded.IsError() {
				return result.OfError[*Response[Resp, Err]](
					fmt.Errorf("decode error response error: %v", decoded.Failure().Error()))
			} else {
				return result.OfValue(&Response[Resp, Err]{
					Value:       option.None[Resp](),
					ErrorEntity: option.Of(decoded.Get()),
					RawBody:     body,
					StatusCode:  res.StatusCode,
				})
			}
		}

		return result.OfValue(&Response[Resp, Err]{
			Value:       option.None[Resp](),
			ErrorEntity: option.None[Err](),
			RawBody:     body,
			StatusCode:  res.StatusCode,
		})
	}
}
