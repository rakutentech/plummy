package client

import (
	"context"
	"io"
)

type Client interface {
	Render(ctx context.Context, engine string, req *RenderRequest) (*RenderResponse, error)
	WaitReady(ctx context.Context) error
}

type SizeReader interface {
	io.Reader
	Size() int64
}

type FileInput struct {
	Name     string
	Metadata []byte
	Reader   SizeReader
}

type FileData struct {
	Name     string
	Metadata []byte
	Data     []byte
}

type RenderRequest struct {
	RawParams []byte
	Files     []FileInput
}

type RenderResponse struct {
	RawParams []byte
	Files     []FileData
}
