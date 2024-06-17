package http

import (
	"github.com/mobilemindtec/go-io/option"
	"github.com/mobilemindtec/go-io/result"
	"github.com/mobilemindtec/go-io/rio"
)

func (this *HttpClient[Req, Resp, Err]) GetRIO(url string, payload ...Req) *rio.IO[*Response[Resp, Err]] {
	return rio.Attempt(func() *result.Result[*Response[Resp, Err]] {
		return this.Get(url, payload...)
	})
}

func (this *HttpClient[Req, Resp, Err]) PostRIO(url string, payload ...Req) *rio.IO[*Response[Resp, Err]] {
	return rio.Attempt(func() *result.Result[*Response[Resp, Err]] {
		return this.Post(url, payload...)
	})
}

func (this *HttpClient[Req, Resp, Err]) PutRIO(url string, payload ...Req) *rio.IO[*Response[Resp, Err]] {

	return rio.Attempt(func() *result.Result[*Response[Resp, Err]] {
		return this.Put(url, payload...)
	})
}

func (this *HttpClient[Req, Resp, Err]) DeleteRIO(url string, payload ...Req) *rio.IO[*Response[Resp, Err]] {
	return rio.Attempt(func() *result.Result[*Response[Resp, Err]] {
		return this.Delete(url, payload...)
	})
}

func (this *HttpClient[Req, Resp, Err]) PatchRIO(url string, payload ...Req) *rio.IO[*Response[Resp, Err]] {
	return rio.Attempt(func() *result.Result[*Response[Resp, Err]] {
		return this.Patch(url, payload...)
	})
}

func (this *HttpClient[Req, Resp, Err]) HeadRIO(url string, payload ...Req) *rio.IO[*Response[Resp, Err]] {
	return rio.Attempt(func() *result.Result[*Response[Resp, Err]] {
		return this.Head(url, payload...)
	})
}

func (this *HttpClient[Req, Resp, Err]) RequestRIO(url string, method HttpMethod, payload *option.Option[Req]) *rio.IO[*Response[Resp, Err]] {
	return rio.Attempt(func() *result.Result[*Response[Resp, Err]] {
		return this.Request(url, method, payload)
	})
}
