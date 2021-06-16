package contentuploader

import (
	"context"
	"io"

	"github.com/go-kit/kit/endpoint"
)

type UploadEndpoint endpoint.Endpoint

type EndPoints struct {
	UploadEndpoint
}

//Upload Endpoint
type uploadRequest struct {
	FileName string
	R        io.ReaderAt
	Size     int64
}

type uploadResponse struct {
	Err error `json:"error,omitempty"`
}

func (r uploadResponse) Error() error { return r.Err }

func MakeUploadEndPoint(s UploadSvc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(uploadRequest)
		s.Upload(ctx, req.FileName, req.R, req.Size)
		return uploadResponse{Err: nil}, nil
	}
}

func (e UploadEndpoint) Upload(ctx context.Context, fileName string, r io.ReaderAt, size int64) {
	request := uploadRequest{
		FileName: fileName,
		R:        r,
		Size:     size,
	}
	_, err := e(ctx, request)
	if err != nil {
		return
	}
}
