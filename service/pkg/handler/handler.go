package pkg

import (
	"context"
	"fmt"

	"github.com/ollama/ollama/api"
)

type RequestHandler struct {
	request  *api.GenerateRequest
	response api.GenerateResponse
}

func NewRequestHandler(model string) *RequestHandler {

	stream := false
	req := &api.GenerateRequest{
		Model:  model,
		Stream: &stream,
		Format: "json",
	}

	return &RequestHandler{
		request: req,
	}

}

func (rh *RequestHandler) ProcessRequest(ctx context.Context, prompt string) error {

	rh.request.Prompt = prompt

	client, err := api.ClientFromEnvironment()
	if err != nil {
		return err
	}

	err = client.Generate(ctx, rh.request, rh.processResponse)
	if err != nil {
		return err
	}

	fmt.Println(rh.response.Response)
	return nil

}

func (rh *RequestHandler) Response() string {
	return rh.response.Response
}

func (rh *RequestHandler) processResponse(resp api.GenerateResponse) error {
	rh.response = resp
	return nil
}
