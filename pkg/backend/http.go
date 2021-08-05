package backend

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"strings"

	"dynexo.de/ufyle/pkg/model/v1"
	"github.com/valyala/fasthttp"

	"dynexo.de/pkg/log"
	"dynexo.de/ufyle/pkg/vars"
)

func jsonError(message string, status int, ctx *fasthttp.RequestCtx) {

	if status == 0 {
		status = fasthttp.StatusNotAcceptable
	}

	if len(message) < 1 {
		message = "Request not acceptable"
	}

	data := map[string]interface{}{
		"success": false,
		"status":  status,
		"error":   message,
	}
	json, _ := json.Marshal(data)

	ctx.SetContentType("application/json")
	ctx.Write(json)
	ctx.SetStatusCode(fasthttp.StatusOK)
}

func jsonSuccess(data interface{}, status int, ctx *fasthttp.RequestCtx) {

	if status == 0 {
		status = fasthttp.StatusOK
	}

	resp := map[string]interface{}{
		"success": true,
		"status":  status,
		"data":    data,
	}
	json, _ := json.Marshal(resp)

	ctx.SetContentType("application/json")
	ctx.Write(json)
	ctx.SetStatusCode(fasthttp.StatusOK)
}

//
// WebService
//

type WebService struct {
	logger  log.ILogger
	ctrl    umodel.IController
	handler fasthttp.RequestHandler
}

func (s *WebService) rpcHandlerColfer(conn net.Conn) {

}

func (s *WebService) rpcHandlerProtoc(conn net.Conn) {

}

func (s *WebService) serviceUfylerV1Rpc(ctx *fasthttp.RequestCtx) {

	// hijack conn based on content-type header value
	switch string(ctx.Request.Header.ContentType()) {

	case "application/colfer":
		ctx.Hijack(s.rpcHandlerColfer)

	case "application/protoc":
		ctx.Hijack(s.rpcHandlerProtoc)

	default:
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
	}
}

func (s *WebService) serviceUfylerV1Post(ctx *fasthttp.RequestCtx) {

	ctype := "application/javascript"

	resp_err := func() {
		jsonError("Request inacceptable", fasthttp.StatusNotAcceptable, ctx)
	}

	if string(ctx.Request.Header.ContentType()) != ctype {
		resp_err()
	}

	// NOTE magic start

	// NOTE magic end

	// default result
	resp_err()
}

func (s *WebService) serviceUfylerV1Put(ctx *fasthttp.RequestCtx) {

	var res_path string

	if base_path, ok := ctx.UserValue("BasePath").(string); !ok {

		ctx.SetStatusCode(fasthttp.StatusNotFound)
		return

	} else {

		res_path = strings.TrimLeft(strings.TrimPrefix(string(ctx.Path()), base_path), "/")
	}

	if strings.HasSuffix(res_path, "/") {
		ctx.Error("Cannot write directory", fasthttp.StatusNotAcceptable)
		ctx.SetStatusCode(fasthttp.StatusNotAcceptable)
		return
	}

	s.logger.Debugf("PUT %s -> %s", string(ctx.Path()), res_path)

	// dbg(path, resPath)

	fs := s.ctrl.FsWrite()

	contentLength, err := ParseUint(ctx.Request.Header.Peek(fasthttp.HeaderContentLength))
	if err != nil || contentLength < 1 {
		ctx.Error("Content-length is missing", fasthttp.StatusExpectationFailed)
		ctx.SetStatusCode(fasthttp.StatusExpectationFailed)
		return
	}

	lockPath := fmt.Sprintf("%s.%s", res_path, vars.File_Lock_Suffix)

	// ensure no lock file exists
	if fi, err := fs.Stat(lockPath); err == nil && fi.Size() > 0 {
		ctx.Error("Cannot open requested path", fasthttp.StatusConflict)
		ctx.SetStatusCode(fasthttp.StatusConflict)
		return
	}

	// we only support files and no folders
	fi, err := fs.Stat(res_path)
	if err == nil && fi.IsDir() {
		ctx.Error("Cannot write directory", fasthttp.StatusNotAcceptable)
		ctx.SetStatusCode(fasthttp.StatusNotAcceptable)
		return
	}

	hdr := &ctx.Response.Header
	statusCode := fasthttp.StatusOK

	r, err := fs.OpenFile(res_path, os.O_TRUNC|os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		s.logger.Errorf("Cannot obtain file writer for path=%q: %s", res_path, err)
		ctx.Error("Internal Server Error", fasthttp.StatusInternalServerError)
		return
	}

	if err := ctx.Request.BodyWriteTo(r); err != nil {
		s.logger.Errorf("Failed to write to path=%q: %s", res_path, err)
		ctx.Error("Internal Server Error", fasthttp.StatusInternalServerError)
		return
	}

	// set content-location and return 200 HTTP OK
	hdr.SetCanonical(strContentLocation, ctx.Path())
	ctx.SetStatusCode(statusCode)
}

