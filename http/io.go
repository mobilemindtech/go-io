package http

import (
	"github.com/mobilemindtec/go-io/io"
	"github.com/mobilemindtec/go-io/option"
	"github.com/mobilemindtec/go-io/result"
	"github.com/mobilemindtec/go-io/types"
)

func (this *HttpClient[Req, Resp, Err]) GetIO(url string, payload ...Req) *types.IO[*Response[Resp, Err]] {
	return io.IO[*Response[Resp, Err]](
		io.Attempt(func() *result.Result[*Response[Resp, Err]] {
			return this.Get(url, payload...)
		}))
}

func (this *HttpClient[Req, Resp, Err]) PostIO(url string, payload ...Req) *types.IO[*Response[Resp, Err]] {
	return io.IO[*Response[Resp, Err]](
		io.Attempt(func() *result.Result[*Response[Resp, Err]] {
			return this.Post(url, payload...)
		}))
}

func (this *HttpClient[Req, Resp, Err]) PutIO(url string, payload ...Req) *types.IO[*Response[Resp, Err]] {
	return io.IO[*Response[Resp, Err]](
		io.Attempt(func() *result.Result[*Response[Resp, Err]] {
			return this.Put(url, payload...)
		}))
}

func (this *HttpClient[Req, Resp, Err]) DeleteIO(url string, payload ...Req) *types.IO[*Response[Resp, Err]] {
	return io.IO[*Response[Resp, Err]](
		io.Attempt(func() *result.Result[*Response[Resp, Err]] {
			return this.Delete(url, payload...)
		}))
}

func (this *HttpClient[Req, Resp, Err]) PatchIO(url string, payload ...Req) *types.IO[*Response[Resp, Err]] {
	return io.IO[*Response[Resp, Err]](
		io.Attempt(func() *result.Result[*Response[Resp, Err]] {
			return this.Patch(url, payload...)
		}))
}

func (this *HttpClient[Req, Resp, Err]) HeadIO(url string, payload ...Req) *types.IO[*Response[Resp, Err]] {
	return io.IO[*Response[Resp, Err]](
		io.Attempt(func() *result.Result[*Response[Resp, Err]] {
			return this.Head(url, payload...)
		}))
}

func (this *HttpClient[Req, Resp, Err]) RequestIO(url string, method HttpMethod, payload *option.Option[Req]) *types.IO[*Response[Resp, Err]] {
	return io.IO[*Response[Resp, Err]](
		io.Attempt(func() *result.Result[*Response[Resp, Err]] {
			return this.Request(url, method, payload)
		}))
}
