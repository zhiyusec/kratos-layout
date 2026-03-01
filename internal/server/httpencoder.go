package server

import (
	"net/http"
	"strings"

	httpx "github.com/go-kratos/kratos/v2/transport/http"
	statusx "github.com/go-kratos/kratos/v2/transport/http/status"

	"google.golang.org/grpc/status"
)

const (
	baseContentType = "application"
)

// ContentType returns the content-type with base prefix.
func ContentType(subtype string) string {
	return baseContentType + "/" + subtype
}

// ContentSubtype returns the content-subtype for the given content-type. The
// given content-type must be a valid content-type that starts with
// but no content-subtype will be returned.
// according rfc7231.
// contentType is assumed to be lowercase already.
func ContentSubtype(contentType string) string {
	left := strings.Index(contentType, "/")
	if left == -1 {
		return ""
	}
	right := strings.Index(contentType, ";")
	if right == -1 {
		right = len(contentType)
	}
	if right < left {
		return ""
	}
	return contentType[left+1 : right]
}

type HTTPResponseEncoder struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

func ResponseEncoder(w http.ResponseWriter, r *http.Request, v any) error {
	if v == nil {
		return nil
	}
	if rd, ok := v.(httpx.Redirector); ok {
		url, code := rd.Redirect()
		http.Redirect(w, r, url, code)
		return nil
	}
	resp := &HTTPResponseEncoder{
		Code: http.StatusOK,
		Msg:  "success",
		Data: v,
	}
	codec, _ := httpx.CodecForRequest(r, "Accept")
	data, err := codec.Marshal(resp)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", ContentType(codec.Name()))
	_, err = w.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func ErrorEncoder(w http.ResponseWriter, r *http.Request, err error) {
	if err == nil {
		return
	}

	resp := &HTTPResponseEncoder{}
	if gs, ok := status.FromError(err); ok {
		resp = &HTTPResponseEncoder{
			Code: statusx.FromGRPCCode(gs.Code()),
			Msg:  gs.Message(),
			Data: nil,
		}
	} else {
		resp = &HTTPResponseEncoder{
			Code: http.StatusInternalServerError,
			Msg:  "internal server error",
		}
	}

	// se := errors.FromError(err)
	codec, _ := httpx.CodecForRequest(r, "Accept")
	body, err := codec.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", ContentType(codec.Name()))
	w.WriteHeader(resp.Code)
	_, _ = w.Write(body)
}
