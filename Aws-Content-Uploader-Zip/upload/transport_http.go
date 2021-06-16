package contentuploader

import (
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"

	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"

	"github.com/gorilla/mux"
)

type ErrorWrapper struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

type Error struct {
	Msg  string `json:"msg"`
	Code int    `json:"code"`
}

func EncodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(&ErrorWrapper{
		Success: true,
		Data:    response,
	})
}

func EncodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(&ErrorWrapper{
		Success: false,
		Data: Error{
			Msg:  "error",
			Code: 409,
		},
	})
}
func MakeHandler(ctx context.Context, s Service, logger kitlog.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(EncodeError),
	}

	uploadHandler := kithttp.NewServer(
		MakeUploadEndPoint(s),
		DecodeUploadRequest,
		EncodeResponse,
		opts...,
	)
	r := mux.NewRouter()
	r.Handle("/contentuploader/upload", uploadHandler).Methods(http.MethodPost)

	return r
}

const (
	multipartFileName = "file"
)

func DecodeUploadRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req uploadRequest
	err := r.ParseMultipartForm(200)
	if err != nil {
		return nil, err
	}

	file := r.MultipartForm.File[multipartFileName]
	if len(file) < 1 {
		return nil, err
	}

	req.FileName = file[0].Filename
	f, err := file[0].Open()
	if err != nil {
		return nil, err
	}

	req.R = f
	size, err := f.Seek(0, 2)
	if err != nil {
		return nil, err
	}
	req.Size = size

	_, err = f.Seek(0, 0)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func EncodeUploadRequest(_ context.Context, r *http.Request, request interface{}) error {
	req := request.(uploadRequest)
	pr, pw := io.Pipe()

	form := multipart.NewWriter(pw)
	go func() {
		defer pw.Close()

		iow, err := form.CreateFormFile(multipartFileName, req.FileName)
		if err != nil {
			return
		}

		_, err = io.Copy(iow, io.NewSectionReader(req.R, 0, req.Size))
		if err != nil {
			return
		}

		form.Close()
	}()

	r.Body = pr
	r.Header.Set("Content-Type", form.FormDataContentType())
	return nil
}
