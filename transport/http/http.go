package chttp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	cerror "Profile/common/model/error"

	renderer "Profile/common/renderer"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/sd/lb"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/google/go-querystring/query"
)

type Errorer interface {
	Error() error
}

type Error struct {
	Msg  string `json:"msg"`
	Code int    `json:"code"`
}

type ErrorWrapper struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

type ErrorWrapperJson struct {
	Success bool            `json:"success"`
	Data    json.RawMessage `json:"data"`
}

func EncodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(Errorer); ok && e.Error() != nil {
		EncodeError(ctx, e.Error(), w)
		return nil
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(&ErrorWrapper{
		Success: true,
		Data:    response,
	})
}

func RenderResponse(renderer renderer.Renderer, templateName string) kithttp.EncodeResponseFunc {
	return func(ctx context.Context, w http.ResponseWriter, response interface{}) error {
		if e, ok := response.(Errorer); ok && e.Error() != nil {
			fmt.Fprintf(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return nil
		}

		w.Header().Set("Content-Type", "text/html")
		return renderer.Render(templateName, w, response)
	}
}

func encodeErrorEncoderToEncodeResponse(encodeError kithttp.ErrorEncoder) kithttp.EncodeResponseFunc {
	return func(ctx context.Context, w http.ResponseWriter, response interface{}) (err error) {
		if e, ok := response.(error); ok {
			encodeError(ctx, e, w)
			return
		}
		return
	}
}

func encodeResponseToErrorEncoder(enc kithttp.EncodeResponseFunc) kithttp.ErrorEncoder {
	return func(ctx context.Context, err error, w http.ResponseWriter) {
		enc(ctx, w, err)
	}
}

// func SignedErrorEncoder(signer apisigner.ServerSigner) kithttp.ServerOption {
// 	return kithttp.ServerErrorEncoder(encodeResponseToErrorEncoder(signer.SignRes(encodeErrorEncoderToEncodeResponse(EncodeError))))
// }

// encode errors from business-logic
func EncodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	var cerrorObj cerror.Cerror
	switch err.(type) {
	case *cerror.Cerror:
		cerrorObj = *err.(*cerror.Cerror)
	default:
		cerrorObj = *cerror.New(-1, err.Error()).(*cerror.Cerror)
	}

	json.NewEncoder(w).Encode(&ErrorWrapper{
		Success: false,
		Data: Error{
			Msg:  cerrorObj.Msg(),
			Code: cerrorObj.Code(),
		},
	})
}

func DecodeResponse(ctx context.Context, r *http.Response, resp interface{}) error {

	if r.StatusCode != 200 {
		return fmt.Errorf("http error, statuscode:%d status:%s", r.StatusCode, r.Status)
	}

	var wrapper ErrorWrapperJson
	err := json.NewDecoder(r.Body).Decode(&wrapper)
	if err != nil {
		return err
	}

	if !wrapper.Success {
		return decodeError(wrapper.Data)
	}

	return json.Unmarshal(wrapper.Data, &resp)
}

func decodeError(r json.RawMessage) error {
	var e Error
	if err := json.Unmarshal(r, &e); err != nil {
		return err
	}

	return cerror.New(e.Code, e.Msg)
}

func ErrorDecoder(r *http.Response) error {
	var wrapper ErrorWrapperJson
	err := json.NewDecoder(r.Body).Decode(&wrapper)
	if err != nil {
		return err
	}

	if wrapper.Success {
		return nil
	}

	return decodeError(wrapper.Data)
}

// EncodeHTTPGenericRequest is a transport/http.EncodeRequestFunc that
// JSON-encodes any request to the request body. Primarily useful in a client.
func EncodeHTTPGenericRequest(_ context.Context, r *http.Request, request interface{}) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		return err
	}
	r.Body = ioutil.NopCloser(&buf)
	return nil
}

func EncodeHTTPGetDeleteGenericRequest(_ context.Context, r *http.Request, request interface{}) error {
	v, _ := query.Values(request)
	r.URL.RawQuery = v.Encode()
	return nil
}

func BalancerToEndpoint(b lb.Balancer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		e, err := b.Endpoint()
		if err != nil {
			return nil, err
		}

		return e(ctx, request)
	}
}