func (s *WebService) serviceUfylerV1Get(ctx *fasthttp.RequestCtx) {

	ctype := "application/octet-stream"

	// resp_err := func() {
	// 	ctx.SetContentType(ctype)
	// 	ctx.SetStatusCode(fasthttp.StatusNotFound)
	// }

	// NOTE magic start

	// TODO check access rights

	base_path, ok := ctx.UserValue("BasePath").(string)
	if !ok {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		return
	}

	res_path := strings.Trim(strings.TrimPrefix(string(ctx.Path()), base_path), "/")

	if strings.HasSuffix(res_path, "/") {
		ctx.Error("Cannot read directory", fasthttp.StatusNotAcceptable)
		ctx.SetStatusCode(fasthttp.StatusNotAcceptable)
		return
	}

	s.logger.Debugf("GET %s -> %s", string(ctx.Path()), res_path)

	// dbg(path, resPath)

	fs := s.ctrl.FsRead()

	fi, err := fs.Stat(res_path)
	if err != nil {
		ctx.Error("Cannot open requested path", fasthttp.StatusNotFound)
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		return
	}

	// we only support files and no folders
	if fi.IsDir() {
		ctx.Error("Cannot read directory", fasthttp.StatusNotFound)
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		return
	}

	hdr := &ctx.Response.Header
	statusCode := fasthttp.StatusOK

	contentLength := int(fi.Size())
	if int64(contentLength) != fi.Size() {
		ctx.Error("Invalid file size", fasthttp.StatusConflict)
		ctx.SetStatusCode(fasthttp.StatusConflict)
	}

	lastModifiedStr := AppendHTTPDate(nil, fi.ModTime())

	r, err := fs.Open(res_path)
	if err != nil {
		s.logger.Errorf("cannot obtain file reader for path=%q: %s", res_path, err)
		ctx.Error("Internal Server Error", fasthttp.StatusInternalServerError)
		return
	}

	hdr.SetCanonical(strLastModified, lastModifiedStr)

	if !ctx.IsHead() {
		ctx.SetBodyStream(r, contentLength)

	} else {
		ctx.Response.ResetBody()
		ctx.Response.SkipBody = true
		ctx.Response.Header.SetContentLength(contentLength)

		if rc, ok := r.(io.Closer); ok {
			if err := rc.Close(); err != nil {
				s.logger.Error("Failed to close file reader")
				ctx.Error("Internal Server Error", fasthttp.StatusInternalServerError)
				return
			}
		}
	}

	// set right content type and return 200 HTTP OK
	ctx.SetContentType(ctype)
	ctx.SetStatusCode(statusCode)
}

func NewWebService(ctrl umodel.IController, logger log.ILogger) *WebService {

	s := &WebService{
		ctrl:   ctrl,
		logger: logger,
	}

	requestHandler := func(ctx *fasthttp.RequestCtx) {

		logger.Debugf("Request: Method=%s, Path=%s, Type=%s", string(ctx.Method()), string(ctx.Path()),
			string(ctx.Request.Header.ContentType()))

		{
			basePath := "/ufyler/v1"

			if strings.HasPrefix(string(ctx.Path()), basePath) {

				ctx.SetUserValue("BasePath", basePath)

				switch string(ctx.Method()) {
				case "RPC":
					s.serviceUfylerV1Rpc(ctx)

				case "POST":
					s.serviceUfylerV1Post(ctx)

				case "GET":
					s.serviceUfylerV1Get(ctx)

				case "PUT":
					s.serviceUfylerV1Put(ctx)
				}

			} else {
				ctx.Error("Unsupported request method", fasthttp.StatusNotFound)
			}

		}
	}

	s.handler = requestHandler

	return s
}

func ListenAndServeService(addr string, s umodel.IController, logger log.ILogger) error {

	address := ":8449"
	if len(addr) > 2 {
		address = addr
	}

	service := NewWebService(s, logger)

	return ListenAndServeFast(address, service.handler, logger)
}
