package rpc

import (
	"encoding/json"

	"github.com/go-errors/errors"
	"gopkg.in/resty.v1"
)

type ErrorResponse struct {
	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}

type BaseClient struct {
	baseAddress string
}

func newBaseClient(baseAddress string) *BaseClient {
	return &BaseClient{
		baseAddress: baseAddress,
	}
}

func (bk *BaseClient) ReadResponse(resp *resty.Response) ([]byte, error) {
	body := resp.Body()
	if resp.StatusCode() != 200 {
		errResponse := ErrorResponse{}
		if err := json.Unmarshal(body, &errResponse); err != nil {
			return nil, errors.Errorf("cannot read error response: %s", string(body))
		}
		return nil, errors.Errorf(errResponse.ErrorText)
	}
	return body, nil
}
