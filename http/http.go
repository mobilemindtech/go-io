package http

import (
	"bytes"
	"fmt"
	"github.com/mobilemindtech/go-io/json"
	"github.com/mobilemindtech/go-io/option"
	"github.com/mobilemindtech/go-io/result"
	"io"
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

var (
	DefaultSuccessStatusCode = []int{200}
)

type HttpError[T any] struct {
	EntityError *option.Option[T]
	Message     string
}

func (this *HttpError[T]) Error() string {
	return fmt.Sprintf(this.Message)
}

type HttpEncoder[T any] interface {
	Encode(T) *result.Result[[]byte]
}

type HttpDecoder[T any] interface {
	Decode(data []byte) *result.Result[T]
}

type Response[T any, E any] struct {
	EntityBody  *option.Option[T] `json:"value"`
	StatusCode  int               `json:"status_code"`
	RawBody     []byte            `json:"-"`
	EntityError *option.Option[E] `json:"error_entity"`
	Header      http.Header       `json:"header"`
}

func (this *Response[T, E]) Body() string {
	return string(this.RawBody)
}

func (this *Response[T, E]) BodyAsResult() *result.Result[T] {
	if this.StatusCode != 200 {
		return result.OfError[T](
			&HttpError[E]{
				Message:     fmt.Sprintf("server return http status %v", this.StatusCode),
				EntityError: this.EntityError,
			})
	}

	if this.EntityBody.IsEmpty() {
		return result.OfError[T](
			&HttpError[E]{
				Message:     "server return empty a body",
				EntityError: this.EntityError,
			})
	}

	return result.OfValue(this.EntityBody.Get())
}

type Responser struct {
	StatusCode int
	Header     http.Header
	Body       io.ReadCloser
	Raw        *option.Option[*http.Response]
}

type DoRequest func(*http.Request) *result.Result[*Responser]

type HttpClient[Req, Resp, Err any] struct {
	debug             bool
	encoder           HttpEncoder[Req]
	decoder           HttpDecoder[Resp]
	errorDecoder      HttpDecoder[Err]
	headers           map[string]string
	successStatusList []int
	Requester         *option.Option[DoRequest]
}

func NewClient[Req, Resp, Err any]() *HttpClient[Req, Resp, Err] {
	return &HttpClient[Req, Resp, Err]{
		headers:           map[string]string{},
		successStatusList: DefaultSuccessStatusCode,
		Requester:         option.None[DoRequest]()}
}

func (this *HttpClient[Req, Resp, Err]) Debug() *HttpClient[Req, Resp, Err] {
	this.debug = true
	return this
}

func (this *HttpClient[Req, Resp, Err]) WithRequester(f DoRequest) *HttpClient[Req, Resp, Err] {
	this.Requester = option.Of(f)
	return this
}

func (this *HttpClient[Req, Resp, Err]) WithSuccessStatus(status ...int) *HttpClient[Req, Resp, Err] {
	for _, st := range status {
		this.successStatusList = append(this.successStatusList, st)
	}
	return this
}

func (this *HttpClient[Req, Resp, Err]) AsJSON() *HttpClient[Req, Resp, Err] {
	this.headers["Content-Type"] = "application/json"
	this.headers["Accept"] = "application/json"
	this.encoder = json.NewJsonEncoder[Req]()
	this.decoder = json.NewJsonDecoder[Resp]()
	this.errorDecoder = json.NewJsonDecoder[Err]()
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

	resResult := option.Or(
		option.Map(this.Requester,
			func(f DoRequest) *result.Result[*Responser] {
				return f(req)
			}), func() *result.Result[*Responser] {
			res, err := client.Do(req)
			if err != nil {
				return result.OfError[*Responser](err)
			}
			return result.OfValue(&Responser{
				Body:       res.Body,
				Header:     res.Header,
				StatusCode: res.StatusCode,
				Raw:        option.Of(res),
			})
		})

	if resResult.HasError() {
		return result.OfError[*Response[Resp, Err]](resResult.Failure())
	}

	res := resResult.Get()

	defer res.Body.Close()

	body, err := gio.ReadAll(res.Body)

	if err != nil {
		return result.OfError[*Response[Resp, Err]](fmt.Errorf("read reponse error: %v", err))
	}

	if this.debug {
		log.Printf("RESPONSE STATUS CODE %V, BODY = %v\n", res.StatusCode, string(body))
	}

	for _, status := range this.successStatusList {

		if status == res.StatusCode {
			if this.decoder != nil {

				decoded := this.decoder.Decode(body)

				if decoded.IsError() {
					return result.OfError[*Response[Resp, Err]](
						fmt.Errorf("response decode error: %v", decoded.Failure().Error()))
				} else {
					return result.OfValue(&Response[Resp, Err]{
						EntityBody:  option.Of(decoded.Get()),
						StatusCode:  res.StatusCode,
						RawBody:     body,
						EntityError: option.None[Err](),
						Header:      res.Header,
					})
				}

			} else if reflect.TypeFor[Resp]().Kind() == reflect.String {

				str := reflect.ValueOf(string(body)).Interface().(Resp)
				return result.OfValue(&Response[Resp, Err]{
					EntityBody:  option.Of(str),
					StatusCode:  res.StatusCode,
					RawBody:     body,
					EntityError: option.None[Err](),
					Header:      res.Header,
				})

			} else {
				return result.OfValue(&Response[Resp, Err]{
					EntityBody:  option.None[Resp](),
					StatusCode:  res.StatusCode,
					RawBody:     body,
					EntityError: option.None[Err](),
					Header:      res.Header,
				})
			}
		}
	}

	if this.errorDecoder != nil {
		decoded := this.errorDecoder.Decode(body)

		if decoded.IsError() {
			return result.OfError[*Response[Resp, Err]](
				fmt.Errorf("decode error response error: %v", decoded.Failure().Error()))
		} else {
			return result.OfValue(&Response[Resp, Err]{
				EntityBody:  option.None[Resp](),
				EntityError: option.Of(decoded.Get()),
				RawBody:     body,
				StatusCode:  res.StatusCode,
				Header:      res.Header,
			})
		}
	}

	return result.OfValue(&Response[Resp, Err]{
		EntityBody:  option.None[Resp](),
		EntityError: option.None[Err](),
		RawBody:     body,
		StatusCode:  res.StatusCode,
		Header:      res.Header,
	})
}
